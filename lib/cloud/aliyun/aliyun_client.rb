require "base64"
require "uri"
require "json"

module Bosh::Aliyun
  class APIBase
    def self.info
      raise "Service Name Missing."
    end

    def self.endpoint
      raise "Service Endpoint Missing."
    end

    def self.default_parameters
      raise "Service Default Parameters Missing."
    end

    def self.http_method
      "GET"
    end
  end

  class ECS < APIBase
    def self.info
      "Aliyun ECS Service"
    end

    def self.endpoint
      "https://ecs.aliyuncs.com/"
    end

    def self.default_parameters
      {
        :Format => "JSON",
        :Version => "2014-05-26",
        :SignatureMethod => "HMAC-SHA1",
        :SignatureVersion => "1.0"
      }
    end

    def self.separator
      "&"
    end

    def self.http_method
      super
    end
  end

  SERVICES = {
    :ecs => ECS
  }

  class Client

    attr_accessor :access_key_id, :access_key
    attr_accessor :options
    attr_accessor :endpoint
    attr_accessor :service

    def initialize(options = {}, logger = nil)
      validate_options options
      options[:service] ||= :ecs

      @@logger = logger || Logger.new(STDERR)
      @@service = SERVICES[options[:service].to_sym]
      @@access_key_id = options[:AccessKeyId]
      @@access_key = options[:AccessKey]
      @@endpoint = options[:endpoint] || @@service.endpoint
      @@options = {:AccessKeyId => @@access_key_id}
    end

    def method_missing(method, *args)
      if args[0].nil?
        raise AliyunException.new "Method missing: #{method}."
      end

      request(method, args[0])
    end

    private
    def request(method, params)
      params = prepare_parameters(method, params)
      uri = URI(@@endpoint)

      # Prepare params if we need to POST
      uri.query = URI.encode_www_form(params)

      http = Net::HTTP.new(uri.host, uri.port)

      # Use SSL for better security
      http.use_ssl = (uri.scheme == "https")

      # Ignore verify
      http.verify_mode = OpenSSL::SSL::VERIFY_NONE

      @@logger.debug "request params is: #{params.inspect}, request endpoint is: #{@@endpoint}"
      request = Net::HTTP::Get.new(uri.request_uri)
      response = http.request(request)

      case response
      when Net::HTTPSuccess
        r = JSON.parse(response.body)
        @@logger.debug "response body is: #{r.inspect}"
        return r
      else
        @@logger.error "request error! reponse code: #{response.code}, message: #{response.body.inspect}"
        raise AliyunException.new "request error! response code: #{response.code}, message: #{response.body}"
      end
    end

    def prepare_parameters(method, params)
      #Add common parameters
      params.merge! @@service.default_parameters

      params.merge! @@options

      params[:Action] = method.to_s
      params[:Timestamp] = Time.now.utc.iso8601
      params[:SignatureNonce] = SecureRandom.uuid
      params[:Signature] = compute_signature(params)

      params
    end

    def compute_signature(params)

      sorted_keys = params.keys.sort

      capitalized_string = ""
      capitalized_string = sorted_keys.map {|key|
        "%s=%s" % [safe_encode(key.to_s), safe_encode(params[key])]
      }.join(@@service.separator)

      length = capitalized_string.length

      string_to_sign = @@service.http_method + @@service.separator + safe_encode("/") + @@service.separator + safe_encode(capitalized_string)

      signature = Base64.strict_encode64(OpenSSL::HMAC.digest(OpenSSL::Digest::SHA1.new, @@access_key+"&", string_to_sign))

      signature
    end

    def safe_encode(value)
      CGI.escape(value).gsub("+", "%20").gsub("*", "%2A").gsub("%7E", "~")
    end

    def validate_options options
      required_keys = [
          :AccessKeyId,
          :AccessKey
      ]

      missing_keys = []
      required_keys.each do |key|
        if !options.has_key?(key)
          missing_keys << "#{key}:"
        end
      end

      @@logger.error "missing configuration parameters > #{missing_keys.join(', ')}" unless missing_keys.empty?
      raise ArgumentError, "missing configuration parameters > #{missing_keys.join(', ')}" unless missing_keys.empty?
    end
  end

end
