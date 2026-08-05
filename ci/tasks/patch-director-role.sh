#!/usr/bin/env bash

set -euo pipefail

: ${METADATA_FILE:=environment/metadata}
: ${provision_role_arn:?}
: ${region:?}

source bosh-cpi-src/ci/tasks/credentials.sh

ram_role_name=$(python3 -c "import json; print(json.load(open('${METADATA_FILE}'))['ram_role'])")

if [[ -z "${ram_role_name}" ]]; then
  echo "environment metadata has no ram_role" >&2
  exit 1
fi

# Two CPIs are configured here, and they cannot share a credential source:
#
#   /instance_groups/name=bosh is the CPI on the Director. It uses the RAM role
#   that terraform attaches to the Director VM, read from instance metadata and
#   refreshed automatically.
#
#   /cloud_provider is the CPI that `bosh create-env` runs locally on the
#   Concourse worker. It cannot use the worker's own instance role: that role is
#   deliberately limited to sts:AssumeRole, so ecs:CreateInstance is denied. It
#   therefore gets a short-lived credential assumed here. The same credential is
#   refreshed before teardown, because a build runs far longer than a session
#   lasts.
ensure_aliyun_cli
assume_pipeline_role "${provision_role_arn}" "create-env"

echo "Director CPI will run under RAM role: ${ram_role_name}"

cat > bosh-deployment/director-ram-role.yml <<EOF
- path: /resource_pools/name=vms/cloud_properties/ram_role_name?
  type: replace
  value: ${ram_role_name}

- path: /instance_groups/name=bosh/properties/alicloud/credential_source?
  type: replace
  value: ecs_ram_role
- path: /instance_groups/name=bosh/properties/alicloud/access_key_id
  type: remove
- path: /instance_groups/name=bosh/properties/alicloud/access_key_secret
  type: remove

- path: /cloud_provider/properties/alicloud/access_key_id?
  type: replace
  value: ${ALIBABA_CLOUD_ACCESS_KEY_ID}
- path: /cloud_provider/properties/alicloud/access_key_secret?
  type: replace
  value: ${ALIBABA_CLOUD_ACCESS_KEY_SECRET}
- path: /cloud_provider/properties/alicloud/security_token?
  type: replace
  value: ${ALIBABA_CLOUD_SECURITY_TOKEN}
EOF

# The file holds a live credential, so print only its shape.
echo "--- director-ram-role.yml (values redacted) ---"
sed -E 's/^(  value: ).+/\1<redacted>/' bosh-deployment/director-ram-role.yml
