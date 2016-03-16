module Bosh::Aliyun
  class AliyunDiskWrapper
    #必传参数：
    def AliyunImgWrapper.createDisk(parameters)
      #parameter check:
      parameters["Action"]= "CreateDisk";
      parameters["Version"]= "2014-05-26";
      AliyunOpenApiHttpsUtil.request(parameters);
    end

    #必传参数：
    def AliyunImgWrapper.describeDisks(parameters)
      #parameter check:
      parameters["Action"]= "DescribeDisks";
      parameters["Version"]= "2014-05-26";
      AliyunOpenApiHttpsUtil.request(parameters);
    end

    #必传参数：
    def AliyunImgWrapper.deleteDisk(parameters)
      #parameter check:
      parameters["Action"]= "DeleteDisk";
      parameters["Version"]= "2014-05-26";
      AliyunOpenApiHttpsUtil.request(parameters);
    end

    #必传参数：
    def AliyunImgWrapper.attachDisk(parameters)
      #parameter check:
      parameters["Action"]= "AttachDisk";
      parameters["Version"]= "2014-05-26";
      AliyunOpenApiHttpsUtil.request(parameters);
    end

    #必传参数：
    def AliyunImgWrapper.detachDisk(parameters)
      #parameter check:
      parameters["Action"]= "DetachDisk";
      parameters["Version"]= "2014-05-26";
      AliyunOpenApiHttpsUtil.request(parameters);
    end

  end
end