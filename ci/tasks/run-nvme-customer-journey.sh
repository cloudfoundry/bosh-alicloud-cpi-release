#!/usr/bin/env bash

set -euo pipefail

source bosh-cpi-src/ci/tasks/utils.sh
source director-state/director.env

METADATA_FILE=$(realpath environment/metadata)
DEPLOYMENT_NAME=e2e-test
NVME_INSTANCE=nvme-upgrade-test/0
MARKER_FILE=/var/vcap/store/nvme-upgrade-test/customer-data

: "${ALICLOUD_ACCESS_KEY_ID:?}"
: "${ALICLOUD_SECRET_ACCESS_KEY:?}"
: "${LEGACY_INSTANCE_TYPE:?}"
: "${NVME_INSTANCE_TYPE:?}"

stemcell_name=$(bosh int <(tar xfO "$(realpath stemcell/*.tgz)" stemcell.MF) --path /name)
region=$(jq -er '.region' "${METADATA_FILE}")

bosh -n update-cloud-config \
  -l "${METADATA_FILE}" \
  -v "legacy_instance_type=${LEGACY_INSTANCE_TYPE}" \
  -v "nvme_instance_type=${NVME_INSTANCE_TYPE}" \
  bosh-cpi-src/ci/assets/e2e-test-release/cloud-config.yml

deploy_customer_phase() {
  local vm_type=$1
  local disk_type=$2

  time bosh -n deploy -d "${DEPLOYMENT_NAME}" \
    -v "stemcell_name=${stemcell_name}" \
    -v "nvme_vm_type=${vm_type}" \
    -v "nvme_disk_type=${disk_type}" \
    -v "access_key=${ALICLOUD_ACCESS_KEY_ID}" \
    -v "secret_key=${ALICLOUD_SECRET_ACCESS_KEY}" \
    -l "${METADATA_FILE}" \
    bosh-cpi-src/ci/assets/e2e-test-release/manifest.yml
}

deployment_row() {
  bosh -d "${DEPLOYMENT_NAME}" instances --details --json | \
    jq -ce '.Tables[]?.Rows[]? | select(.instance | startswith("nvme-upgrade-test/"))'
}

vm_cid() {
  deployment_row | jq -er '.vm_cid | select(length > 0)'
}

disk_cid() {
  deployment_row | jq -er '
    .disk_cids
    | if type == "array" then .[0] else split(",")[0] end
    | select(length > 0)'
}

assert_remote_state() {
  local expected_marker=$1

  bosh -n -d "${DEPLOYMENT_NAME}" ssh "${NVME_INSTANCE}" -c \
    "sudo test \"\$(cat '${MARKER_FILE}')\" = '${expected_marker}' && findmnt -rn /var/vcap/data && findmnt -rn /var/vcap/store"

  deployment_row | jq -e '.process_state == "running"' >/dev/null
}

assert_iaas_state() {
  local expected_instance_type=$1
  local expected_disk_category=$2
  local current_vm_cid
  local current_disk_cid
  local instance_json
  local disk_json

  current_vm_cid=$(vm_cid)
  current_disk_cid=$(disk_cid)

  instance_json=$(aliyun ecs DescribeInstances \
    --InstanceIds "[\"${current_vm_cid}\"]" \
    --access-key-id "${ALICLOUD_ACCESS_KEY_ID}" \
    --access-key-secret "${ALICLOUD_SECRET_ACCESS_KEY}" \
    --region "${region}")
  echo "${instance_json}" | jq -e \
    --arg expected "${expected_instance_type}" \
    '.Instances.Instance[0].InstanceType == $expected' >/dev/null

  disk_json=$(aliyun ecs DescribeDisks \
    --DiskIds "[\"${current_disk_cid}\"]" \
    --access-key-id "${ALICLOUD_ACCESS_KEY_ID}" \
    --access-key-secret "${ALICLOUD_SECRET_ACCESS_KEY}" \
    --region "${region}")
  echo "${disk_json}" | jq -e \
    --arg expected "${expected_disk_category}" \
    '.Disks.Disk[0].Category == $expected and .Disks.Disk[0].Encrypted == true' >/dev/null
}

# Establish the customer state that predates NVMe support.
deploy_customer_phase nvme_upgrade_legacy nvme_upgrade_legacy

marker="nvme-customer-journey-$(date +%s)"
bosh -n -d "${DEPLOYMENT_NAME}" ssh "${NVME_INSTANCE}" -c \
  "sudo mkdir -p '$(dirname "${MARKER_FILE}")' && echo '${marker}' | sudo tee '${MARKER_FILE}' >/dev/null"
assert_remote_state "${marker}"
assert_iaas_state "${LEGACY_INSTANCE_TYPE}" cloud_efficiency

# Upgrade compute and disks together, as a customer does when moving to c9i.
deploy_customer_phase nvme_upgrade_target nvme_upgrade_target
assert_remote_state "${marker}"
assert_iaas_state "${NVME_INSTANCE_TYPE}" cloud_essd

# Exercise Director's copy/migrate fallback for a requested disk shrink.
deploy_customer_phase nvme_upgrade_target nvme_upgrade_target_smaller
assert_remote_state "${marker}"
assert_iaas_state "${NVME_INSTANCE_TYPE}" cloud_essd

# Recreate on NVMe and verify agent bootstrap, both disk mounts and customer data.
time bosh -n -d "${DEPLOYMENT_NAME}" recreate "${NVME_INSTANCE}"
assert_remote_state "${marker}"
assert_iaas_state "${NVME_INSTANCE_TYPE}" cloud_essd

# Roll compute back to the legacy family while retaining the upgraded ESSD disk.
deploy_customer_phase nvme_upgrade_legacy nvme_upgrade_target_smaller
assert_remote_state "${marker}"
assert_iaas_state "${LEGACY_INSTANCE_TYPE}" cloud_essd
