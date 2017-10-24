#!/usr/bin/env bash

set -e

source bosh-cpi-src/ci/tasks/utils.sh
#source /etc/profile.d/chruby-with-ruby-2.1.2.sh

METADATA_FILE=$(pwd)/terraform-metadata/ci/assets/terraform/metadata

: ${ALICLOUD_ACCESS_KEY_ID:?}
: ${ALICLOUD_ACCESS_KEY_SECRET:?}
: ${CPI_STEMCELL_ID:?}
: ${CPI_INTERNAL_IP:?}
: ${CPI_INTERNAL_GW:?}

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

exportMetadata2Env(){
  export $1=$(cat ${METADATA_FILE} | grep $2 | awk -F : '{print $2}' | grep -o "[^ ]\+\( \+[^ ]\+\)*")
}


export CPI_ACCESS_KEY_ID=${ALICLOUD_ACCESS_KEY_ID}
export CPI_ACCESS_KEY_SECRET=${ALICLOUD_ACCESS_KEY_SECRET}
export CPI_STEMCELL_ID=${CPI_STEMCELL_ID}
export CPI_INTERNAL_IP=${CPI_INTERNAL_IP}
export CPI_INTERNAL_GW=${CPI_INTERNAL_GW}

exportMetadata2Env CPI_REGION region
exportMetadata2Env CPI_ZONE availability_zone
exportMetadata2Env CPI_SECURITY_GROUP_ID security_group_id
exportMetadata2Env CPI_VSWITCH_ID vswitch_id
exportMetadata2Env CPI_INTERNAL_CIDR cidr_block
exportMetadata2Env CPI_EXTERNAL_IP external_ip_address
export CIDR_NOTATION=$(getCidrNotation $CPI_INTERNAL_CIDR)
export CPI_INTERNAL_NETMASK=$(cdr2mask $CIDR_NOTATION)


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
