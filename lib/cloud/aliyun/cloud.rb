module Bosh::Aliyun
  class Cloud < Bosh::Cloud

    attr_reader :options

    def initialize(options)

      @options = options.dup.freeze

      @aliyun_client = Bosh::Aliyun::Client.new @options[:aliyun]
      @region_id = @options[:aliyun][:region_id]
    end

    def create_stemcell(image_path = nil, cloud_properties = nil)

      # TODO
      # Upload image from local path

      current_stemcell_id
    end

    def delete_stemcell(stemcell_id)
      nil
    end

    def create_vm(agent_id=nil, stemcell_id=nil, resource_pool=nil, networks=nil, disk_locality=nil, env=nil)

      param = {}
      param[:ImageId] = stemcell_id || current_stemcell_id
      param.merge! prepare_create_vm_parameters

      begin
        vm_created = true
        res = @aliyun_client.CreateInstance param

        ins_id = res["InstanceId"]

        vm_started = true
        start_vm ins_id

        ins_id
      rescue => e
        stop_vm ins_id if vm_started
        delete_vm ins_id if vm_created
        raise AliyunException.new "Failed to start a new vm. #{e.inspect}\n#{e.backtrace.join("\n")}"
      end

    end

    def delete_vm(vm_id)

      param={
        :InstanceId => vm_id
      }
      @aliyun_client.DeleteInstance param if is_vm_stopped? vm_id
    end

    def has_vm?(vm_id)
      param = {
        :RegionId => @region_id,
        :InstanceIds => "[\"#{vm_id}\"]"
      }

      r = @aliyun_client.DescribeInstances param

      r["Instances"]["Instance"].count != 0
    end

    def reboot_vm(vm_id)

      stop_vm vm_id if is_vm_running? vm_id
      start_vm vm_id if is_vm_stopped? vm_id
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
      puts "Starting"
      until is_vm_running? vm_id
        puts "DOWN. Ping #{count} time"
        count += 1
        sleep 10
      end
      puts "Started"
    end

    def stop_vm vm_id
      param = {
        :InstanceId => vm_id
      }
      r = @aliyun_client.StopInstance param if is_vm_running? vm_id

      count = 1
      puts "Stopping"
      until is_vm_stopped? vm_id
        puts "UP. Ping #{count} time"
        count += 1
        sleep 10
      end
      puts "Stopped"
    end

    def vm_status vm_id

      param = {
        :RegionId => @region_id,
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

    def current_stemcell_id
      "m-23g9tihvk"
    end

    def prepare_create_vm_parameters
      para = {
        :RegionId => @region_id,
        :InstanceType => "ecs.t1.small",
        :SecurityGroupId => "sg-237p56jii",
        :InternetChargeType => "PayByTraffic",
        :InternetMaxBandwidthOut => "10",
        :InstanceName => "bosh_aliyun_cpi_test",
        :Description => "",
        :HostName => "",
        :Password => "1qaz@WSX"
      }
    end
  end
end
