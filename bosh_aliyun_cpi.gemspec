# coding: utf-8
lib = File.expand_path('../lib', __FILE__)
$LOAD_PATH.unshift(lib) unless $LOAD_PATH.include?(lib)
require 'version'

Gem::Specification.new do |spec|
  spec.name          = "bosh_aliyun_cpi"
  spec.version       = BoshAliyunCpi::VERSION
  spec.authors       = ["顺岭"]
  spec.email         = ["changlong.wcl@alibaba-inc.com"]

  spec.summary       = 'BOSH Aliyun CPI'
  spec.description   = 'This is the BOSH cloud platform interface for Aliyun, which is the biggest Infrastructure as a service in China.'

#  spec.files         = `git ls-files -z`.split("\x0").reject { |f| f.match(%r{^(test|spec|features)/}) }
  spec.bindir        = "bin"
  spec.executables   = spec.files.grep(%r{^bin/}) { |f| File.basename(f) }
  spec.require_paths = ["lib"]

  spec.add_development_dependency "bundler", "~> 1.10"
  spec.add_development_dependency "rake", "~> 10.0"
  spec.add_development_dependency "rspec", "~> 3.0"
  spec.add_development_dependency "webmock"

  spec.add_dependency 'bosh_common'
  spec.add_dependency 'bosh_cpi'
  spec.add_dependency 'bosh-registry'
  spec.add_dependency 'yajl-ruby',     '>=0.8.2'
  spec.add_dependency 'httpclient',    '=2.7.1'
end
