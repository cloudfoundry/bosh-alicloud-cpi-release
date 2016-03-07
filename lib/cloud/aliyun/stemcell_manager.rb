module Bosh::Aliyun
  class StemcellManager

    FAKE_STEMCELL_NAME = "test_stemcell_name"

    def initialize(aliyun_properties)
      @logger = Bosh::Clouds::Config.logger
    end

    def create_stemcell(image_path, cloud_properties)
      @logger.info("create_stemcell(#{image_path})")      
      FAKE_STEMCELL_NAME
    end

    def delete_stemcell(name)
      @logger.info("delete_stemcell(#{name})")
    end
  end
end
