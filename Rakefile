require "bundler/setup"
require "bundler/gem_tasks"
require "rspec/core/rake_task"
require "yaml"

require 'bosh_aliyun_cpi'

RSpec::Core::RakeTask.new(:spec)

task :default => :spec

namespace :tools do
  desc "Delete created vms"
  task :delete_vm, :id, :running do |t, args|
    cloud_config = OpenStruct.new(:logger => Logger.new(STDERR))

    cloud_config[:logger].datetime_format = '%Y-%m-%d %H:%M:%S'
    cloud_config[:logger].formatter = proc do |severity, datetime, progname, msg|
      "[#{severity}], #{datetime} #{caller[4]}:#{__LINE__}: #{msg}\n"
    end
    cloud_config[:logger].level = Logger::INFO
    Bosh::Clouds::Config.configure(cloud_config)

    o = YAML.load_file('spec/assets/cpi_config')
    c = Bosh::Aliyun::Cloud.new o

    c.stop_it args[:id] # if args[:running] == "1"

    c.delete_vm args[:id]
  end

  desc "Allocate public ip address"
  task :allocate_public_ip, :id do |t, args|
    cloud_config = OpenStruct.new(:logger => Logger.new(STDERR))

    cloud_config[:logger].datetime_format = '%Y-%m-%d %H:%M:%S'
    cloud_config[:logger].formatter = proc do |severity, datetime, progname, msg|
      "[#{severity}], #{datetime} #{caller[4]}:#{__LINE__}: #{msg}\n"
    end
    cloud_config[:logger].level = Logger::INFO
    Bosh::Clouds::Config.configure(cloud_config)

    logger = Bosh::Clouds::Config.logger

    o = YAML.load_file('spec/assets/client_config')
    c = {}
    c[:AccessKeyId] = o["AccessKeyId"]
    c[:AccessKey] = o["AccessKey"]
    cli = Bosh::Aliyun::Client.new c, logger

    logger.info("start to allocate public ip for vm")
    param={
      :InstanceId => args[:id]
    }
    r = cli.AllocatePublicIpAddress param
    logger.info("the ip adderss is allocated")

    r
   end

end
