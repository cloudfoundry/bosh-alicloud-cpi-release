require 'spec_helper'

describe Bosh::Aliyun::Cloud do

  describe "Classic network test" do
    it 'can create stemcell' do
      # Init
      o = load_cloud_options
      o["aliyun"]["SecurityGroupId"] = "sg-237p56jii"
      c = Bosh::Aliyun::Cloud.new o

      # Create stemcell
      para = {}
      r = c.create_stemcell

      expect(r).to match(/[\w]{1}-[\w]{9}/)
    end

    it 'can delete stemcell' do
      o = load_cloud_options
      o["aliyun"]["SecurityGroupId"] = "sg-237p56jii"
      c = Bosh::Aliyun::Cloud.new o

      # Delete stemcell
      c.delete_stemcell ""

      expect(true).to eq(true)
    end

    # TODO We need to find a proper way(dynamic ip`) to reuse this test
    # it 'can create, reboot and delete a vm' do
    #   o = load_cloud_options
    #   o["aliyun"]["SecurityGroupId"] = "sg-237p56jii"
    #   c = Bosh::Aliyun::Cloud.new o
    #
    #   # Create VM
    #   ins_id = c.create_vm
    #
    #   expect(ins_id).to match(/[\w]{1}-[\w]{9}/)
    #
    #   c.reboot_vm ins_id
    #
    #   r = c.stop_it ins_id
    #   r = c.delete_vm ins_id
    #
    #   expect(r).to have_key("RequestId")
    # end

    it 'can check vm status' do
      o = load_cloud_options
      o["aliyun"]["SecurityGroupId"] = "sg-237p56jii"
      c = Bosh::Aliyun::Cloud.new o

      # TODO Right now just try to match the bootstrap vm
      # if it failed, please use a new vm id.
      ins_id = "i-236kv7mlm"

      r = c.has_vm? ins_id

      expect(r).to eq(true)
    end
  end

  describe "VPC network test", :debug => true do

    it "can create a vm with both private and public network" do
      o = load_cloud_options
      c = Bosh::Aliyun::Cloud.new o

      # Create VM with specific network
      ins_id = c.create_vm nil, nil, nil, o["networks"]

      expect(ins_id).to match(/[\w]{1}-[\w]{9}/)

      r = c.stop_it ins_id
      r = c.delete_vm ins_id

      expect(r).to have_key("RequestId")
    end
  end


end
