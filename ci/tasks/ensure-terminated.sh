#!/usr/bin/env bash

set -e

: ${ALICLOUD_DEFAULT_REGION:?}
: ${cleanup_role_arn:?}

source bosh-cpi-src/ci/tasks/utils.sh
source bosh-cpi-src/ci/tasks/credentials.sh

# Assume the provisioning role from scratch: this task runs in the `ensure` of a
# failed job, so it cannot rely on anything the main task produced.
region="${ALICLOUD_DEFAULT_REGION}" assume_pipeline_role "${cleanup_role_arn}" "cleanup"

metadata=$(cat environment/metadata)
vpc_id=$(echo ${metadata} | jq --raw-output ".vpc_id")

echo "Checking whether there still exists instances in the VPC ${vpc_id} and then delete them..."
if [[ ! -z "${vpc_id}" ]] ; then
  instance_ids="$(echo $(aliyun ecs DescribeInstances \
    --VpcId ${vpc_id} \
    --region ${ALICLOUD_DEFAULT_REGION}
    ) | jq -r '.Instances.Instance[].InstanceId'
    )"

  # if it's not an empty string (of any length)...
  for inst in ${instance_ids}
  do
    echo -e "Deleting instance ${inst}..."
    aliyun ecs DeleteInstance \
            --InstanceId ${inst} \
            --Force true \
            --region ${ALICLOUD_DEFAULT_REGION}
  done
fi
