module Bosh::Aliyun
  class AliyunInstanceWrapper
    #必传参数：地域(RegionId)、镜像文件 ID()、实例的资源规则、安全组代码
    def AliyunInstanceWrapper.createInstance(parameters)
      #parameter check:
      parameters["Action"]= "CreateInstance";
      parameters["Version"]= "2014-05-26";
      AliyunOpenApiHttpsUtil.request(parameters);
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

  end
end