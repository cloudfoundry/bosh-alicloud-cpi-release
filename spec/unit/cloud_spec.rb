require 'spec_helper'

describe Bosh::Alicloud::Cloud do

  describe "Unit test" do
    it 'can create stemcell' do
      # Init
      o = load_cloud_options
      o["alicloud"]["SecurityGroupId"] = "sg-237p56jii"
      c = Bosh::Alicloud::Cloud.new o

      # Create stemcell
      para = {}
      r = c.create_stemcell

      expect(r).to match(/[\w]{1}-[\w]{9}/)
    end

    it 'can delete stemcell' do
      o = load_cloud_options
      o["alicloud"]["SecurityGroupId"] = "sg-237p56jii"
      c = Bosh::Alicloud::Cloud.new o

      # Delete stemcell
      c.delete_stemcell ""

      expect(true).to eq(true)
    end

    it 'can create, reboot and delete a vm', :debug => true do
      o = load_cloud_options

      c = Bosh::Alicloud::Cloud.new o

      # Create VM
      ins_id = c.create_vm nil, nil, o["resource_pool"], o["network"]

      expect(ins_id).to match(/[\w]{1}-[\w]{9}/)

      # Reboot VM
      # c.reboot_vm ins_id

      # Delete VM
      # r = c.delete_vm ins_id

      # expect(r).to have_key("RequestId")
    end

    it 'can check vm status' do
      o = load_cloud_options
      # o["alicloud"]["SecurityGroupId"] = "sg-237p56jii"
      c = Bosh::Alicloud::Cloud.new o

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
      c = Bosh::Alicloud::Cloud.new o

      # Create VM with specific network
      ins_id = c.create_vm nil, nil, o["resource_pool"], o["network"]

      expect(ins_id).to match(/[\w]{1}-[\w]{9}/)

      r = c.delete_vm ins_id

      expect(r).to have_key("RequestId")
    end
  end

  describe "Bosh registry test" do

    it "can init a registry" do
      o = load_cloud_options
      c = Bosh::Alicloud::Cloud.new o

      registry_options = recursive_symbolize_keys(o["registry"])
      c.initialize_registry registry_options

      expect(true).to eq(true)
    end

    it "can initilize agent settings" do
      o = load_cloud_options
      c = Bosh::Alicloud::Cloud.new o

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
      c = Bosh::Alicloud::Cloud.new o

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

      if not r.nil?
        s = c.initial_agent_settings(ins_id, agent_id, networks, "")
        r.update_settings(ins_id, s)
      end
    end

  end

  describe "Disk management test" do

    it "can create a disk and delete it" do
      o = load_cloud_options
      c = Bosh::Alicloud::Cloud.new o

      # Create VM
      ins_id = c.create_vm nil, nil, nil, nil

      size = "1024" # default
      size = o["persistent_disk"] if o.has_key? "persistent_disk"

      disk_id = c.create_disk "1024", nil, ins_id
      expect(ins_id).to match(/[\w]{1}-[\w]{9}/)

      c.delete_disk disk_id
      c.delete_vm ins_id
    end

    it "can create a disk, attach it, detach it and delete it" do
      o = load_cloud_options
      c = Bosh::Alicloud::Cloud.new o

      # Create VM
      ins_id = c.create_vm nil, nil, o["resource_pool"], o["network"]

      size = "1024" # default
      size = o["persistent_disk"] if o.has_key? "persistent_disk"

      disk_id = c.create_disk "1024", nil, ins_id
      expect(disk_id).to match(/[\w]{1}-[\w]{9}/)

      c.attach_disk ins_id, disk_id
      c.detach_disk ins_id, disk_id

      c.delete_disk disk_id
      c.delete_vm ins_id
    end

  end
end
