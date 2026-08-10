#!/usr/bin/env bash

set -euo pipefail

source bosh-cpi-src/ci/tasks/utils.sh
source bosh-cpi-src/ci/tasks/credentials.sh
source director-state/director.env

: "${observer_role_arn:?}"
: "${region:?}"
: "${LEGACY_INSTANCE_TYPE:?}"
: "${NVME_INSTANCE_TYPE:?}"

METADATA_FILE=$(realpath environment/metadata)
MANIFEST=bosh-cpi-src/ci/assets/e2e-test-release/nvme-journey.yml
CLOUD_CONFIG=bosh-cpi-src/ci/assets/e2e-test-release/cloud-config.yml
# Its own deployment, so the errands in manifest.yml are untouched and neither
# test waits on the other's VMs.
DEPLOYMENT_NAME=nvme-journey
INSTANCE_GROUP=nvme-upgrade-test
NVME_INSTANCE="${INSTANCE_GROUP}/0"
MARKER_FILE="/var/vcap/store/${INSTANCE_GROUP}/customer-data"

# The journey only describes instances, disks and images, so it runs under the
# read-only role rather than the one that can build and tear down the
# environment.
ensure_aliyun_cli

# The task is allowed to run for longer than an STS credential lives, so every
# round of cloud assertions starts from a fresh one instead of the journey dying
# of an expired token somewhere after the first hour.
refresh_cloud_credentials() {
  assume_pipeline_role "${observer_role_arn}" "nvme-journey"
}
refresh_cloud_credentials

# The e2e task in this same job already uploaded the light stemcell, which shares
# its name and version with the full one. Uploading the full stemcell as-is would
# be a no-op, and `version: latest` in the manifest would then silently select the
# light stemcell -- the journey would report success without ever exercising the
# image path it exists to test. Repacking under a name of its own makes the
# stemcell this journey runs on unambiguous.
FULL_STEMCELL_NAME=nvme-journey-full-stemcell
FULL_STEMCELL_VERSION=0.1
repacked_stemcell=$(mktemp -d)/${FULL_STEMCELL_NAME}.tgz

echo "Repacking the full stemcell as ${FULL_STEMCELL_NAME}/${FULL_STEMCELL_VERSION}"
time bosh repack-stemcell \
  --name "${FULL_STEMCELL_NAME}" \
  --version "${FULL_STEMCELL_VERSION}" \
  "$(realpath stemcell/*.tgz)" \
  "${repacked_stemcell}"
time bosh -n upload-stemcell "${repacked_stemcell}"

# is_nvme_instance_type asks the same question the CPI asks when it decides which
# by-id path to hand the agent: DescribeInstanceTypes filtered by
# NvmeSupport=required.
is_nvme_instance_type() {
  local instance_type=$1
  aliyun ecs DescribeInstanceTypes --region "${region}" \
    --InstanceTypes.1 "${instance_type}" --NvmeSupport required 2>/dev/null |
    jq -e '(.InstanceTypes.InstanceType | length) > 0' >/dev/null
}

# Without this check the whole journey can pass while proving nothing: a legacy
# type that turns out to be NVMe-capable, or a target type that is not, makes
# every path assertion below vacuous.
if ! is_nvme_instance_type "${NVME_INSTANCE_TYPE}"; then
  echo "NVME_INSTANCE_TYPE=${NVME_INSTANCE_TYPE} is not NVMe-capable; the journey would assert nothing" >&2
  exit 1
fi
if is_nvme_instance_type "${LEGACY_INSTANCE_TYPE}"; then
  echo "LEGACY_INSTANCE_TYPE=${LEGACY_INSTANCE_TYPE} is NVMe-capable; the journey needs a pre-NVMe starting point" >&2
  exit 1
fi
echo "Preconditions: ${LEGACY_INSTANCE_TYPE} is not NVMe, ${NVME_INSTANCE_TYPE} is"

