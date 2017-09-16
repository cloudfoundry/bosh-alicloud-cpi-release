#!/bin/bash

set -x

source /etc/profile.d/chruby.sh
chruby 2.1.2

semver=`cat version-semver/number`

# install bosh
echo "installing bosh..."
curl -O curl -O https://s3.amazonaws.com/bosh-cli-artifacts/bosh-cli-2.0.28-linux-amd64
chmod +x ./bosh-cli-*
mv ./bosh-cli-* /usr/local/bin/bosh

pushd bosh-cpi-src
  echo "running unit tests"
  pushd src/bosh_alicloud_cpi
    bundle install
    # bundle exec rspec spec/unit/*
    bundle exec rspec --tag debug spec/unit/cloud_spec.rb
  popd

  echo "using bosh CLI version..."
  bosh -v

  cpi_release_name="bosh-alicloud-cpi"

  echo "building CPI release..."
  # refers: https://bosh.io/docs/cli-v2#create-release
  bosh create-release --name $cpi_release_name --version $semver
popd

ls bosh-cpi-src
ls candidate

mv bosh-cpi-src/dev_releases/$cpi_release_name/$cpi_release_name-$semver.tgz candidate/