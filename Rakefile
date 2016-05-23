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

  desc "Allocate public IP address"
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

  desc "Remove unbund elastic IP addresses"
  task :remove_elastic_ip, :region_id do |t, args|
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

    logger.info("get all elastic IP addresses")
    param = {
      :RegionId => args[:region_id],
      :Status   => "Available"
    }

    r = cli.DescribeEipAddresses param

    a_ids = r["EipAddresses"]["EipAddress"].map {|item| item["AllocationId"]}
    logger.info("Avaiable elastic IP addresses are #{a_ids}")

    a_ids.each do |id|
      param = {
        :AllocationId => id
      }

      r = cli.ReleaseEipAddress param
      logger.info("Release #{a_ids} success") if r["IncorrectEipStatus"].nil?
    end
  end

end
