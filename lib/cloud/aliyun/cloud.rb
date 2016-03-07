module Bosh::Aliyun
  class Cloud < Bosh::Cloud

    attr_reader :options

    def initialize(options)
      @options = options.dup.freeze

      @stemcell_manager = Bosh::Aliyun::StemcellManager.new(aliyun_properties)
    end

    def create_stemcell(image_path, cloud_properties)
      with_thread_name("create_stemcell(#{image_path}...)") do
        @stemcell_manager.create_stemcell(image_path, cloud_properties)
      end
    end

    def delete_stemcell(stemcell_id)
      with_thread_name("delete_stemcell(#{stemcell_id})") do
        @stemcell_manager.delete_stemcell(stemcell_id)
      end
    end

    private

    def aliyun_properties
      @aliyun_properties ||= options.fetch('aliyun')
    end

  end
end
