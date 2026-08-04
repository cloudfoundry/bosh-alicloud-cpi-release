#!/usr/bin/env bash

set -euo pipefail

: ${METADATA_FILE:=environment/metadata}

ram_role_name=$(python3 -c "import json; print(json.load(open('${METADATA_FILE}'))['ram_role'])")

if [[ -z "${ram_role_name}" ]]; then
  echo "environment metadata has no ram_role; terraform must create the director role" >&2
  exit 1
fi

echo "Configuring the director to run the CPI under RAM role: ${ram_role_name}"

# Applied after bosh-deployment/alicloud/cpi.yml, so it can drop the access key
# that cpi.yml puts in the manifest.
#
# Two CPIs are configured here:
#   - /cloud_provider is the CPI that `bosh create-env` runs locally, which picks
#     up the Concourse worker's own role by leaving ram_role_name unset.
#   - /instance_groups/name=bosh is the CPI on the Director, which gets the role
#     that terraform attaches to the Director VM.
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

- path: /cloud_provider/properties/alicloud/credential_source?
  type: replace
  value: ecs_ram_role
- path: /cloud_provider/properties/alicloud/access_key_id
  type: remove
- path: /cloud_provider/properties/alicloud/access_key_secret
  type: remove
EOF

echo "--- director-ram-role.yml ---"
cat bosh-deployment/director-ram-role.yml
