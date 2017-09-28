#!/usr/bin/env bash

set -e

: ${ALICLOUD_ACCESS_KEY_ID:?}
: ${ALICLOUD_SECRET_ACCESS_KEY:?}
: ${ALICLOUD_DEFAULT_REGION:?}

source bosh-cpi-src/ci/tasks/utils.sh

export GOPATH=${PWD}/bosh-cpi-src
export PATH=${GOPATH}/bin:$PATH

export ACCESS_KEY_ID=${ALICLOUD_ACCESS_KEY_ID}
export ACCESS_KEY_SECRET=${ALICLOUD_SECRET_ACCESS_KEY}

check_go_version $GOPATH

#cd ${PWD}/bosh-cpi-src/src/bosh-alicloud-cpi
cd ${PWD}/bosh-cpi-src
source .envrc

# fix https fetch failed: Get https://golang.org/x/net/html/charset?go-get=1: dial tcp 216.239.37.1:443: i/o timeout
# git config --global http.proxy http://127.0.0.1:1080 
# git config --global https.proxy https://127.0.0.1:1080

pushd src
  #go get github.com/golang/net
  mkdir -p golang.org/x/net
  pushd golang.org/x/net
    git clone https://github.com/golang/net
  popd
  go install
popd
make test
