require 'spec_helper'

describe BoshAliyunCpi do
  it 'has a version number' do
    expect(BoshAliyunCpi::VERSION).not_to be nil
  end

  it 'can load client options' do
    o = load_client_options

    expect(o).to have_key(:access_key_id)
    expect(o).to have_key(:secret)
  end

  # TODO : Read from asset yaml file
  it 'can load cloud options' do
    o = load_cloud_options

    expect(o).to have_key(:aliyun)
  end

end
