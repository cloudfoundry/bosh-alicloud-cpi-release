module Bosh
  module Aliyun

    class AliyunException < RuntimeError
    end

  end

end

require "httpclient"
require "yajl"

require "common/exec"
require "common/thread_pool"
require "common/thread_formatter"
require "common/common"

require "bosh/registry/client"

require "cloud"
require "cloud/aliyun/helpers"
require "cloud/aliyun/cloud"
require "cloud/aliyun/network_manager"
require "cloud/aliyun/aliyun_client"

module Bosh
  module Clouds
    Aliyun = Bosh::Aliyun::Cloud
  end
end
