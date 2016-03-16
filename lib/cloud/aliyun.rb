module Bosh
  module Aliyun; end
end

require "common/exec"
require "common/thread_pool"
require "common/thread_formatter"
require "common/common"

require "cloud"
require "cloud/aliyun/cloud"
require "cloud/aliyun/aliyunConstants"
require "cloud/aliyun/aliyunHttpsUtil"
require "cloud/aliyun/aliyunDiskWrapper"
require "cloud/aliyun/aliyunImgWrapper"
require "cloud/aliyun/aliyunInstanceWrapper"
require "cloud/aliyun/aliyunSecurityGroupWrapper"
require "cloud/aliyun/aliyunSnapshotWrapper"

module Bosh
  module Clouds
    Aliyun = Bosh::Aliyun::Cloud
  end
end
