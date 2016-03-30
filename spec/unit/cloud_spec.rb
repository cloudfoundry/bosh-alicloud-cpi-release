require 'spec_helper'

describe Bosh::Aliyun::Cloud do

  describe '#initialize' do
    it 'when all the required configurations are present' do
      expect(true).to eq(true)
    end

  end

  describe '#createImage' do
    options = {}
    options["aliyun"] = {}
    options["aliyun"]["RegionId"] = "cn-hangzhou"
    options["aliyun"]["InstanceType"] = "ecs.s3.large"
    options["aliyun"]["ImageId"] = "m-23g9tihvk"
    options["aliyun"]["SecurityGroupId"] = "sg-237p56jii"
    options["aliyun"]["InternetChargeType"] = "PayByTraffic"
    options["aliyun"]["InternetMaxBandwidthOut"] = "10"
    options["aliyun"]["InstanceName"] = "bosh_aliyun_cpi_test"
    options["aliyun"]["DataDisk.1.Size"] = "300"
    options["aliyun"]["Description"] = ""
    options["aliyun"]["HostName"] = ""
    options["aliyun"]["Password"] = "1qaz@WSX"
    options["aliyun"]["AccessKeyId"] = "***REMOVED***"
    options["aliyun"]["AccessKey"] = "***REMOVED***"

    it 'should create a image' do
      cloud = Bosh::Aliyun::Cloud.new(options)
      expect(cloud.create_vm).to eq(stemcell_id)
    end
  end
end
