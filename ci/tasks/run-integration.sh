#!/usr/bin/env bash

set -e

source bosh-cpi-src/ci/tasks/utils.sh
#source /etc/profile.d/chruby-with-ruby-2.1.2.sh

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

echo "get metadata..."
ls bosh-cpi-src/ci/assets/terraform
cat environment/ci/assets/terraform/metadata

#metadata=$(cat ${METADATA_FILE})

export BOSH_ALICLOUD_ACCESS_KEY_ID=${ALICLOUD_ACCESS_KEY_ID}
export BOSH_ALICLOUD_SECRET_ACCESS_KEY=${ALICLOUD_SECRET_ACCESS_KEY}
export ACCESS_KEY_CONFIG=${ALICLOUD_SECRET_ACCESS_KEY}
export ACCESS_KEY_ID=${ALICLOUD_ACCESS_KEY_ID}
#export BOSH_ALICLOUD_DEFAULT_KEY_NAME=$(echo ${metadata} | jq -e --raw-output ".default_key_name")
#export BOSH_ALICLOUD_REGION=$(echo ${metadata} | jq -e --raw-output ".region")
#export BOSH_ALICLOUD_SUBNET_ID=$(echo ${metadata} | jq -e --raw-output ".subnet_id")
#export BOSH_ALICLOUD_SUBNET_ZONE=$(echo ${metadata} | jq -e --raw-output ".az")
#export BOSH_ALICLOUD_LIFECYCLE_MANUAL_IP=$(echo ${metadata} | jq -e --raw-output ".internal_ip")
#export BOSH_ALICLOUD_ELB_ID=$(echo ${metadata} | jq -e --raw-output ".elb")
#export BOSH_ALICLOUD_TARGET_GROUP_NAME=$(echo ${metadata} | jq -e --raw-output ".alb_target_group")
#export BOSH_ALICLOUD_ELASTIC_IP=$(echo ${metadata} | jq -e --raw-output ".bats_eip")
#export BOSH_ALICLOUD_IPV6_IP=$(echo ${metadata} | jq -e --raw-output ".static_ipv6")
#export BOSH_ALICLOUD_KMS_KEY_ARN=$(echo ${metadata} | jq -e --raw-output ".alicloud_kms_key_arn")

export BOSH_CLI_SILENCE_SLOW_LOAD_WARNING=true

echo "ak: "
echo $ACCESS_KEY_ID
echo $ACCESS_KEY_CONFIG

ls ${PWD}
# Setup Go and run tests
export GOPATH=${PWD}/bosh-cpi-src
export PATH=${GOPATH}/bin:$PATH

check_go_version $GOPATH

cd ${PWD}/bosh-cpi-src
env
make
make testintci
