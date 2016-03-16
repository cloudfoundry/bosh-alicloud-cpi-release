module Bosh::Aliyun
  class Cloud < Bosh::Cloud

    attr_reader :options

    def initialize(options)
      @options = options.dup.freeze
    end

    def create_stemcell(image_path, cloud_properties)
      parameters={};
      AliyunImgWrapper.createImage(parameters);
    end

    def delete_stemcell(stemcell_id)
      parameters={};
      AliyunImgWrapper.deleteImage(parameters);
    end

    def create_vm(agent_id, stemcell_id, resource_pool, networks, disk_locality, env)
      aliyun_properties = options.fetch('aliyun');

      AliyunInstanceWrapper.createInstance(parameters);
    end

    def delete_vm(vm_id)
      parameters={};
      AliyunInstanceWrapper.deleteInstance(parameters);
    end

    def has_vm?(vm_id)
      parameters={};
      AliyunInstanceWrapper.describeInstances(parameters);
    end

    def has_disk?(disk_id)
      parameters={};
      AliyunImgWrapper.describeDisks(parameters);
    end

    def reboot_vm(vm_id)
      parameters={};
      AliyunInstanceWrapper.rebootInstance(parameters);
    end

    def set_vm_metadata(vm, metadata)
      parameters={};
      AliyunInstanceWrapper.modifyInstanceAttribute(parameters);
    end

    def create_disk(size, cloud_properties, vm_locality)
      parameters={};
      AliyunImgWrapper.createDisk(parameters);
    end

    def delete_disk(disk_id)
      parameters={};
      AliyunImgWrapper.deleteDisk(parameters);
    end

    def attach_disk(vm_id, disk_id)
      parameters={};
      AliyunImgWrapper.attachDisk(parameters);
    end

    def snapshot_disk(disk_id, metadata)
      parameters={};
      AliyunImgWrapper.createSnapshot(parameters);
    end

    def delete_snapshot(snapshot_id)
      parameters={};
      AliyunImgWrapper.deleteSnapshot(parameters);
    end

    def detach_disk(vm_id, disk_id)
      parameters={};
      AliyunImgWrapper.detachDisk(parameters);
    end

    def get_disks(vm_id)
      parameters={};
      AliyunImgWrapper.describeDisks(parameters);
    end

    private

  end
end
