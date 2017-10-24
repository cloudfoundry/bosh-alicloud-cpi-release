#!/bin/bash

set -e

: ${GIT_USER_EMAIL:?}
: ${GIT_USER_NAME:?}
: ${GIT_USER_ID:?}
: ${GIT_USER_PASSWORD:?}
: ${BOSH_REPO_HOST:?}
: ${BOSH_REPO_BRANCH:?}

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

# todo: get email and user from params
git config --global user.email "demon.wy@alibaba-inc.com"
git config --global user.name "demonwy"

pushd candidate/repo
  git fetch https://${GIT_USER_ID}:${GIT_USER_PASSWORD}@${BOSH_REPO_HOST} ${BOSH_REPO_BRANCH}
  echo "****** git merge ******"
  git merge FETCH_HEAD

  echo "using bosh CLI version..."
  bosh2 -v

  cpi_release_name="bosh-alicloud-cpi"

  # fix cannot find package "bosh-alicloud-cpi/action
  source .envrc
  #echo $GOPATH

  make
  
  # add go cpi blob
  #ls ../../go-cpi-blobs
  bosh2 add-blob ../../go-cpi-blobs/go1.8.1.linux-amd64.tar.gz go1.8.1.linux-amd64.tar.gz

  # export TERM=msys
  echo "git status..."
  git status



  git add .
  git commit -m 'do nothing'
  #git pull

  ls -al

  echo "building CPI release..."
  # refers: https://bosh.io/docs/cli-v2#create-release
  bosh2 create-release --name $cpi_release_name --version $semver --tarball $cpi_release_name-$semver.tgz

  ls -al dev-release-artifacts
  mkdir dev-release-artifacts
  rm -rf dev-release-artifacts/*
  mv $cpi_release_name-$semver.tgz dev-release-artifacts/
  ls -al dev-release-artifacts

  echo "git status..."
  git status

  git add .
  git commit -m 'create cpi release'
popd