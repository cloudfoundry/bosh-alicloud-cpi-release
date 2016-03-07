require 'spec_helper'

describe Bosh::Aliyun::Cloud do

  FAKE_CLOUD_PROPERTIES = {"a" => "b"}

  let(:stemcell_id) { "test_stemcell_name" }

  describe '#initialize' do
    it 'when all the required configurations are present' do
      expect(true).to eq(true) 
    end

  end

  describe '#create_stemcell' do
    let(:cloud_properties) { {} }
    let(:image_path) { "fake-image-path" }

    it 'should create a stemcell' do
      stemcell_manager = Bosh::Aliyun::StemcellManager.new(FAKE_CLOUD_PROPERTIES)
      expect(stemcell_manager).to receive(:create_stemcell).
        with(image_path, cloud_properties).and_return(stemcell_id)

      expect(stemcell_manager.create_stemcell(image_path, cloud_properties)).to eq(stemcell_id)
    end
  end
end
