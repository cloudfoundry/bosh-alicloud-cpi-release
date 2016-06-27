$LOAD_PATH.unshift File.expand_path('../../lib', __FILE__)
require 'bosh_aliyun_cpi'
require 'yaml'
require 'yajl'

def load_client_options
  conf = YAML.load_file('spec/assets/client_config')

  c = {}
  c[:AccessKeyId] = conf["AccessKeyId"]
  c[:AccessKey] = conf["AccessKey"]

  c
end

def load_cloud_options
  YAML.load_file('spec/assets/cpi_config')
end

def recursive_symbolize_keys(h)
  case h
  when Hash
    Hash[
      h.map do |k, v|
        [ k.respond_to?(:to_sym) ? k.to_sym : k, recursive_symbolize_keys(v) ]
      end
    ]
  when Enumerable
    h.map { |v| recursive_symbolize_keys(v) }
  else
    h
  end
end

def mock_registry
  registry = double('registry',
    :endpoint => mock_registry_properties['endpoint'],
    :user     => mock_registry_properties['user'],
    :password => mock_registry_properties['password']
  )
  allow(Bosh::Registry::Client).to receive(:new).and_return(registry)
  registry
end

RSpec.configure do |config|
  config.before do
    logger = Logger.new(STDOUT)
    logger.level = Logger::DEBUG
    logger.datetime_format = '%Y-%m-%d %H:%M:%S'
    logger.formatter = proc do |severity, datetime, progname, msg|
      "[#{severity}], #{datetime} #{caller[4]}:#{__LINE__}: #{msg}\n"
    end

    allow(Bosh::Clouds::Config).to receive(:logger).and_return(logger)
  end
end
