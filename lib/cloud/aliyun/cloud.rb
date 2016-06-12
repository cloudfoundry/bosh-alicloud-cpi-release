module Bosh::Aliyun

  class Cloud < Bosh::Cloud

    include Helpers

    attr_reader :options
    attr_reader :registry
    attr_reader :agent_options

    def initialize options
      @logger = Bosh::Clouds::Config.logger

      @options = recursive_symbolize_keys(options["aliyun"])
      validate_options

      registry_options = recursive_symbolize_keys(options["registry"])
      initialize_registry registry_options

      @agent_options = options["agent"]
      @logger.debug "current options is #{options.inspect}"

      @aliyun_client = Bosh::Aliyun::Client.new @options, @logger
    end

    def create_stemcell image_path = nil, cloud_properties = nil

      @logger.info "fake stemcell creation for now"
      # TODO
      # Upload image from local path

      @options[:ImageId]
    end

    def delete_stemcell stemcell_id
      @logger.info "fake stemcell deletion for now"
      nil
    end

    #  {"network_a" =>
    #    {
    #      "netmask"          => "255.255.248.0",
    #      "ip"               => "172.30.41.40",
    #      "gateway"          => "172.30.40.1",
    #      "dns"              => ["172.30.22.153", "172.30.22.154"],
    #      "cloud_properties" => {"name" => "VLAN444"}
    #    }
    #  }

    def create_vm(agent_id=nil, stemcell_id=nil, resource_pool=nil, networks=nil, disk_locality=nil, env=nil)

      @logger.info "start to create a vm"
      param = {}
      param[:ImageId] = stemcell_id || @options[:ImageId]

      nm = nil
      if not networks.nil?
        @logger.debug "networks param is: #{networks.inspect}"
        nm = NetworkManager.new networks
        param.merge! nm.network.configure
      end

      param.merge! prepare_create_vm_parameters
      @logger.debug "current param is: #{param.inspect}"

      vm_created = false
      begin
        res = @aliyun_client.CreateInstance param
        ins_id = res["InstanceId"]

        while is_vm_pending? ins_id
          sleep 10
        end
        vm_created = true
        @logger.debug "created a vm, the vm id is #{ins_id}"

        if not networks.nil?
          if nm.vip?
            # Do not assign, just bind
            bind_public_eip ins_id, nm
            @logger.debug "assigned a public ip for the newly created vm, try to start it"
          end
        end

        # VM will be started after creation
        start_vm ins_id
        @logger.debug "the vm creation is done"

        if not @registry.nil?
          registry_settings = initial_agent_settings(
            ins_id,
            agent_id,
            networks,
            env
          )
          @logger.debug "registry_settings is #{registry_settings}"
          @registry.update_settings(ins_id, registry_settings)
        end

        ins_id
      rescue => e
        delete_vm ins_id if vm_created
        @logger.error("failed to start a new vm. #{e.inspect}\n#{e.backtrace.join('\n')}")
        raise Bosh::Clouds::VMCreationFailed.new(false), "failed to start a new vm. #{e.inspect}\n#{e.backtrace.join("\n")}"
      end

    end

    def delete_vm vm_id
      @logger.info "start to delete a vm"
      param={
        :InstanceId => vm_id
      }

      vm_stopped = true
      begin
        if not is_vm_stopped? vm_id
          @logger.debug "stopping the vm"
          vm_stopped = false

          stop_vm vm_id

          vm_stopped = true
          @logger.debug "the vm is stopped"
        end
        r = @aliyun_client.DeleteInstance param if vm_stopped
        @logger.debug "the vm is deleted"

        r
      rescue => e
        @logger.error("failed to delete the vm. #{e.inspect}\n#{e.backtrace.join('\n')}")
        raise Bosh::Clouds::VMNotFound.new(false), "failed to delete the vm. #{e.inspect}\n#{e.backtrace.join("\n")}"
      end
    end

    def has_vm? vm_id
      param = {
        :RegionId => @options[:RegionId],
        :InstanceIds => "[\"#{vm_id}\"]"
      }
      r = @aliyun_client.DescribeInstances param
      r["Instances"]["Instance"].count != 0
    end

    def reboot_vm vm_id
      @logger.info "start to reboot a vm"
      stop_vm vm_id if is_vm_running? vm_id
      start_vm vm_id if is_vm_stopped? vm_id
      @logger.info "the vm is rebooted"
    end

    def stop_it vm_id
      stop_vm vm_id
    end

    def initialize_registry registry_properties

      registry_endpoint   = registry_properties.fetch(:endpoint)
      registry_user       = registry_properties.fetch(:user)
      registry_password   = registry_properties.fetch(:password)

      @registry = Bosh::Registry::Client.new(registry_endpoint, registry_user, registry_password)

      begin
        @registry.read_settings "check"
      rescue Errno::ECONNREFUSED => e
        @logger.info("failed to read settings from registry endpoint. Continue...")
        @registry = nil
      end

    end

    # Generates initial agent settings. These settings will be read by agent
    # from BOSH registry on a target instance. Disk
    # conventions for amazon are:
    # system disk: /dev/sda
    # ephemeral disk: /dev/sdb
    # EBS volumes can be configured to map to other device names later (sdf
    # through sdp, also some kernels will remap sd* to xvd*).
    #
    # @param [String] agent_id Agent id (will be picked up by agent to
    #   assume its identity
    # @param [Hash] network_spec Agent network spec
    # @param [Hash] environment
    #   keys are device type ("ephemeral", "raw_ephemeral") and values are array of strings representing the
    #   path to the block device. It is expected that "ephemeral" has exactly one value.
    # @return [Hash]
    def initial_agent_settings(ins_id, agent_id, network_spec, environment)
      settings = {
          "vm" => {
            "name" => ins_id
          },
          "agent_id" => agent_id,
          "networks" => agent_network_spec,
           "disks" => {
               "system" => "/dev/xvda",
               "ephemeral" => "/dev/xvdb"
           }
      }

      # TODO Will add this two later
      # @param [String] root_device_name root device, e.g. /dev/sda1
      # @param [Hash] block_device_agent_info disk attachment information to merge into the disks section.
      # settings["disks"].merge!(block_device_agent_info)
      # settings["disks"]["ephemeral"] = settings["disks"]["ephemeral"][0]["path"]

      settings["env"] = environment if environment

      @logger.debug "current agent settings is #{@agent_options.inspect}"
      @logger.debug "current settings is #{settings.inspect}"

      settings.merge(@agent_options)
    end

    private

    def agent_network_spec
      {
          "private"=> {
            "type"=> "vip"
          },
          "public"=> {
            "type"=> "vip"
          }
      }
    end

    def bind_public_eip vm_id, nm
      @logger.info "Start to bind public eip with the specific vm"
      param = {
        :RegionId => @options[:RegionId],
        :Status => "Available",
        :EipAddress => nm.vip_network.ip
      }

      r = @aliyun_client.DescribeEipAddresses param
      @logger.debug "check eip exist or not, #{r.inspect}"

      allocation_id = r['EipAddresses']['EipAddress'][0]['AllocationId']
      param = {
        :InstanceId => vm_id,
        :AllocationId => allocation_id
      }
      r = @aliyun_client.AssociateEipAddress param
      @logger.debug "bond the newly created eip with the vm"
    end

    def start_vm vm_id
      param = {
        :InstanceId => vm_id
      }
      @aliyun_client.StartInstance param if is_vm_stopped? vm_id

      count = 1
      @logger.debug "starting the vm"
      until is_vm_running? vm_id
        @logger.debug "down. ping #{count} time"
        count += 1
        sleep 10
      end
      @logger.debug "the vm is started"
    end

    def stop_vm vm_id
      param = {
        :InstanceId => vm_id,
        :ForceStop => "true"
      }
      r = @aliyun_client.StopInstance param if is_vm_running? vm_id

      count = 1
      @logger.debug "stopping the vm"
      until is_vm_stopped? vm_id
        @logger.debug "up. ping #{count} time"
        count += 1
        sleep 10
      end
      @logger.debug "the vm is stopped"
    end

    def vm_status vm_id
      param = {
        :RegionId => @options[:RegionId],
        :InstanceIds => "[\"#{vm_id}\"]"
      }
      r = @aliyun_client.DescribeInstances param

      r["Instances"]["Instance"][0]["Status"]
    end

    def is_vm_pending? vm_id
      vm_status(vm_id) == "Pending"
    end

    def is_vm_running? vm_id
      r = vm_status(vm_id)
      until r != "Pending"
        r = vm_status vm_id
        sleep 10
      end
      r == "Running"
    end

    def is_vm_stopped? vm_id
      r = vm_status(vm_id)
      until r != "Pending"
        r = vm_status vm_id
        sleep 10
      end
      r == "Stopped"
    end

    def prepare_create_vm_parameters
      para = {
        :RegionId => @options[:RegionId],
        :InstanceType => @options[:InstanceType],
        :InstanceName => @options[:InstanceName],
        :Description => @options[:Description],
        :HostName => @options[:HostName],
        :Password => @options[:Password]
      }
    end

    def validate_options
      required_keys = [
          :RegionId,
          :InstanceType,
          :ImageId,
          :AccessKeyId,
          :AccessKey,
          :Password,
          :InstanceName,
          :HostName
      ]

      missing_keys = []
      required_keys.each do |key|
        if !@options.has_key?(key)
          missing_keys << "#{key}:"
        end
      end

      raise ArgumentError, "missing configuration parameters > #{missing_keys.join(', ')}" unless missing_keys.empty?
    end

  end
end
