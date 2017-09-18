#!/bin/bash

set -x

source /etc/profile.d/chruby.sh
chruby 2.1.2

semver=`cat version-semver/number`

# install bosh
echo "installing bosh..."
curl -O https://s3.amazonaws.com/bosh-cli-artifacts/bosh-cli-2.0.28-linux-amd64
chmod +x ./bosh-cli-*
mv ./bosh-cli-* /usr/local/bin/bosh

cp -r bosh-cpi-src candidate/repo

pushd candidate/repo
  echo "running unit tests"
  pushd src/bosh_alicloud_cpi
    bundle install
    # bundle exec rspec spec/unit/*
    bundle exec rspec --tag debug spec/unit/cloud_spec.rb
  popd

  echo "using bosh CLI version..."
  bosh -v

  cpi_release_name="bosh-alicloud-cpi"

  # add ruby cpi blob
  ls ../../ruby-cpi-blobs
  bosh add-blob ../../ruby-cpi-blobs/bundler-1.10.6.gem ruby_alicloud_cpi/bundler-1.10.6.gem
  bosh add-blob ../../ruby-cpi-blobs/ruby-2.1.7.tar.gz ruby_alicloud_cpi/ruby-2.1.7.tar.gz
  bosh add-blob ../../ruby-cpi-blobs/rubygems-2.4.8.tgz ruby_alicloud_cpi/rubygems-2.4.8.tgz
  bosh add-blob ../../ruby-cpi-blobs/yaml-0.1.5.tar.gz ruby_alicloud_cpi/yaml-0.1.5.tar.gz

  export TERM=msys
  git status

  # todo: get email and user from params
  git config --global user.email "demon.wy@alibaba-inc.com"
  git config --global user.name "demonwy"

  git add .
  git commit -m 'do nothing'
  #git pull
  git remote update

  echo "building CPI release..."
  # refers: https://bosh.io/docs/cli-v2#create-release
  bosh create-release --name $cpi_release_name --version $semver --tarball $cpi_release_name-$semver.tgz

  mkdir dev-release-artifacts
  mv $cpi_release_name-$semver.tgz dev-release-artifacts/
  git add .
  git commit -m 'create cpi release'
popd