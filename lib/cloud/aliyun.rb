module Bosh
  module Aliyun

    class AliyunException < RuntimeError
    end

  end

end

require "common/exec"
require "common/thread_pool"
require "common/thread_formatter"
require "common/common"

require "cloud"
require "cloud/aliyun/helpers"
require "cloud/aliyun/cloud"
require "cloud/aliyun/aliyun_client"

module Bosh
  module Clouds
    Aliyun = Bosh::Aliyun::Cloud
  end
end
