#!/usr/bin/env bash

set -euo pipefail

: ${provision_role_arn:?}
: ${region:?}
: ${DIRECTOR_MANIFEST:=director-state/director.yml}

source bosh-cpi-src/ci/tasks/credentials.sh

# `bosh delete-env` in the teardown task reuses the credential that was baked
# into the director manifest at deploy time. A build runs for hours while an
# assumed-role session lasts an hour, so that credential is expired by then and
# teardown would fail, leaving the director VM behind. Mint a fresh one.
if [[ ! -f "${DIRECTOR_MANIFEST}" ]]; then
  echo "${DIRECTOR_MANIFEST} does not exist, nothing to refresh"
  exit 0
fi

ensure_aliyun_cli
assume_pipeline_role "${provision_role_arn}" "teardown"

python3 - "${DIRECTOR_MANIFEST}" <<'PY'
import os, sys, yaml

path = sys.argv[1]
with open(path) as f:
    manifest = yaml.safe_load(f)

alicloud = manifest.get('cloud_provider', {}).get('properties', {}).get('alicloud')
if alicloud is None:
    raise SystemExit('cloud_provider.properties.alicloud missing from %s' % path)

alicloud['access_key_id'] = os.environ['ALIBABA_CLOUD_ACCESS_KEY_ID']
alicloud['access_key_secret'] = os.environ['ALIBABA_CLOUD_ACCESS_KEY_SECRET']
alicloud['security_token'] = os.environ['ALIBABA_CLOUD_SECURITY_TOKEN']

with open(path, 'w') as f:
    yaml.safe_dump(manifest, f, default_flow_style=False)

print('  refreshed cloud_provider credentials in %s' % path)
PY

# Fail loudly rather than let teardown discover a broken manifest.
python3 -c "
import yaml, sys
m = yaml.safe_load(open('${DIRECTOR_MANIFEST}'))
a = m['cloud_provider']['properties']['alicloud']
missing = [k for k in ('access_key_id', 'access_key_secret', 'security_token') if not a.get(k)]
if missing:
    sys.exit('refresh left these empty: %s' % missing)
print('  manifest still parses and carries a complete credential')
"
