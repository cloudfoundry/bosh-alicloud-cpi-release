#!/usr/bin/env bash

set -e

: ${ALICLOUD_DEFAULT_REGION:?}
: ${cleanup_role_arn:?}
: ${instance_termination_timeout:=600}

source bosh-cpi-src/ci/tasks/utils.sh
source bosh-cpi-src/ci/tasks/credentials.sh

# Assume the provisioning role from scratch: this task runs in the `ensure` of a
# failed job, so it cannot rely on anything the main task produced.
region="${ALICLOUD_DEFAULT_REGION}" assume_pipeline_role "${cleanup_role_arn}" "cleanup"

metadata=$(cat environment/metadata)
vpc_id=$(echo ${metadata} | jq --raw-output ".vpc_id")

vpc_instance_ids() {
  aliyun ecs DescribeInstances \
    --VpcId "${vpc_id}" \
    --PageSize 100 \
    --region "${ALICLOUD_DEFAULT_REGION}" \
    | jq -r '.Instances.Instance[].InstanceId'
}

echo "Checking whether there still exists instances in the VPC ${vpc_id} and then delete them..."
if [[ ! -z "${vpc_id}" ]] ; then
  for inst in $(vpc_instance_ids)
  do
    echo -e "Deleting instance ${inst}..."
    aliyun ecs DeleteInstance \
            --InstanceId ${inst} \
            --Force true \
            --region ${ALICLOUD_DEFAULT_REGION}
  done

  # DeleteInstance is asynchronous. Returning before the instances are actually
  # gone makes the terraform destroy that follows fail on DependencyViolation
  # for the security group and on the key pair still being attached, which
  # leaks the whole VPC and then collides with the next build.
  echo "Waiting for instances in ${vpc_id} to terminate..."
  deadline=$(( $(date +%s) + instance_termination_timeout ))
  while true; do
    remaining=$(vpc_instance_ids | wc -l | tr -d ' ')
    if [[ "${remaining}" == "0" ]]; then
      echo "No instances left in ${vpc_id}."
      break
    fi
    if [[ $(date +%s) -ge ${deadline} ]]; then
      echo "Still ${remaining} instance(s) in ${vpc_id} after ${instance_termination_timeout}s;" \
           "the terraform destroy that follows will likely fail." >&2
      vpc_instance_ids >&2
      break
    fi
    echo "  ${remaining} instance(s) remaining..."
    sleep 15
  done
fi
