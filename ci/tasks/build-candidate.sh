#!/bin/bash

set -e

: ${GIT_USER_EMAIL:?}
: ${GIT_USER_NAME:?}

source bosh-cpi-src/ci/tasks/utils.sh
source /etc/profile.d/chruby-with-ruby-2.1.2.sh

cpi_release_name="bosh-alicloud-cpi"
semver=`cat version-semver/number`

CURRENT_PATH=$(pwd)
DESC=$CURRENT_PATH/bosh-cpi-dev-artifacts/

# install bosh
echo "Installing Bosh CLI..."
curl -O https://s3.amazonaws.com/bosh-cli-artifacts/bosh-cli-2.0.28-linux-amd64
chmod +x ./bosh-cli-*
mv ./bosh-cli-* /usr/local/bin/bosh2

git config --global user.email ${GIT_USER_EMAIL}
git config --global user.name ${GIT_USER_NAME}

pushd bosh-cpi-src
  echo "using bosh CLI version..."
  bosh2 -v

  #echo $GOPATH
  source .envrc

  # add go cpi blob
  bosh2 add-blob ../go-cpi-blobs/go1.8.1.linux-amd64.tar.gz go1.8.1.linux-amd64.tar.gz

  make

  # Git repository has local modifications:
  # M ci/tasks/build-candidate.sh
  git add .
  git commit -m 'do nothing'

  echo "building CPI release..."
  # refers: https://bosh.io/docs/cli-v2#create-release
  bosh2 create-release --name $cpi_release_name --version $semver --tarball $cpi_release_name-$semver.tgz

  rm -rf ${DESC}/*
  mv $cpi_release_name-$semver.tgz ${DESC}/
popd

cp -r bosh-cpi-dev-artifacts candidate/repo
pushd candidate/repo
  ls -al ${DESC}/
  echo "git status..."
  git status
  git add .
  git commit -m 'create cpi release $cpi_release_name-$semver.tgz'
popd