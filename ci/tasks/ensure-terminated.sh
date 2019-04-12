#!/usr/bin/env bash

set -e

: ${ALICLOUD_ACCESS_KEY_ID:?}
: ${ALICLOUD_ACCESS_KEY_SECRET:?}
: ${ALICLOUD_DEFAULT_REGION:?}

source bosh-cpi-src/ci/tasks/utils.sh

metadata=$(cat environment/metadata)
vpc_id=$(echo ${metadata} | jq --raw-output ".vpc_id")

echo "Checking whether there still exists instances in the VPC ${vpc_id} and then delete them..."
if [[ ! -z "${vpc_id}" ]] ; then
  instance_ids="$(echo $(aliyun ecs DescribeInstances \
    --VpcId ${vpc_id} \
    --access-key-id ${ALICLOUD_ACCESS_KEY_ID} \
    --access-key-secret ${ALICLOUD_ACCESS_KEY_SECRET} \
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
            --access-key-id ${ALICLOUD_ACCESS_KEY_ID} \
            --access-key-secret ${ALICLOUD_ACCESS_KEY_SECRET} \
            --region ${ALICLOUD_DEFAULT_REGION}
  done
fi