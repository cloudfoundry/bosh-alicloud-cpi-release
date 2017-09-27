#!/bin/bash

set -x

source bosh-cpi-src/ci/tasks/utils.sh
source /etc/profile.d/chruby-with-ruby-2.1.2.sh

cpi_release_name="bosh-alicloud-cpi"
semver=`cat version-semver/number`

# install bosh
echo "Installing Bosh CLI..."
curl -O https://s3.amazonaws.com/bosh-cli-artifacts/bosh-cli-2.0.28-linux-amd64
chmod +x ./bosh-cli-*
mv ./bosh-cli-* /usr/local/bin/bosh2

cp -r bosh-cpi-src candidate/repo

pushd candidate/repo
  #echo "running unit tests"
  #pushd src/bosh_alicloud_cpi
  #  bundle install
  #  # bundle exec rspec spec/unit/*
  #  bundle exec rspec --tag debug spec/unit/cloud_spec.rb
  #popd

  echo "using bosh CLI version..."
  bosh2 -v

  cpi_release_name="bosh-alicloud-cpi"

  pwd
  make
  
  # add go cpi blob
  #ls ../../go-cpi-blobs
  bosh2 add-blob ../../go-cpi-blobs/go1.9.linux-amd64.tar.gz go1.9.linux-amd64.tar.gz

  export TERM=msys
  git status

  # todo: get email and user from params
  git config --global user.email "demon.wy@alibaba-inc.com"
  git config --global user.name "demonwy"

  git add .
  git commit -m 'do nothing'
  #git pull

  echo "building CPI release..."
  # refers: https://bosh.io/docs/cli-v2#create-release
  bosh2 create-release --name $cpi_release_name --version $semver --tarball $cpi_release_name-$semver.tgz

  mkdir dev-release-artifacts
  #pushd dev-release-artifacts
  #  rm -f *
  #popd

  mv $cpi_release_name-$semver.tgz dev-release-artifacts/
  git add .
  git commit -m 'create cpi release'
popd