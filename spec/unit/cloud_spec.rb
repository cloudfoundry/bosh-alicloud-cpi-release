require 'spec_helper'

describe Bosh::Aliyun::Cloud do

  it 'can create stemcell' do
    # Init
    o = load_cloud_options
    c = Bosh::Aliyun::Cloud.new o

    # Create stemcell
    para = {}
    r = c.create_stemcell

    expect(r).to match(/[\w]{1}-[\w]{9}/)
  end

  it 'can delete stemcell' do
    o = load_cloud_options
    c = Bosh::Aliyun::Cloud.new o

    # Delete stemcell
    c.delete_stemcell ""

    expect(true).to eq(true)
  end

  it 'can create, reboot and delete a vm', :debug => true do
    o = load_cloud_options
    c = Bosh::Aliyun::Cloud.new o

    # Create VM
    ins_id = c.create_vm

    expect(ins_id).to match(/[\w]{1}-[\w]{9}/)

    c.reboot_vm ins_id

    r = c.stop_it ins_id
    r = c.delete_vm ins_id

    expect(r).to have_key("RequestId")

  end

  it 'can check vm status' do
    o = load_cloud_options
    c = Bosh::Aliyun::Cloud.new o

    ins_id = "i-236kv7mlm"

    r = c.has_vm? ins_id

    expect(r).to eq(true)
  end

end
