module Bosh::Aliyun
  class AliyunDiskWrapper
    #必传参数：
    def self.createDisk(parameters)
      #parameter check:
      parameters["Action"]= "CreateDisk";
      parameters["Version"]= "2014-05-26";
      AliyunOpenApiHttpsUtil.request(parameters);
    end

    #必传参数：
    def self.describeDisks(parameters)
      #parameter check:
      parameters["Action"]= "DescribeDisks";
      parameters["Version"]= "2014-05-26";
      AliyunOpenApiHttpsUtil.request(parameters);
    end

    #必传参数：
    def self.deleteDisk(parameters)
      #parameter check:
      parameters["Action"]= "DeleteDisk";
      parameters["Version"]= "2014-05-26";
      AliyunOpenApiHttpsUtil.request(parameters);
    end

    #必传参数：
    def self.attachDisk(parameters)
      #parameter check:
      parameters["Action"]= "AttachDisk";
      parameters["Version"]= "2014-05-26";
      AliyunOpenApiHttpsUtil.request(parameters);
    end

    #必传参数：
    def self.detachDisk(parameters)
      #parameter check:
      parameters["Action"]= "DetachDisk";
      parameters["Version"]= "2014-05-26";
      AliyunOpenApiHttpsUtil.request(parameters);
    end

  end
end
