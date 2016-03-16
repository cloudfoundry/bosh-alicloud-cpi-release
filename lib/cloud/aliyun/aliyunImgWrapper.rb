require 'json'

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

    #必传参数：
    def AliyunImgWrapper.describeImages(parameters)
      #parameter check:
      parameters["Action"]= "DescribeImages";
      parameters["Version"]= "2014-05-26";
      AliyunOpenApiHttpsUtil.request(parameters);
    end

    #必传参数：
    def AliyunImgWrapper.hasImg(parameters)
      #parameter check:
      flag=false;
      imgs = AliyunImgWrapper.describeImages(parameters);
      if imgs.fetch("PageSize")==1
        flag=true;
      end
      return flag;
    end

  end
end