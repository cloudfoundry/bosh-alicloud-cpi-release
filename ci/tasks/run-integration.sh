#!/usr/bin/env bash

set -e

source bosh-cpi-src/ci/tasks/utils.sh
#source /etc/profile.d/chruby-with-ruby-2.1.2.sh

METADATA_FILE=$(pwd)/terraform-metadata/ci/assets/terraform/metadata

check_param ALICLOUD_ACCESS_KEY_ID
check_param ALICLOUD_DEFAULT_REGION

: ${ALICLOUD_ACCESS_KEY_ID:?}
: ${ALICLOUD_SECRET_ACCESS_KEY:?}
: ${ALICLOUD_DEFAULT_REGION:?}

# Stemcell stuff
export STEMCELL_VERSION=`cat stemcell/version`
export STEMCELL_FILE=`pwd`/stemcell/image.tgz
pushd stemcell
  tar -zxvf stemcell.tgz
  mv image image.tgz
popd

echo "get metadata from output..."
cat terraform-metadata/ci/assets/terraform/metadata

metadata=$(cat ${METADATA_FILE})

export BOSH_ALICLOUD_ACCESS_KEY_ID=${ALICLOUD_ACCESS_KEY_ID}
export BOSH_ALICLOUD_SECRET_ACCESS_KEY=${ALICLOUD_SECRET_ACCESS_KEY}
export ACCESS_KEY_CONFIG=${ALICLOUD_SECRET_ACCESS_KEY}
export ACCESS_KEY_ID=${ALICLOUD_ACCESS_KEY_ID}
export BOSH_ALICLOUD_REGION_ID=cn-beijing
#export BOSH_ALICLOUD_ZONE_ID=$(echo ${metadata} | jq -e --raw-output ".availability_zone")
export BOSH_ALICLOUD_ZONE_ID=$(cat metadata | grep availability_zone | awk -F : '{print $2}' | grep -o "[^ ]\+\( \+[^ ]\+\)*")
export BOSH_ALICLOUD_SECURITY_GROUP_ID=$(cat metadata | grep security_group_id | awk -F : '{print $2}' | grep -o "[^ ]\+\( \+[^ ]\+\)*")
export BOSH_ALICLOUD_VSWITCH_ID=$(cat metadata | grep vswitch_id | awk -F : '{print $2}' | grep -o "[^ ]\+\( \+[^ ]\+\)*")
export BOSH_CLI_SILENCE_SLOW_LOAD_WARNING=true

echo "vswitch id: "
echo $BOSH_ALICLOUD_VSWITCH_ID

# Setup Go and run tests
echo "set go path..."
export GOPATH=${PWD}/bosh-cpi-src
export PATH=${GOPATH}/bin:$PATH

echo "go version..."
check_go_version $GOPATH

echo "do integration test..."
cd ${PWD}/bosh-cpi-src
env
make
make testintci
