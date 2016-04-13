require 'spec_helper'

describe Bosh::Aliyun::Client do

  it 'can query aliyun ecs regions' do
    # Init
    o = load_client_options
    o[:service] = :ecs
    logger = Bosh::Clouds::Config.logger
    c = Bosh::Aliyun::Client.new o, logger

    # Get region info
    para = {}
    r = c.DescribeRegions para

    expect(r).to have_key("Regions")
    expect(r["Regions"]).to have_key("Region")
  end

  it 'can create aliyun ecs vm and then delete it' do

    # Init
    o = load_client_options
    o[:service] = :ecs
    logger = Bosh::Clouds::Config.logger
    c = Bosh::Aliyun::Client.new o, logger

    # Create an instance
    para = {
      :RegionId => "cn-hangzhou",
      :InstanceType => "ecs.t1.small",
      :ImageId => "m-23qx965sh",
      :SecurityGroupId => "sg-237p56jii",
      :InternetChargeType => "PayByTraffic",
      :InternetMaxBandwidthOut => "10",
      :InstanceName => "bosh_aliyun_cpi_test",
      :Description => "",
      :HostName => "",
      :Password => "1qaz@WSX"
    }
    r = c.CreateInstance para

    expect(r).to have_key("InstanceId")
    expect(r["InstanceId"]).to match(/[\w]{1}-[\w]{9}/)

    sleep 15

    # Delete an instance
    para = {
      :InstanceId => r["InstanceId"]
    }
    r = c.DeleteInstance para

    expect(r).to have_key("RequestId")
  end

  it 'can describe aliyun ecs vm status' do

    # Init
    o = load_client_options
    o[:service] = :ecs
    logger = Bosh::Clouds::Config.logger
    c = Bosh::Aliyun::Client.new o, logger

    # Get instance status
    para = {
      :RegionId => "cn-hangzhou"
    }
    r = c.DescribeInstanceStatus para

    expect(r).to have_key("InstanceStatuses")

  end

end
