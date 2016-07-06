module Bosh::Aliyun
  class DiskManager
    
    include Bosh::Exec

    def initialize(aliyun_properties)
      @aliyun_properties = aliyun_properties

      @logger = Bosh::Clouds::Config.logger
    end

    # Create a disk lazily, it could be attached to a vm later.
    #
    # @param [Integer] size: Disk size in GB
    # @param [Hash] cloud_properties: Cloud properties to create the disk
    # @return [String] disk_name: a unique ID to name a disk
    def create_disk(size, cloud_properties)
      @logger.info("create_disk(#{size}, #{cloud_properties})")
      "test-disk-name"
    end

    def has_disk(disk_name)
      @logger.info("has_disk(#{disk_name})")q
    end

    def delete_disk(disk_name)
      @logger.info("delete_disk(#{disk_name})")
    end


