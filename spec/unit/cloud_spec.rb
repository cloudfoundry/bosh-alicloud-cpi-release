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

  describe "VPC network test" do

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

  describe "Bosh registry test", :debug => true do

    it "can init a registry" do
      o = load_cloud_options
      c = Bosh::Aliyun::Cloud.new o

      registry_options = recursive_symbolize_keys(o["registry"])
      c.initialize_registry registry_options

      expect(true).to eq(true)
    end

    it "can initilize agent settings" do
      o = load_cloud_options
      c = Bosh::Aliyun::Cloud.new o

      registry_options = recursive_symbolize_keys(o["registry"])
      c.initialize_registry registry_options

      ins_id = "aaa"
      agent_id = "dddd"
      networks = {
        :type => "manual",
        :ip => "127.0.0.1",
        :cloud_properties => {
          :VSwitchId => "vsw-23priv72n",
          :SecurityGroupId => "sg-238ux30qe"
        }
      }

      s = c.initial_agent_settings(ins_id, agent_id, networks, "")

      expect(s).to have_key("vm")
      expect(s).to have_key("agent_id")
      expect(s).to have_key("networks")
      expect(s).to have_key("disks")
    end

    it "can update agent settings" do
      o = load_cloud_options
      c = Bosh::Aliyun::Cloud.new o

      registry_options = recursive_symbolize_keys(o["registry"])
      r = c.initialize_registry registry_options

      ins_id = "aaa"
      agent_id = "dddd"
      networks = {
        :type => "manual",
        :ip => "127.0.0.1",
        :cloud_properties => {
          :VSwitchId => "vsw-23priv72n",
          :SecurityGroupId => "sg-238ux30qe"
        }
      }

      s = c.initial_agent_settings(ins_id, agent_id, networks, "")

      p s

      r.update_settings(ins_id, s)

    end

  end
end
