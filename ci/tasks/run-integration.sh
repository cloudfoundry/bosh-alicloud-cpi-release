#!/usr/bin/env bash

set -e

: ${ALICLOUD_ACCESS_KEY_ID:?}
: ${ALICLOUD_SECRET_ACCESS_KEY:?}
: ${ALICLOUD_DEFAULT_REGION:?}

# NOTE: To run with specific line numbers, set:
RSPEC_ARGUMENTS="spec/integration/lifecycle_spec.rb:mm:nn"
: ${RSPEC_ARGUMENTS:=spec/integration}
#: ${METADATA_FILE:=environment/metadata}

echo "begin test1.."

release_dir="$( cd $(dirname $0) && cd ../.. && pwd )"

if [[ -f "/etc/profile.d/chruby.sh" ]] ; then
  source /etc/profile.d/chruby.sh
  chruby 2.1.2
fi

echo "begin test2.."

#metadata=$(cat ${METADATA_FILE})

echo ${ALICLOUD_ACCESS_KEY_ID}

export BOSH_ALICLOUD_ACCESS_KEY_ID=${ALICLOUD_ACCESS_KEY_ID}
export BOSH_ALICLOUD_SECRET_ACCESS_KEY=${ALICLOUD_SECRET_ACCESS_KEY}
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

pushd ${release_dir}/src/bosh_alicloud_cpi > /dev/null
  ls
  echo "bundle install"

  bundle install
  bundle exec rspec ${RSPEC_ARGUMENTS}
popd > /dev/null