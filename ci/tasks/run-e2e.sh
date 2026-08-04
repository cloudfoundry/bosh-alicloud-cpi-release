#!/usr/bin/env bash

set -euo pipefail

source bosh-cpi-src/ci/tasks/utils.sh
source director-state/director.env

CURRENT_PATH=$(pwd)
JQ_BLOB_PATH=$CURRENT_PATH/jq-blob
ALIYUN_CLI_PATH=$CURRENT_PATH/aliyun-cli
METADATA_FILE=$CURRENT_PATH/environment/metadata

: "${ALICLOUD_ACCESS_KEY_ID:?}"
: "${ALICLOUD_SECRET_ACCESS_KEY:?}"
: "${LEGACY_INSTANCE_TYPE:?}"
: "${NVME_INSTANCE_TYPE:?}"

# add blobs for e2e test
# when e2e-test-release/config/final set to remote storage, this code can be removed
function addBlobs {
   bosh add-blob ${JQ_BLOB_PATH}/jq-linux64 jq/jq-linux64
   bosh add-blob ${ALIYUN_CLI_PATH}/aliyun-cli-linux-*.tgz aliyuncli/aliyun-cli-linux-amd64.tgz
}

# CREATE TEST RELEASE
pushd bosh-cpi-src/ci/assets/e2e-test-release
  addBlobs
  time bosh -n create-release --force --name e2e-test --version 1.0.0
  time bosh -n upload-release
popd

# UPLOAD STEMCELL
echo "upload light stemcell..."
time bosh -n upload-stemcell "$(realpath stemcell/*.tgz)"
#echo "upload heavy stemcell..."
#time bosh -n upload-stemcell "$(realpath heavy-stemcell/*.tgz)"

stemcell_name="$( bosh int <( tar xfO $(realpath stemcell/*.tgz) stemcell.MF ) --path /name )"
#heavy_stemcell_name="$( bosh int <( tar xfO $(realpath heavy-stemcell/*.tgz) stemcell.MF ) --path /name )"
# alicloud_kms_key_arn="$(cat environment/metadata | jq --raw-output ".alicloud_kms_key_arn")"

# --cloud-properties "{\"encrypted\": true, \"kms_key_arn\": \"${alicloud_kms_key_arn}\"}" \
#time bosh repack-stemcell \
#  --name e2e-encrypted-heavy-stemcell \
#  --version 0.1 \
#  "$(realpath heavy-stemcell/*.tgz)" \
#  /tmp/e2e-encrypted-heavy-stemcell.tgz
#time bosh -n upload-stemcell /tmp/e2e-encrypted-heavy-stemcell.tgz
#encrypted_heavy_stemcell_img_id="$( bosh stemcells | grep e2e-encrypted-heavy-stemcell | awk '{print $NF;}' )"

#UPDATE CLOUD CONFIG
time bosh -n ucc \
  -l ${METADATA_FILE} \
  -v "legacy_instance_type=${LEGACY_INSTANCE_TYPE}" \
  -v "nvme_instance_type=${NVME_INSTANCE_TYPE}" \
  bosh-cpi-src/ci/assets/e2e-test-release/cloud-config.yml

# BOSH DEPLOY
time bosh -n deploy -d e2e-test \
  -v "stemcell_name=${stemcell_name}" \
  -v "nvme_vm_type=nvme_upgrade_legacy" \
  -v "nvme_disk_type=nvme_upgrade_legacy" \
  -v "access_key=${ALICLOUD_ACCESS_KEY_ID}" \
  -v "secret_key=${ALICLOUD_SECRET_ACCESS_KEY}" \
  -l "${METADATA_FILE}" \
  bosh-cpi-src/ci/assets/e2e-test-release/manifest.yml

# RUN ERRANDS
#time bosh -n run-errand -d e2e-test iam-instance-profile-test
time bosh -n run-errand -d e2e-test spot-instance-test
time bosh -n run-errand -d e2e-test slb-registration-test
#time bosh -n run-errand -d e2e-test heavy-stemcell-test
#time bosh -n run-errand -d e2e-test encrypted-heavy-stemcell-test
