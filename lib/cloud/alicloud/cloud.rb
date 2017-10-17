module Bosh::Alicloud

  class Cloud < Bosh::Cloud

    include Helpers

    attr_reader :registry_options
    attr_reader :alicloud_options
    attr_reader :agent_options

    def initialize options
      @logger = Bosh::Clouds::Config.logger

      @alicloud_options = recursive_symbolize_keys(options["alicloud"])
      validate_options

      @registry_options = recursive_symbolize_keys(options["registry"])
      initialize_registry

      @agent_options = options["agent"]
      @logger.debug "current options is #{options.inspect}"

      @alicloud_client = Bosh::Alicloud::Client.new @alicloud_options, @logger
    end

    ##
    # Creates a stemcell
    #
    # Creates a reusable VM image in the IaaS from the stemcell image
    #
    # @param image_path [String] Path to the stemcell image extracted from the stemcell tarball on a local filesystem
    # @param cloud_properties [Hash] Cloud properties hash extracted from the stemcell tarball
    #
    # @return stemcell_id [String] Cloud ID of the created stemcell
    #
    def create_stemcell image_path = nil, cloud_properties = nil

      @logger.info "Start to create a stemcel;"

      @logger.info "Cloud properties are: #{cloud_properties}"
      region = @alicloud_options[:RegionId]

      @logger.debug "Region: #{region}"

      if not cloud_properties.nil?

        begin

          cloud_properties = recursive_symbolize_keys(cloud_properties)

          if cloud_properties.has_key?(:image_id)
            images = cloud_properties[:image_id]
            @logger.debug "images: #{images}"

            if images.has_key?(region.to_sym)

              image_id = images[region.to_sym]
              @logger.debug "create image: #{image_id} with region: #{region}"

              @alicloud_options[:ImageId] = image_id

              return image_id

            end

          end

        rescue => e
          @logger.error("failed to create stemcell. #{e.inspect}\n#{e.backtrace.join('\n')}")
          raise Bosh::Clouds::StemcellCreationFailed.new(false), "failed to create stemcell. #{e.inspect}\n#{e.backtrace.join("\n")}"
        end

      end

      @logger.info "current alicloud options: #{alicloud_options}"

      @alicloud_options[:ImageId]
    end

    ##
    # Deletes a stemcell
    #
    # Delete previously created stemcell
    #
    # @param stemcell_id [String] Cloud ID of the stemcell to delete. Returned from {#create_stemcell}
    #
    # @return [void]
    #
    def delete_stemcell stemcell_id
      @logger.info "fake stemcell deletion for now"
      nil
    end

    ##
    # Creates a VM
    #
    #    Creates a new VM based on the stemcell.
    #
    #    @param agent_id [String] Bosh director will use it communicate with agent
    #    @param stemcell_id [String] An UUId specifies the stemcell when creating a vm. Returned by {#create_stemcell}
    #    @param resource_pool [Hash] Resource properties define this VM resource. Cloud specified.
    #    @param network [Hash] Network properties define this VM's network. Used by network manager
    #    @param disk_locality [Array of strings] Array of disk cloud IDs for each disk that created VM will most likely be attached
    #    @param env [Hash] environment that will be passed to this vm
    #    @return vm_id [String] Cloud ID of the created VM. Used by {#reboot_vm}, {#delete_vm}, {#attach_disk} and {#detach_disk}
    def create_vm(agent_id=nil, stemcell_id=nil, resource_pool=nil, network=nil, disk_locality=nil, env=nil)

      @logger.info "start to create a vm"
      @logger.info "current alicloud options: #{alicloud_options}"
      param = {}

      nm = nil
      if not network.nil?
        @logger.debug "networks param is: #{network.inspect}"
        nm = NetworkManager.new network
        param.merge! nm.network.configure
      end

      param.merge! prepare_create_vm_parameters

      # disk params
      if not resource_pool.nil?
        @logger.info "resource_pool param is: #{resource_pool.inspect}"
        resource_pool = recursive_symbolize_keys(resource_pool)
        validate_resource_pool resource_pool

        param[:ImageId] = stemcell_id || resource_pool[:image_id]

        begin

          param[:InstanceType] = resource_pool[:instance_type]
          # param[:ImageId] = resource_pool[:image_id]

          ephemeral_disk = resource_pool[:ephemeral_disk]
          param[:"DataDisk.1.Size"] = ephemeral_disk[:size].to_s
          param[:"DataDisk.1.Category"] = ephemeral_disk[:type]
          #param[:"DataDisk.1.Device"] = "/dev/xvdb"

          if resource_pool.has_key? :instance_name
            param[:InstanceName] = resource_pool[:instance_name]

          end

          if resource_pool.has_key? :availability_zone
            param[:ZoneId] = resource_pool[:availability_zone]
          end

          if resource_pool.has_key? :system_disk
            system_disk = resource_pool[:system_disk]
            param[:"SystemDisk.Size"] = system_disk[:size].to_s
            param[:"SystemDisk.Category"] = system_disk[:type]

          end

        rescue => e
          @logger.error("failed to prepare parameters before creating a new vm. #{e.inspect}\n#{e.backtrace.join('\n')}")
          raise Bosh::Clouds::VMCreationFailed.new(false), "failed to prepare parameters before creating a new vm. #{e.inspect}\n#{e.backtrace.join("\n")}"
        end
      end

      ##
      # Add user data here
      #
      #user_data = '{"Registry": {"Endpoint": "http://'+ @registry_options[:user] + ':' + @registry_options[:password] + '@' + @registry_options[:endpoint] + '"}}'

      user_data = <<EOF
{"Registry":{"Endpoint":"http://#{@registry_options[:user]}:#{@registry_options[:password]}@#{@registry_options[:endpoint].sub(/^http:\/\//, '')}"}}
EOF

      @logger.debug "user data info #{user_data}"

      param[:IoOptimized] = "optimized"
      param[:UserData] = Base64.strict_encode64(user_data)

      @logger.debug "current param is: #{param.inspect}"

      vm_created = false
      begin
        res = @alicloud_client.CreateInstance param
        ins_id = res["InstanceId"]

        while is_vm_pending? ins_id
          sleep 10
        end
        vm_created = true
        @logger.debug "created a vm, the vm id is #{ins_id}"

        if not network.nil?
          if nm.vip?
            # Do not assign, just bind
            bind_public_eip ins_id, nm
            @logger.debug "assigned a public ip for the newly created vm, try to start it"
          end
        end

        #Query disk device for instance
        disk_param = {
            :RegionId => @alicloud_options[:RegionId],
            :InstanceId => ins_id,
            :DiskType => "data"
        }
        disk_res = @alicloud_client.DescribeDisks disk_param
        data_device = disk_res["Disks"]["Disk"][0]["Device"]
        if disk_res["Disks"]["Disk"][0]["Category"] != "cloud"
          data_device[-4] = ""
        end

        @logger.debug "the vm ephemeral data device is #{data_device}"

        # VM will be started after creation
        start_vm ins_id
        @logger.debug "the vm creation is done"

        if not @registry_client.nil?
          registry_settings = initial_agent_settings(
              ins_id,
              agent_id,
              network,
              data_device,
              env
          )
          @logger.debug "registry_settings is #{registry_settings}"
          @registry_client.update_settings(ins_id, registry_settings)
        end

        ins_id
      rescue => e
        delete_vm ins_id if vm_created
        @logger.error("failed to start a new vm. #{e.inspect}\n#{e.backtrace.join('\n')}")
        raise Bosh::Clouds::VMCreationFailed.new(false), "failed to start a new vm. #{e.inspect}\n#{e.backtrace.join("\n")}"
      end

    end

    ##
    # Delete a VM
    #
    # @param vm_id [String] Cloud ID of the VM to delete. Returned from {#create_vm}
    # @return [void]
    #
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
        r = @alicloud_client.DeleteInstance param if vm_stopped
        @logger.debug "the vm is deleted"

        r
      rescue => e
        @logger.error("failed to delete the vm. #{e.inspect}\n#{e.backtrace.join('\n')}")
        raise Bosh::Clouds::VMNotFound.new(false), "failed to delete the vm. #{e.inspect}\n#{e.backtrace.join("\n")}"
      end
    end

    ##
    # Has the vm or not
    #
    # Checks for VM presence in the IaaS
    # @param vm_id [String] Cloud ID of the VM to check. Returned from {#create_vm}
    # @return [Boolean] True if vm is present
    #
    def has_vm? vm_id
      param = {
          :RegionId => @alicloud_options[:RegionId],
          :InstanceIds => "[\"#{vm_id}\"]"
      }
      r = @alicloud_client.DescribeInstances param
      r["Instances"]["Instance"].count != 0
    end

    ##
    # Reboots the vm
    #
    # @param vm_id [String] Cloud ID of the VM to reboot. Returned from {#create_vm}
    # @return [void]
    #
    def reboot_vm vm_id
      @logger.info "start to reboot a vm"
      stop_vm vm_id if is_vm_running? vm_id
      start_vm vm_id if is_vm_stopped? vm_id
      @logger.info "the vm is rebooted"
    end

    ##
    # Creates a disk
    #
    # Creates disk with specific size. Disk does not belong to any given VM
    #
    # @param size [Integer] Size of the disk in MB
    # @param cloud_properties [Hash] Cloud properties hash specified in the deployment manifest under the disk pool
    # @param ins_id [String] Cloud ID of the VM created disk will most likely be attached
    #
    # @return disk_id [String] Cloud ID of the created disk
    def create_disk(size, cloud_properties, ins_id)
      @logger.info("Check the vm belongs to which realiability zone")

      param = {
          :RegionId => @alicloud_options[:RegionId],
          :InstanceIds => "[\"#{ins_id}\"]"
      }

      r = @alicloud_client.DescribeInstances param

      zone_id = "undefined"
      # io_optimized = "undefined"
      if r["TotalCount"] == 1
        zone_id = r["Instances"]["Instance"][0]["ZoneId"]
        # io_optimized = r["Instances"]["Instance"][0]["IoOptimized"]
      else
        raise Bosh::Clouds::VMNotFound.new(false), "failed to get the vm realiability zone id. Response is #{r.inspect}"
      end

      @logger.info("The vm is in realiability zone #{zone_id}")

      @logger.info("Start to create a disk with #{size} GB. it will be band to vm #{ins_id}")

      disk_category = "undefined"
      if r["Instances"]["Instance"][0]["IoOptimized"]
        disk_category = "cloud_efficiency"
      else
        disk_category = "cloud"
      end

      begin
        # TODO size can not be small than the snapshoted vm size
        param = {
            :RegionId => @alicloud_options[:RegionId],
            :ZoneId => zone_id,
            :Size => size.to_s,
            # :DiskSnapShotId => @alicloud_options[:DiskSnapShotId],
            :DiskCategory => disk_category
        }
        @logger.info("Creating a disk with param #{param.inspect}")
        r = @alicloud_client.CreateDisk param
        disk_id = r["DiskId"]

        while disk_status(disk_id) != "Available"
          sleep 10
        end

        disk_id
      rescue => e
        @logger.error("failed to create a disk. #{e.inspect}\n#{e.backtrace.join('\n')}")
        raise Bosh::Clouds::VMNotFound.new(false), "failed to create a disk. #{e.inspect}\n#{e.backtrace.join("\n")}"
      end

      @logger.info("Created a disk #{disk_id}")

      disk_id
    end
    ##
    # Check disk presence
    #
    # Check if a disk is present or not
    #
    # @param disk_id [String] Cloud ID of the disk to check. Returned from {#create_disk}
    #
    # @return [Boolean] True if disk is present
    #
    def has_disk? disk_id
      @logger.info("Check if the disk #{disk_id} exists or not")
      status = disk_status disk_id

      @logger.info("the disk #{disk_id} is #{status}")

      status != "NotExist"
    end

    ##
    # Deletes a disk
    #
    # Deletes a disk. Assume that disk was dettached from all VMs
    #
    # @param disk_id [String] Cloud ID of the disk to delete. Returned from {#create_disk}
    #
    # @return [void]
    #
    def delete_disk disk_id
      if disk_status(disk_id) != "Available"
        @logger.error("The disk #{disk_id} is not available, can not delete directly.")
        raise Bosh::Clouds::DiskNotFound.new(false), "the disk #{disk_id} is not available, can not delete directly"
      end

      @logger.info("Start to delete the disk #{disk_id}")

      param = {
          :DiskId => disk_id
      }
      begin
        r = @alicloud_client.DeleteDisk param
        while has_disk? disk_id
          sleep 10
        end
      rescue => e
        @logger.error("failed to delete the disk. #{e.inspect}\n#{e.backtrace.join('\n')}")
        raise Bosh::Clouds::DiskNotFound.new(false), "failed to delete the disk. #{e.inspect}\n#{e.backtrace.join("\n")}"
      end

      @logger.info("Deleted #{disk_id}")
    end

    ##
    # Attaches disk to the VM
    #
    # @param vm_id [String] Cloud ID of the VM
    # @param disk_id [String] Cloud ID of the disk
    #
    # @return [void]
    #
    def attach_disk ins_id, disk_id
      # TODO if a vm is tagged in LockReason: security, we need to return failed
      if disk_status(disk_id) != "Available"
        @logger.error("The disk #{disk_id} is not available, can not attach")
        raise Bosh::Clouds::DiskNotFound.new(false), "the disk #{disk_id} is not avaiable, can not attach"
      end

      # TODO Currently, we only allow running vm attach a disk
      if vm_status(ins_id) != "Running"
        @logger.error("The vm #{ins_id} can not attach a disk #{disk_id}")
        raise Bosh::Clouds::VMNotFound.new(false), "the vm #{ins_id} can not attach a disk #{disk_id}"
      end

      @logger.info("Start to attach the disk #{disk_id} to the vm #{ins_id}")

      begin
        param = {
            :InstanceId => ins_id,
            :DiskId => disk_id
        }
        @alicloud_client.AttachDisk param

        while disk_status(disk_id) != "In_use"
          sleep 10
        end

        if not @registry_client.nil?
          registry_settings = @registry_client.read_settings(ins_id)
          registry_settings["disks"] ||= {}
          registry_settings["disks"]["persistent"] ||= {}
          registry_settings["disks"]["persistent"][disk_id] = disk_device(disk_id)
          @registry_client.update_settings(ins_id, registry_settings)

          @logger.debug "registry_settings is #{registry_settings}"
        end

      rescue => e
        @logger.error("the vm #{ins_id} failed to attach a disk #{disk_id}. #{e.inspect}\n#{e.backtrace.join('\n')}")
        raise Bosh::Clouds::DiskNotAttached.new(false), "the vm #{ins_id} failed to attach a disk #{disk_id}. #{e.inspect}\n#{e.backtrace.join("\n")}"
      end

      @logger.info("Attached #{disk_id} to #{ins_id}")
    end

    ##
    # Detaches disk from the VM
    #
    # @param vm_id [String] Cloud ID of the VM
    # @param disk_id [String] Cloud ID of the disk
    #
    # @return [void]
    #
    def detach_disk ins_id, disk_id
      if disk_status(disk_id) != "In_use"
        @logger.error("The disk #{disk_id} can not be detached.")
        raise Bosh::Clouds::DiskNotFound.new(false), "the disk #{disk_id} can not be detached"
      end

      # TODO Currently, we only allow running vm detach the disk
      if vm_status(ins_id) != "Running"
        @logger.error("The vm #{ins_id} can not detach the disk #{disk_id}")
        raise Bosh::Clouds::VMNotFound.new(false), "the vm #{ins_id} can not detach the disk #{disk_id}"
      end

      @logger.info("Start to detach the disk #{disk_id} to the vm #{ins_id}")

      begin
        param = {
            :InstanceId => ins_id,
            :DiskId => disk_id
        }
        r = @alicloud_client.DetachDisk param

        if not @registry_client.nil?
          registry_settings = @registry_client.read_settings(ins_id)
          registry_settings["disks"] ||= {}
          registry_settings["disks"]["persistent"] ||= {}
          registry_settings["disks"]["persistent"].delete(disk_id)
          @registry_client.update_settings(ins_id, registry_settings)

          @logger.debug "registry_settings is #{registry_settings}"
        end

        while disk_status(disk_id) != "Available"
          sleep 10
        end
      rescue => e
        @logger.error("the vm #{ins_id} failed to detach the disk #{disk_id}. #{e.inspect}\n#{e.backtrace.join('\n')}")
        raise Bosh::Clouds::DiskNotAttached.new(false), "the vm #{ins_id} failed to detach the disk #{disk_id}. #{e.inspect}\n#{e.backtrace.join("\n")}"
      end

      @logger.info("Detached #{disk_id} from #{ins_id}")
    end

    ##
    # Return list of disks currently attached to the VM
    #
    # @param vm_id [String] Cloud ID of the VM
    #
    # @return disk_ids [Array of strings] Array of disk ids that are currently attached to VM
    #
    def get_disks ins_id
      @logger.info("Check the vm #{ins_id} disks")

      param = {
          :RegionId => @alicloud_options[:RegionId],
          :InstanceId => ins_id,
          :Status => "In_use"
      }
      r = @alicloud_client.DescribeDisks param

      # TODO: We assume that each vm has no more than 10 disks
      # https://help.aliyun.com/document_detail/25514.html
      disk_ids = []
      disk_ids = r["Disks"]["Disk"].map{|d| d["DiskId"]}

      @logger.info("The vm has disks #{disk_ids.inspect}")

      disk_ids
    end

    # Configures networking an existing VM.
    #
    # @param vm_id [String] Cloud ID of the VM to modify. Returned by {#create_vm}
    # @param networks [Hash] Network hashes that specify networks VM must be configured
    #
    # @return [void]
    def configure_networks vm_id, networks=nil
      @logger.info("configure_networks, current vm id is #{vm_id}, current network properties are #{networks})")

      raise Bosh::Clouds::NotSupported
    end

    ##
    # Takes a snapshot of the disk
    #
    # @param disk_id [String] Cloud ID of the disk
    # @param metadata [Hash] Collection of key-value pairs
    #
    # @return snapshot_id [String] Cloud ID of the disk snapshot
    #
    def snapshot_disk disk_id, metadata
    end

    ##
    # Deletes the disk snapshot
    #
    # @param snapshot_id [String] snapshot id to delete
    #
    # @return [void]
    #
    def delete_snapshot
    end

    ##
    # Deprecated method
    # http://bosh.io/docs/cpi-api-v1.html#current_vm_id
    #
    def current_vm_id
    end

    def stop_it vm_id
      stop_vm vm_id
    end

    def initialize_registry

      registry_endpoint   = @registry_options.fetch(:endpoint)
      registry_user       = @registry_options.fetch(:user)
      registry_password   = @registry_options.fetch(:password)

      @registry_client = Bosh::Registry::Client.new(registry_endpoint, registry_user, registry_password)

      begin
        @registry_client.read_settings "check"
      rescue Errno::ECONNREFUSED => e
        @logger.info("failed to read settings from registry endpoint. Continue...")
        @registry_client = nil
      rescue HTTPClient::ConnectTimeoutError => e
        @logger.info("failed to read settings from registry endpoint. Continue...")
        @registry_client = nil
      end


    end

    #
    # Check the disk is attached or not
    # @return [String] if does not exist, return NotExist;
    #                  else return value may be In_use | Available
    #                  | Attaching | Detaching | Creating | ReIniting
    #
    def disk_status disk_id
      @logger.debug("Check the disk status")
      param = {
          :RegionId => @alicloud_options[:RegionId],
          :DiskIds => "[\"#{disk_id}\"]"
      }
      r = @alicloud_client.DescribeDisks param
      if r["TotalCount"] == 1
        @logger.debug("The disk #{disk_id} exists")
        return r["Disks"]["Disk"][0]["Status"]
      else
        @logger.debug("The disk #{disk_id} does not exist")
        return "NotExist"
      end

    end

    #
    # Check the disk device
    # @return [String] if does not exist, return NotExist;
    #                  else return value may be /dev/xvdb | /dev/vdb
    #
    def disk_device disk_id
      @logger.debug("Query the disk device")
      param = {
          :RegionId => @alicloud_options[:RegionId],
          :DiskIds => "[\"#{disk_id}\"]"
      }
      r = @alicloud_client.DescribeDisks param
      if r["TotalCount"] == 1
        @logger.debug("The disk #{disk_id} exists")
        device = r["Disks"]["Disk"][0]["Device"]
        if r["Disks"]["Disk"][0]["Category"] == "cloud"
          return device
        else
          # 如果非普通云盘，需要去除x字母，如: xvdb -> vdb
          device[-4] = ""
          return device
        end

      else
        @logger.debug("The disk #{disk_id} does not exist")
        return "NotExist"
      end

    end

    # Generates initial agent settings. These settings will be read by agent
    # from BOSH registry on a target instance. Disk
    # conventions for alicloud are:
    # system disk: /dev/xvda
    # ephemeral disk: /dev/xvdb
    #
    # @param [String] agent_id Agent id (will be picked up by agent to
    #   assume its identity
    # @param [Hash] network_spec Agent network spec
    # @param [Hash] environment
    #   keys are device type ("ephemeral", "raw_ephemeral") and values are array of strings representing the
    #   path to the block device. It is expected that "ephemeral" has exactly one value.
    # @return [Hash]
    def initial_agent_settings(ins_id, agent_id, network_spec, data_device, environment)
      sys_device = String.new(data_device)
      sys_device[-1] = "a"
      settings = {
          "vm" => {
              "name" => ins_id
          },
          "agent_id" => agent_id,
          "networks" => agent_network_spec,
          "disks" => {
              "system" => sys_device,
              "ephemeral" => data_device,
              "persistent" => {}
          }
      }

      # TODO Will add this two later
      # @param [String] root_device_name root device, e.g. /dev/xvda1
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
          :RegionId => @alicloud_options[:RegionId],
          :Status => "Available",
          :EipAddress => nm.vip_network.ip
      }

      r = @alicloud_client.DescribeEipAddresses param
      @logger.debug "check eip exist or not, #{r.inspect}"

      allocation_id = r['EipAddresses']['EipAddress'][0]['AllocationId']
      param = {
          :InstanceId => vm_id,
          :AllocationId => allocation_id
      }
      r = @alicloud_client.AssociateEipAddress param
      @logger.debug "bond the newly created eip with the vm"
    end

    def start_vm vm_id
      param = {
          :InstanceId => vm_id
      }
      @alicloud_client.StartInstance param if is_vm_stopped? vm_id

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
      r = @alicloud_client.StopInstance param if is_vm_running? vm_id

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
          :RegionId => @alicloud_options[:RegionId],
          :InstanceIds => "[\"#{vm_id}\"]"
      }
      r = @alicloud_client.DescribeInstances param

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
          :RegionId => @alicloud_options[:RegionId],
          :Password => @alicloud_options[:Password]
      }
    end

    def validate_options
      required_keys = [
          :RegionId,
          :AccessKeyId,
          :AccessKey,
          :Password,
      ]

      missing_keys = []
      required_keys.each do |key|
        if !@alicloud_options.has_key?(key)
          missing_keys << "#{key}:"
        end
      end

      raise ArgumentError, "missing configuration parameters > #{missing_keys.join(', ')}" unless missing_keys.empty?
    end

    def validate_resource_pool resource_pool
      missing_keys = []

      if ! resource_pool.has_key? :image_id
        missing_keys << "image_id:"
      end

      if ! resource_pool.has_key? :instance_type
        missing_keys << "instance_type:"
      end

      if ! resource_pool.has_key? :ephemeral_disk
        missing_keys << "ephemeral_disk:"
      else
        required_keys = [
            :type,
            :size
        ]

        required_keys.each do |key|
          if !resource_pool[:ephemeral_disk].has_key? key
            missing_keys << "#{key}:"
          end
        end
      end

      raise ArgumentError, "missing configuration parameters > #{missing_keys.join(', ')}" unless missing_keys.empty?

    end

  end
end
