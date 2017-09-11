require 'spec_helper'

describe Bosh::Alicloud::Cloud do

  describe "Unit test" do
    it 'can create stemcell' do
      # Init
      o = load_cloud_options
      o["alicloud"]["SecurityGroupId"] = "sg-237p56jii"
      c = Bosh::Alicloud::Cloud.new o

      # Create stemcell
      para = {
        :image_id => {
          "cn-beijing" => "m-2zeggz4i4n2z510ajcvw",
          "cn-hangzhou" => "m-bp1bidv1aeiaynlyhmu9"
        }
      }
      r = c.create_stemcell nil, para

      expect(r).to match(/[\w]{1}-[\w]{9}/)
    end

    it 'can delete stemcell' do
      o = load_cloud_options
      o["alicloud"]["SecurityGroupId"] = "sg-237p56jii"
      c = Bosh::Alicloud::Cloud.new o

      # Delete stemcell
      c.delete_stemcell "m-bp1bidv1aeiaynlyhmu9"

      expect(true).to eq(true)
    end

    it 'can create, reboot and delete a vm', :debug => true do
      o = load_cloud_options

      c = Bosh::Alicloud::Cloud.new o

      # Create VM
      ins_id = c.create_vm nil, nil, o["resource_pool"], o["network"]

      expect(ins_id).to match(/[\w]{1}-[\w]{9}/)

      # Check VM status
      r = c.has_vm? ins_id
      expect(r).to eq(true)

      # Reboot VM
      c.reboot_vm ins_id

      # Delete VM
      r = c.delete_vm ins_id

      expect(r).to have_key("RequestId")
    end
  end

  describe "Bosh registry test" do

    it "can init a registry" do
      o = load_cloud_options
      c = Bosh::Alicloud::Cloud.new o

      registry_options = recursive_symbolize_keys(o["registry"])
      c.initialize_registry

      expect(true).to eq(true)
    end

    it "can initilize agent settings" do
      o = load_cloud_options
      c = Bosh::Alicloud::Cloud.new o

      registry_options = recursive_symbolize_keys(o["registry"])
      c.initialize_registry

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
      r = c.initialize_registry

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
      ins_id = c.create_vm nil, nil, o["resource_pool"], o["network"]

      size = "20"
      size = o["persistent_disk"] if o.has_key? "persistent_disk"

      disk_id = c.create_disk size, nil, ins_id
      expect(ins_id).to match(/[\w]{1}-[\w]{9}/)

      c.delete_disk disk_id
      c.delete_vm ins_id
    end

    it "can create a disk, attach it, detach it and delete it" do
      o = load_cloud_options
      c = Bosh::Alicloud::Cloud.new o

      # Create VM
      ins_id = c.create_vm nil, nil, o["resource_pool"], o["network"]

      size = "20"
      size = o["persistent_disk"] if o.has_key? "persistent_disk"

      disk_id = c.create_disk "20", nil, ins_id
      expect(disk_id).to match(/[\w]{1}-[\w]{9}/)

      c.attach_disk ins_id, disk_id
      c.detach_disk ins_id, disk_id

      c.delete_disk disk_id
      c.delete_vm ins_id
    end

  end
end
