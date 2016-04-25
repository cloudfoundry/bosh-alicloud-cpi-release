module Bosh::Aliyun

  class NetworkManager
    include Helpers

    attr_reader :vip_network, :network
    attr_accessor :logger

    #  {"network_a" =>
    #    {
    #      "netmask"          => "255.255.248.0",
    #      "ip"               => "172.30.41.40",
    #      "gateway"          => "172.30.40.1",
    #      "dns"              => ["172.30.22.153", "172.30.22.154"],
    #      "cloud_properties" => {"name" => "VLAN444"}
    #    }
    #  }

    def initialize spec
      unless spec.is_a? Hash
        raise ArgumentError, "invalid spec, hash expected"
      end

      spec = recursive_symbolize_keys(spec)

      @logger = Bosh::Clouds::Config.logger

      @logger.debug("networks: #{spec}")

      @vip_network = nil
      @network = nil
      spec.each_pair do |name, network_spec|
        network_type = network_spec[:type] || "manual"

        case network_type
        when "manual"
          # @logger.error "exactly one dynamic or manual network per instance is required" if @network
          # TODO raise
          @network = ManualNetwork.new name, network_spec
        when "vip"
          # @logger.error "more than one vip network for '#{name}'" if @vip_network
          # TODO raise
          @vip_network = VipNetwork.new name, network_spec
        when "dynamic"
          @logger.error "exactly one dynamic or manual network per instance is required" if @network
          # Not support right now
          # TODO raise
          # @network = DynamicNetwork.new name, network_spec
        else
          @logger.error "invalid network type for #{network_type}, can only handle 'dynamic', 'vip' or 'manual' network types"
          # TODO raise
        end
      end
    end

    def subnet_name
      @network.subnet
    end

    def private_ip
      vpc? ? @network.private_ip : nil
    end

    def vpc?
      @network.is_a? ManualNetwork
    end

    def vip?
      @vip_network.is_a? VipNetwork
    end
  end


  class Network

    def initialize name, spec
      unless spec.is_a? Hash
        raise ArgumentError, "Invalid spec, Hash expected"
      end

      @logger = Bosh::Clouds::Config.logger

      @name = name
      @ip = spec[:ip]
      @cloud_properties = spec[:cloud_properties]

      configure
    end

    def private_ip
      @ip
    end

    def configure
      @logger.error "`configure` is not implemented by #{self.class}"
    end
  end

  class DynamicNetwork < Network

    include Helpers

    def initialize name, spec
      super
    end

    def configure
      {
        :SecurityGroupId => @cloud_properties[:SecurityGroupId]
      }
    end
  end

  class ManualNetwork < Network
    include Helpers

    def initialize name, spec
      super
      if @cloud_properties.nil? || !@cloud_properties.has_key?(:VSwitchId) || !@cloud_properties.has_key?(:SecurityGroupId)
        @logger.error "subnet is required for manual network"
      end
    end

    def configure
      {
        :SecurityGroupId => @cloud_properties[:SecurityGroupId],
        :VSwitchId => @cloud_properties[:VSwitchId],
        :PrivateIpAddress => @ip
      }
    end
  end

  class VipNetwork < Network
    def initialize name, spec
      super
    end

    def configure
      {
        :Bandwidth => @cloud_properties[:InternetMaxBandwidthOut],
        :InternetChargeType => @cloud_properties[:InternetChargeType]
      }
    end
  end
end
