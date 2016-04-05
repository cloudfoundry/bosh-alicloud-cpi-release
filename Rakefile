require "bundler/setup"
require "bundler/gem_tasks"
require "rspec/core/rake_task"

RSpec::Core::RakeTask.new(:spec)

task :default => :spec

namespace :tools do
  desc "Delete created vms"
  task :delete_vm, :id, :running do |t, args|
    $LOAD_PATH.unshift File.expand_path('../../lib', __FILE__)
    require 'bosh_aliyun_cpi'

    o = {:aliyun => {:access_key_id => ENV['ACCESS_KEY_ID'], :secret => ENV['SECRET'], :region_id => ENV['RegionId'] }}

    c = Bosh::Aliyun::Cloud.new o

    c.stop_it args[:id] if args[:running] == "1"

    c.delete_vm args[:id]
  end
end
