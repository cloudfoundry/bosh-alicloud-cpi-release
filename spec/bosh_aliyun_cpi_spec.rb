require 'spec_helper'

describe BoshAliyunCpi do
  it 'has a version number' do
    expect(BoshAliyunCpi::VERSION).not_to be nil
  end

  it 'can load client options' do
    o = load_client_options

    expect(o).to have_key(:AccessKeyId)
    expect(o).to have_key(:AccessKey)
  end

  it 'can load cloud options' do
    o = load_cloud_options

    expect(o).to have_key("aliyun")
  end

end
