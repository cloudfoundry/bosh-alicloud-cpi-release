module Bosh::Aliyun
  class AliyunImgWrapper
    #必传参数：
    def AliyunImgWrapper.createImage(parameters)
      #parameter check:
      parameters["Action"]= "CreateImage";
      parameters["Version"]= "2014-05-26";
      AliyunOpenApiHttpsUtil.request(parameters);
    end

    #必传参数：
    def AliyunImgWrapper.deleteImage(parameters)
      #parameter check:
      parameters["Action"]= "DeleteImage";
      parameters["Version"]= "2014-05-26";
      AliyunOpenApiHttpsUtil.request(parameters);
    end
  end
end