$LOAD_PATH.unshift File.expand_path('../../lib', __FILE__)
require 'bosh_aliyun_cpi'

RSpec.configure do |config|
  config.before do
    logger = Logger.new('/dev/null')
    allow(Bosh::Clouds::Config).to receive(:logger).and_return(logger)
  end
end
