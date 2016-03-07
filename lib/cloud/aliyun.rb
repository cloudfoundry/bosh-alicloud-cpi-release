module Bosh
  module Aliyun; end
end

require "common/exec"
require "common/thread_pool"
require "common/thread_formatter"
require "common/common"

require "cloud"
require "cloud/aliyun/cloud"
require "cloud/aliyun/stemcell_manager"

module Bosh
  module Clouds
    Aliyun = Bosh::Aliyun::Cloud
  end
end
