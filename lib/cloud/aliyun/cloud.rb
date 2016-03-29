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

    def create_vm(agent_id=nil, stemcell_id=nil, resource_pool=nil, networks=nil, disk_locality=nil, env=nil)
      aliyun_properties = options.fetch('aliyun');
      parameters = getCreatVmParameter(aliyun_properties);
      body = AliyunInstanceWrapper.createInstance(parameters);

      #磁盘挂载
      #外网设置
      #实例启动
      parameters={};
      parameters["InstanceId"] = body["InstanceId"];
      initCommonParameter(aliyun_properties, parameters)
      AliyunInstanceWrapper.startInstance(parameters)
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
    body = AliyunInstanceWrapper.createInstance(parameters);

    def snapshot_disk(disk_id, metadata)
      parameters={};
      AliyunImgWrapper.createSnapshot(parameters);
      parameters["InstanceId"] = body["InstanceId"];
      initCommonParameter(aliyun_properties, parameters)
      AliyunInstanceWrapper.startInstance(parameters)
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

    def getCreatVmParameter(aliyun_properties)
      parameters={};
      #regionID：cn-hangzhou
      #实例类型：ecs.s3.large
      #imageID：m-23g9tihvk
      #安全组ID：sg-237p56jii
      #计费类型：PayByTraffic
      #公网入带宽：10M
      keys=["RegionId", "ImageId", "InstanceType", "SecurityGroupId", "InternetChargeType", "InternetMaxBandwidthOut", "Password","DataDisk.1.Size"];
      keys.each { |key|
        if aliyun_properties.has_key?(key)
          parameters[key]=aliyun_properties[key];
        end
      }
      initCommonParameter(aliyun_properties, parameters);
      return parameters;
    end

    def initCommonParameter(aliyun_properties, parameters)
      #AccessKeyId:***REMOVED***
      #AccessKey:***REMOVED***
      parameters["AccessKeyId"]=aliyun_properties["AccessKeyId"];
      parameters["Secret"]=aliyun_properties["Secret"];
    end
  end
end
