require "spec_helper"

describe Bosh::Aliyun::NetworkManager do

  let(:manual) {
    {
      :type => "manual",
      :ip => "127.0.0.1",
      :cloud_properties => {
        :VSwitchId => "vsw-23priv72n",
        :SecurityGroupId => "sg-238ux30qe"
      }
    }
  }

  let(:vip) {
    {
      :type => "vip",
      :cloud_properties => {
        :InternetMaxBandwidthOut => "10",
        :InternetChargeType => "PayByBandwidth"
      }
    }
  }

  it "should raise an error if the spec is not a hash" do
    expect {
      Bosh::Aliyun::NetworkManager.new :foo
    }.to raise_error ArgumentError
  end

  describe "Manual network" do
    it "should return a valid manual netowrk parameters"do
      nm = Bosh::Aliyun::NetworkManager.new("manual_test" => manual)
      expect(nm.network.configure).to eq({
        :PrivateIpAddress => "127.0.0.1",
        :VSwitchId => "vsw-23priv72n",
        :SecurityGroupId => "sg-238ux30qe"
      })

      expect(nm.vpc?).to be(true)

      expect(nm.private_ip).to eq("127.0.0.1")
    end
  end

  describe "VIP network" do
    it "should return a valid vip network parameters", :debug => true do
      nm = Bosh::Aliyun::NetworkManager.new("vip_test" => vip)
      expect(nm.vip_network.configure).to eq({
        :Bandwidth => "10",
        :InternetChargeType => "PayByBandwidth"
        })
    end
  end

end
