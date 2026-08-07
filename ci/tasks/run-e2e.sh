#!/usr/bin/env bash

set -e

source bosh-cpi-src/ci/tasks/utils.sh
source bosh-cpi-src/ci/tasks/credentials.sh
source director-state/director.env

: ${observer_role_arn:?}
: ${region:?}

CURRENT_PATH=$(pwd)
JQ_BLOB_PATH=$CURRENT_PATH/jq-blob
ALIYUN_CLI_PATH=$CURRENT_PATH/aliyun-cli
METADATA_FILE=$CURRENT_PATH/environment/metadata

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
  bosh-cpi-src/ci/assets/e2e-test-release/cloud-config.yml

# BOSH DEPLOY
#-v "heavy_stemcell_name=${heavy_stemcell_name}" \
#-v "encrypted_heavy_stemcell_img_id=${encrypted_heavy_stemcell_img_id}" \
# The errands query ECS and SLB to assert what the CPI did, so they need a
# credential of their own. They get a short-lived assumed-role one rather than a
# RAM role on their VM, which would require the director's CPI to hold
# ram:PassRole.
#
# observer_role_arn is read-only. The errands only describe resources, so they
# must not receive the provisioning role, which can create and delete the whole
# test environment.
assume_pipeline_role "${observer_role_arn}" "e2e-observer"

# The credential goes through a vars file rather than -v, which would put it in
# the process arguments where any user on this container can read it.
deploy_vars=$(mktemp)
chmod 0600 "${deploy_vars}"
trap 'rm -f "${deploy_vars}"' EXIT

cat > "${deploy_vars}" <<EOF
access_key: ${ALIBABA_CLOUD_ACCESS_KEY_ID}
secret_key: ${ALIBABA_CLOUD_ACCESS_KEY_SECRET}
security_token: ${ALIBABA_CLOUD_SECURITY_TOKEN}
EOF

time bosh -n deploy -d e2e-test \
  -v "stemcell_name=${stemcell_name}" \
  -l "${deploy_vars}" \
  -l ${METADATA_FILE} \
  bosh-cpi-src/ci/assets/e2e-test-release/manifest.yml

rm -f "${deploy_vars}"
trap - EXIT

# RUN ERRANDS
#time bosh -n run-errand -d e2e-test iam-instance-profile-test
time bosh -n run-errand -d e2e-test spot-instance-test
time bosh -n run-errand -d e2e-test slb-registration-test
#time bosh -n run-errand -d e2e-test heavy-stemcell-test
#time bosh -n run-errand -d e2e-test encrypted-heavy-stemcell-test