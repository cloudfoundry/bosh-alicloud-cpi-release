module Bosh::Aliyun
  class AliyunInstanceWrapper
    #必传参数：地域(RegionId)、镜像文件 ID()、实例的资源规则、安全组代码
    def AliyunInstanceWrapper.createInstance(parameters)
      #parameter check:
      #if !has_img(aliyun_properties)
      #  raise "image not exist";
      #end

      #if !has_securityGroup(aliyun_properties)
      #  raise "securityGroup not exist";
      #end

      parameters["Action"]= "CreateInstance";
      parameters["Version"]= "2014-05-26";
      return AliyunOpenApiHttpsUtil.request(parameters);
    end

    #必传参数：实例 ID(InstanceId)
    def AliyunInstanceWrapper.deleteInstance(parameters)
      #parameter check:
      parameters["Action"]= "DeleteInstance";
      parameters["Version"]= "2014-05-26";
      AliyunOpenApiHttpsUtil.request(parameters);
    end

    #必传参数：
    def AliyunInstanceWrapper.describeInstances(parameters)
      #parameter check:
      parameters["Action"]= "DescribeInstances";
      parameters["Version"]= "2014-05-26";
      AliyunOpenApiHttpsUtil.request(parameters);
    end

    #必传参数：
    def AliyunInstanceWrapper.rebootInstance(parameters)
      #parameter check:
      parameters["Action"]= "RebootInstance";
      parameters["Version"]= "2014-05-26";
      AliyunOpenApiHttpsUtil.request(parameters);
    end

    #必传参数：
    def AliyunInstanceWrapper.modifyInstanceAttribute(parameters)
      #parameter check:
      parameters["Action"]= "ModifyInstanceAttribute";
      parameters["Version"]= "2014-05-26";
      AliyunOpenApiHttpsUtil.request(parameters);
    end

    #必传参数：
    def AliyunInstanceWrapper.startInstance(parameters)
      #parameter check:
      parameters["Action"]= "StartInstance";
      parameters["Version"]= "2014-05-26";
      AliyunOpenApiHttpsUtil.request(parameters);
    end

    #必传参数：
    def AliyunInstanceWrapper.allocatePublicIpAddress(parameters)
      #parameter check:
      parameters["Action"]= "AllocatePublicIpAddress";
      parameters["Version"]= "2014-05-26";
      AliyunOpenApiHttpsUtil.request(parameters);
    end

    private

    def AliyunInstanceWrapper.has_img(paramInput)
      parameters={};
      parameters["RegionId"] = paramInput.fetch("RegionId");
      parameters["ImageId"] = paramInput.fetch("ImageId");
      parameters["Status"] = "Available";
      return AliyunImgWrapper.hasImg(parameters);
    end

    def AliyunInstanceWrapper.has_securityGroup(paramInput)
      parameters={};
      parameters["RegionId"] = paramInput.fetch("RegionId");
      parameters["SecurityGroupId"] = paramInput.fetch("SecurityGroupId");
      return AliyunSecurityGroupWrapper.hasSecurityGroup(parameters);
    end

  end
end