#!/usr/bin/env bash

set -e

source bosh-cpi-src/ci/tasks/utils.sh
source bosh-cpi-src/ci/tasks/credentials.sh

: ${test_role_arn:?}
: ${METADATA_FILE:=environment/metadata}
: ${CPI_CREDENTIAL_SOURCE:=static}

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
exportMetadata2Env CPI_SLB_SERVER_GROUP_ID server_group_slb
exportMetadata2Env RAM_ROLE_NAME ram_role
export CIDR_NOTATION=$(getCidrNotation $CPI_INTERNAL_CIDR)
export CPI_INTERNAL_NETMASK=$(cdr2mask $CIDR_NOTATION)

# Exchange the worker's instance role for the test role. The aliyun CLI reads the
# result from the environment, so no credential appears on a command line.
region="${CPI_REGION}" assume_pipeline_role "${test_role_arn}" "integration"

# The CPI under test authenticates the same way. `static` feeds it the assumed
# short-lived credential, which is also what exercises the static code path;
# `ecs_ram_role` makes it discover the role attached to this container instead.
export CPI_CREDENTIAL_SOURCE
if [[ "${CPI_CREDENTIAL_SOURCE}" == "static" ]]; then
  export CPI_ACCESS_KEY_ID=${ALIBABA_CLOUD_ACCESS_KEY_ID}
  export CPI_ACCESS_KEY_SECRET=${ALIBABA_CLOUD_ACCESS_KEY_SECRET}
  export CPI_SECURITY_TOKEN=${ALIBABA_CLOUD_SECURITY_TOKEN}
fi
echo "CPI credential source: ${CPI_CREDENTIAL_SOURCE}"

echo "Uploading raw image ${stemcell_image_name} to ${CPI_STEMCELL_OSS_BUCKET}..."
aliyun oss cp ${stemcell_image} oss://${CPI_STEMCELL_OSS_BUCKET}/${stemcell_image_name} --region ${CPI_REGION} --force

# Setup Go and run tests
echo "set go path..."
export GOPATH=${PWD}/bosh-cpi-src
export PATH=${GOPATH}/bin:$PATH

echo "go version..."
check_go_version $GOPATH

echo "do integration test..."
cd ${PWD}/bosh-cpi-src

make testintci

echo "Deleting raw image ${stemcell_image_name} from ${CPI_STEMCELL_OSS_BUCKET}..."
aliyun oss rm oss://${CPI_STEMCELL_OSS_BUCKET}/${stemcell_image_name} -r -f --region ${CPI_REGION}
