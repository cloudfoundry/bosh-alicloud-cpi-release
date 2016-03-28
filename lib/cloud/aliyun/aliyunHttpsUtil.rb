require 'cgi'
require 'hmac-sha1'
require 'digest/sha1'
require 'base64'
require 'uri'
require 'net/http'
require 'time'
require 'json'
module Bosh::Aliyun
  class AliyunOpenApiHttpsUtil
    @@Aliyun_OpenApi_Url="http://ecs.aliyuncs.com/?";
    @@Secret_Key="Secret";
    @@Signature="Signature";

    def AliyunOpenApiHttpsUtil.request(parameters)
      initCommonParameters(parameters);
      uri = URI.parse(@@Aliyun_OpenApi_Url);
      response = Net::HTTP.post_form(uri, parameters);
      body=JSON.parse(response.body);
      #if body.has_key("Code")
      #  raise body["Code"]+":"+body["Message"];
      #end
      puts body;
      return body;
    end

    def AliyunOpenApiHttpsUtil.initCommonParameters(parameters)
      parameters["Format"]= "JSON";
      parameters["SignatureMethod"]= "HMAC-SHA1";
      parameters["Timestamp"]= Time.now.utc.iso8601.to_s;
      parameters["SignatureVersion"]= "1.0";
      parameters["SignatureNonce"]= rand(9999999999999999).to_s;
      getSignature(parameters);
    end

    def AliyunOpenApiHttpsUtil.getSignature(parameters)
      secretKey=parameters[@@Secret_Key];

      parameters.delete(@@Secret_Key);

      keys = parameters.keys;
      keys = keys.sort;
      start=true;
      query="";

      keys.each{ |item|
        if(start)
          start=false;
        else
          query<<"&";
        end
        query<<item<<"="<<percentEncode(parameters[item]);
      }

      stringToSign = "POST" + "&" + percentEncode("/") + "&" + percentEncode(query);
      parameters[@@Signature]=Base64.encode64(HMAC::SHA1.digest(secretKey, stringToSign)).strip;
    end

    def AliyunOpenApiHttpsUtil.percentEncode(str)
      flag = CGI.escape(str).gsub("+", "%20").gsub("*", "%2A").gsub("%7E", "~");
      return flag;
    end

  end
end
