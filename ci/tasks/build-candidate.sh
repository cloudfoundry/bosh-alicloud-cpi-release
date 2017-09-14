#!/bin/bash

set -x

source /etc/profile.d/chruby.sh
chruby 2.1.2

touch hope.log

semver=`cat version-semver/number`

pushd bosh-cpi-src
  echo "running unit tests"
  pushd src/bosh_alicloud_cpi
    bundle install
    bundle exec rspec spec/unit/*
  popd

  echo "using bosh CLI version..."
  bosh version

  cpi_release_name="bosh-alicloud-cpi"

  echo "building CPI release..."
  bosh create release --name $cpi_release_name --version $semver --with-tarball
popd

mv bosh-cpi-src/dev_releases/$cpi_release_name/$cpi_release_name-$semver.tgz candidate/