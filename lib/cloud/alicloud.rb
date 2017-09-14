module Bosh
  module Alicloud

    class AlicloudException < RuntimeError
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
require "cloud/alicloud/helpers"
require "cloud/alicloud/cloud"
require "cloud/alicloud/network_manager"
require "cloud/alicloud/alicloud_client"

module Bosh
  module Clouds
    Alicloud = Bosh::Alicloud::Cloud
  end
end
