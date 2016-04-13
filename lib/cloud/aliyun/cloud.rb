module Bosh::Aliyun

  class Cloud < Bosh::Cloud

    include Helpers

    attr_reader :options

    def initialize(options)
      @options = recursive_symbolize_keys(options["aliyun"])
      validate_options

      @logger = Bosh::Clouds::Config.logger

      @logger.debug "current options is #{options.inspect}"

      @aliyun_client = Bosh::Aliyun::Client.new(@options, @logger)
    end

    def create_stemcell(image_path = nil, cloud_properties = nil)

      @logger.info("fake stemcell creation for now")
      # TODO
      # Upload image from local path

      @options[:ImageId]
    end

    def delete_stemcell(stemcell_id)
      @logger.info("fake stemcell deletion for now")
      nil
    end

    def create_vm(agent_id=nil, stemcell_id=nil, resource_pool=nil, networks=nil, disk_locality=nil, env=nil)

      @logger.info("start to create a vm")
      param = {}
      param[:ImageId] = stemcell_id || @options[:ImageId]
      param.merge! prepare_create_vm_parameters
      @logger.debug("current param is: #{param.inspect}")

      begin
        vm_created = true
        res = @aliyun_client.CreateInstance param
        ins_id = res["InstanceId"]
        @logger.debug("created a vm, the vm id is #{ins_id}. Try to start it")

        # VM will be started after creation
        vm_started = true
        start_vm ins_id
        @logger.debug("the vm creation is done")

        ins_id
      rescue => e
        stop_vm ins_id if ! ins_id.nil? && vm_started
        delete_vm ins_id if ! ins_id.nil? && vm_created
        @logger.error(%Q[failed to start a new vm. #{e.inspect}\n#{e.backtrace.join("\n")}])
        raise Bosh::Clouds::VMCreationFailed.new(false), "failed to start a new vm. #{e.inspect}\n#{e.backtrace.join("\n")}"
      end

    end

    def delete_vm(vm_id)

      @logger.info("start to delete a vm")
      param={
        :InstanceId => vm_id
      }
      r = @aliyun_client.DeleteInstance param if is_vm_stopped? vm_id
      @logger.info("the vm is deleted")

      r
    end

    def has_vm?(vm_id)

      param = {
        :RegionId => @options[:RegionId],
        :InstanceIds => "[\"#{vm_id}\"]"
      }
      r = @aliyun_client.DescribeInstances param

      r["Instances"]["Instance"].count != 0
    end

    def reboot_vm(vm_id)
      @logger.info "start to reboot a vm"
      stop_vm vm_id if is_vm_running? vm_id
      start_vm vm_id if is_vm_stopped? vm_id
      @logger.info "the vm is rebooted"
    end

    def stop_it vm_id
      stop_vm vm_id
    end

    private

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

    def is_vm_running? vm_id
      vm_status(vm_id) == "Running"
    end

    def is_vm_stopped? vm_id
      vm_status(vm_id) == "Stopped"
    end

    def prepare_create_vm_parameters
      para = {
        :RegionId => @options[:RegionId],
        :InstanceType => @options[:InstanceType],
        :SecurityGroupId => @options[:SecurityGroupId],
        :InternetChargeType => @options[:InternetChargeType],
        :InternetMaxBandwidthOut => @options[:InternetMaxBandwidthOut],
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
          :SecurityGroupId,
          :AccessKeyId,
          :AccessKey,
          :Password,

          :InternetChargeType,
          :InternetMaxBandwidthOut,
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