# The image the director uploaded must carry the NVMe feature, or a 9th-gen VM
# built from it cannot see its disks. Encryption is on for this director, so
# create_stemcell reached this image through CopyImage, where the feature does
# not carry over from the source and the CPI has to re-apply it.
#
# The lookup pins both name and version: picking the first row that merely
# matches a name would go back to reading whichever stemcell the e2e task
# uploaded.
assert_stemcell_is_nvme_capable() {
  local image_id
  refresh_cloud_credentials

  image_id=$(bosh -n stemcells --json |
    jq -er --arg name "${FULL_STEMCELL_NAME}" --arg version "${FULL_STEMCELL_VERSION}" '
      [ .Tables[0].Rows[]
        | select(.name == $name and (.version | sub("\\*$"; "")) == $version) ]
      | if length == 1 then .[0].cid
        else error("expected exactly one \($name)/\($version) stemcell, found \(length)")
        end')

  echo "Checking NVMe support on stemcell image ${image_id}"
  aliyun ecs DescribeImages --region "${region}" --ImageId "${image_id}" |
    jq -e '.Images.Image[0].Features.NvmeSupport == "supported"' >/dev/null
}
assert_stemcell_is_nvme_capable

bosh -n update-cloud-config \
  -l "${METADATA_FILE}" \
  -v "legacy_instance_type=${LEGACY_INSTANCE_TYPE}" \
  -v "nvme_instance_type=${NVME_INSTANCE_TYPE}" \
  "${CLOUD_CONFIG}"

# The journey's own manifest carries no credential properties, so nothing has to
# be passed in here. The stemcell is pinned by version as well as name so the
# deploy cannot drift onto the light stemcell the e2e task uploaded.
deploy_customer_phase() {
  local vm_type=$1
  local disk_type=$2

  time bosh -n deploy -d "${DEPLOYMENT_NAME}" \
    -v "stemcell_name=${FULL_STEMCELL_NAME}" \
    -v "stemcell_version=${FULL_STEMCELL_VERSION}" \
    -v "nvme_vm_type=${vm_type}" \
    -v "nvme_disk_type=${disk_type}" \
    -l "${METADATA_FILE}" \
    "${MANIFEST}"
}

deployment_row() {
  bosh -d "${DEPLOYMENT_NAME}" instances --details --json | \
    jq -ce --arg group "${INSTANCE_GROUP}/" \
      '.Tables[]?.Rows[]? | select(.instance | startswith($group))'
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
  local expected_performance_level=${3:-}
  local current_vm_cid
  local current_disk_cid
  local instance_json
  local disk_json

  refresh_cloud_credentials

  current_vm_cid=$(vm_cid)
  current_disk_cid=$(disk_cid)

  # The credential comes from the environment, which the aliyun CLI reads on its
  # own. Passing it as --access-key-id would put it in the process arguments,
  # where any user on this container can read it.
  instance_json=$(aliyun ecs DescribeInstances \
    --InstanceIds "[\"${current_vm_cid}\"]" \
    --region "${region}")
  echo "${instance_json}" | jq -e \
    --arg expected "${expected_instance_type}" \
    '.Instances.Instance[0].InstanceType == $expected' >/dev/null

  disk_json=$(aliyun ecs DescribeDisks \
    --DiskIds "[\"${current_disk_cid}\"]" \
    --region "${region}")
  echo "${disk_json}" | jq -e \
    --arg expected "${expected_disk_category}" \
    '.Disks.Disk[0].Category == $expected and .Disks.Disk[0].Encrypted == true' >/dev/null

  if [[ -n "${expected_performance_level}" ]]; then
    echo "${disk_json}" | jq -e \
      --arg expected "${expected_performance_level}" \
      '.Disks.Disk[0].PerformanceLevel == $expected' >/dev/null
  fi
}

# assert_device_paths checks the device path the CPI handed to the agent, which a
# mounted-filesystem check cannot see: a regression from the by-id path back to
# /dev/vdb still mounts and still boots, right up until the kernel enumerates the
# disks in a different order. The agent records the CPI's answer verbatim in its
# settings, so that is what gets compared.
#
# create_vm treats a failed resolution as fatal, but attach_disk only warns and
# carries on with a best-effort path, which makes the persistent disk the more
# likely of the two to regress unnoticed.
assert_device_paths() {
  local expected_prefix=$1
  local settings ephemeral persistent

  settings=$(bosh -n -d "${DEPLOYMENT_NAME}" ssh "${NVME_INSTANCE}" -c \
    "sudo cat /var/vcap/bosh/settings.json" --column=stdout --json |
    jq -r '.Tables[0].Rows[0].stdout')

  ephemeral=$(echo "${settings}" | jq -er '.disks.ephemeral')
  # The persistent entry is keyed by disk CID, and is either the path itself or an
  # object carrying it, depending on the agent version.
  persistent=$(echo "${settings}" | jq -er '
    .disks.persistent | to_entries[0].value
    | if type == "object" then (.path // .device_path // empty) else . end')

  echo "  ephemeral=${ephemeral}"
  echo "  persistent=${persistent}"

  local path
  for path in "${ephemeral}" "${persistent}"; do
    case "${path}" in
      "${expected_prefix}"*) ;;
      *)
        echo "expected a device path under ${expected_prefix}, got ${path}" >&2
        return 1
        ;;
    esac
  done

  # The link the CPI reported has to exist on the instance, not merely look
  # plausible: that is the difference between naming a device and finding one.
  bosh -n -d "${DEPLOYMENT_NAME}" ssh "${NVME_INSTANCE}" -c \
    "sudo test -e '${ephemeral}' && sudo test -e '${persistent}'"
}

# assert_disk_kept_in_place distinguishes an in-place update_disk from the
# director falling back to create, copy and attach. Both end with a disk that
# matches the request, so only the CID says which one ran, and that distinction is
# the whole point of the update_disk method.
assert_disk_kept_in_place() {
  local previous=$1 current=$2

  if [[ "${previous}" != "${current}" ]]; then
    echo "expected an in-place change but the disk was replaced: ${previous} -> ${current}" >&2
    return 1
  fi
  echo "  disk kept in place: ${current}"
}

NVME_BY_ID=/dev/disk/by-id/nvme-Alibaba_Cloud_Elastic_Block_Storage_
VIRTIO_BY_ID=/dev/disk/by-id/virtio-

echo "### Phase 1: where a customer stands before moving to 9th-gen"
deploy_customer_phase nvme_upgrade_legacy nvme_upgrade_legacy

marker="nvme-customer-journey-$(date +%s)"
bosh -n -d "${DEPLOYMENT_NAME}" ssh "${NVME_INSTANCE}" -c \
  "sudo mkdir -p '$(dirname "${MARKER_FILE}")' && echo '${marker}' | sudo tee '${MARKER_FILE}' >/dev/null"
assert_remote_state "${marker}"
assert_iaas_state "${LEGACY_INSTANCE_TYPE}" cloud_efficiency
assert_device_paths "${VIRTIO_BY_ID}"
legacy_disk=$(disk_cid)

echo "### Phase 2: move compute and disks to 9th-gen and ESSD together"
deploy_customer_phase nvme_upgrade_target nvme_upgrade_target
assert_remote_state "${marker}"
assert_iaas_state "${NVME_INSTANCE_TYPE}" cloud_essd
assert_device_paths "${NVME_BY_ID}"
upgraded_disk=$(disk_cid)
assert_disk_kept_in_place "${legacy_disk}" "${upgraded_disk}"

echo "### Phase 3: grow the ESSD disk, which update_disk must do in place"
deploy_customer_phase nvme_upgrade_target nvme_upgrade_target_larger
assert_remote_state "${marker}"
assert_iaas_state "${NVME_INSTANCE_TYPE}" cloud_essd
assert_device_paths "${NVME_BY_ID}"
grown_disk=$(disk_cid)
assert_disk_kept_in_place "${upgraded_disk}" "${grown_disk}"

# A shrink is deliberately not exercised here. The CPI rejects it with a plain
# error rather than Bosh::Clouds::NotSupported, so the director does not fall back
# to copy-and-migrate -- the deploy just fails. Shrinking also has no meaning for
# a customer upgrading: the data on the larger disk would not fit.

echo "### Phase 4: recreate on 9th-gen and check the agent returns to both mounts"
time bosh -n -d "${DEPLOYMENT_NAME}" recreate "${NVME_INSTANCE}"
assert_remote_state "${marker}"
assert_iaas_state "${NVME_INSTANCE_TYPE}" cloud_essd
assert_device_paths "${NVME_BY_ID}"
assert_disk_kept_in_place "${grown_disk}" "$(disk_cid)"

echo "### Phase 5: roll compute back to the legacy family, keeping the ESSD disk"
deploy_customer_phase nvme_upgrade_legacy nvme_upgrade_target_larger
assert_remote_state "${marker}"
assert_iaas_state "${LEGACY_INSTANCE_TYPE}" cloud_essd
# The same ESSD disk must now be reported as virtio, because the path depends on
# what the instance supports, not on the disk.
assert_device_paths "${VIRTIO_BY_ID}"
assert_disk_kept_in_place "${grown_disk}" "$(disk_cid)"

echo "NVMe customer journey passed"
