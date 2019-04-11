#!/usr/bin/env bash

set -e

source bosh-cpi-src/ci/tasks/utils.sh
#source /etc/profile.d/chruby-with-ruby-2.1.2.sh

: ${ALICLOUD_ACCESS_KEY_ID:?}
: ${ALICLOUD_ACCESS_KEY_SECRET:?}
: ${CPI_STEMCELL_ID:?}
: ${METADATA_FILE:=environment/metadata}

tar -xzf aliyun-cli/aliyun-cli-linux-amd64.tar.gz -C /usr/bin

# Stemcell stuff
export CPI_STEMCELL_VERSION=`cat stemcell/version`
export CPI_STEMCELL_FILE=`pwd`/stemcell/image.tgz

pushd stemcell
  tar -zxvf stemcell.tgz
  cp image image.tgz

  tar -xzvf image
  stemcell_image="$(echo ${PWD}/root.*)"
  stemcell_image_name="$(basename ${stemcell_image})"
  export CPI_STEMCELL_OSS_OBJECT=${stemcell_image_name}
popd

echo "get metadata from output..."
cat ${METADATA_FILE}

exportMetadata2Env(){
  value=$(cat ${METADATA_FILE} | grep '^"'$2 | awk -F : '{print $2}' | grep -o "[^ ]\+\( \+[^ ]\+\)*")
  value2=${value//,/}
  export $1=${value2//\"/}
}


export CPI_ACCESS_KEY_ID=${ALICLOUD_ACCESS_KEY_ID}
export CPI_ACCESS_KEY_SECRET=${ALICLOUD_ACCESS_KEY_SECRET}
# todo: get stemcell id from STEMCELL.MF
export CPI_STEMCELL_ID=${CPI_STEMCELL_ID}

exportMetadata2Env CPI_REGION region
exportMetadata2Env CPI_ZONE zone
exportMetadata2Env CPI_SECURITY_GROUP_ID security_group_id
exportMetadata2Env CPI_VSWITCH_ID vswitch_id
exportMetadata2Env CPI_INTERNAL_CIDR internal_cidr
exportMetadata2Env CPI_INTERNAL_IP internal_ip
exportMetadata2Env CPI_INTERNAL_GW internal_gw
exportMetadata2Env CPI_EXTERNAL_IP external_ip
exportMetadata2Env CPI_STEMCELL_OSS_BUCKET integration_bucket
exportMetadata2Env CPI_SLB_ID slb
exportMetadata2Env RAM_ROLE_NAME ram_role
export CIDR_NOTATION=$(getCidrNotation $CPI_INTERNAL_CIDR)
export CPI_INTERNAL_NETMASK=$(cdr2mask $CIDR_NOTATION)

echo "Uploading raw image ${stemcell_image_name} to ${CPI_STEMCELL_OSS_BUCKET}..."
aliyun oss cp ${stemcell_image} oss://${CPI_STEMCELL_OSS_BUCKET}/${stemcell_image_name} --access-key-id ${CPI_ACCESS_KEY_ID} --access-key-secret ${CPI_ACCESS_KEY_SECRET} --region ${CPI_REGION}

# Setup Go and run tests
echo "set go path..."
export GOPATH=${PWD}/bosh-cpi-src
export PATH=${GOPATH}/bin:$PATH

echo "go version..."
check_go_version $GOPATH

echo "do integration test..."
cd ${PWD}/bosh-cpi-src
env

make testintci

echo "Deleting raw image ${stemcell_image_name} from ${CPI_STEMCELL_OSS_BUCKET}..."
aliyun oss rm oss://${CPI_STEMCELL_OSS_BUCKET}/${stemcell_image_name} -r -f --access-key-id ${CPI_ACCESS_KEY_ID} --access-key-secret ${CPI_ACCESS_KEY_SECRET} --region ${CPI_REGION}
