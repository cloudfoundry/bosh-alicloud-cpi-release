module Bosh::Aliyun
  class AliyunSnapshotWrapper
    #必传参数：
    def AliyunImgWrapper.createSnapshot(parameters)
      #parameter check:
      parameters["Action"]= "CreateSnapshot";
      parameters["Version"]= "2014-05-26";
      AliyunOpenApiHttpsUtil.request(parameters);
    end

    #必传参数：
    def AliyunImgWrapper.deleteSnapshot(parameters)
      #parameter check:
      parameters["Action"]= "DeleteSnapshot";
      parameters["Version"]= "2014-05-26";
      AliyunOpenApiHttpsUtil.request(parameters);
    end
  end
end