#!/usr/bin/env bash

set -uo pipefail

: ${provision_role_arn:?}
: ${region:?}
: ${DIRECTOR_MANIFEST:=director-state/director.yml}

source bosh-cpi-src/ci/tasks/credentials.sh

# `bosh delete-env` in the teardown task reuses the credential that was baked
# into the director manifest at deploy time. A build runs for hours while an
# assumed-role session lasts an hour, so that credential is expired by then and
# teardown would fail, leaving the director VM behind. Mint a fresh one.
#
# This runs first in an `ensure` block, so it must never fail: a non-zero exit
# would skip teardown and ensure-terminated and leak the whole environment. Warn
# and let the cleanup chain continue instead.
warn_and_exit() {
  echo "WARNING: $1; teardown will use the credential from deploy time" >&2
  exit 0
}

if [[ ! -f "${DIRECTOR_MANIFEST}" ]]; then
  echo "${DIRECTOR_MANIFEST} does not exist, nothing to refresh"
  exit 0
fi

command -v bosh >/dev/null 2>&1 || warn_and_exit "the bosh CLI is not on PATH"

ensure_aliyun_cli || warn_and_exit "could not install the aliyun CLI"
assume_pipeline_role "${provision_role_arn}" "teardown" || warn_and_exit "AssumeRole failed"

ops=$(mktemp)
cat > "${ops}" <<EOF
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

# bosh int keeps the manifest valid; the image's python has no yaml module.
refreshed=$(mktemp)
if ! bosh int "${DIRECTOR_MANIFEST}" -o "${ops}" > "${refreshed}" 2>/tmp/bosh-int-err; then
  rm -f "${ops}"
  echo "  bosh int failed: $(head -3 /tmp/bosh-int-err | tr '\n' ' ')" >&2
  warn_and_exit "could not rewrite ${DIRECTOR_MANIFEST}"
fi

if [[ ! -s "${refreshed}" ]]; then
  rm -f "${ops}" "${refreshed}"
  warn_and_exit "bosh int produced an empty manifest"
fi

mv "${refreshed}" "${DIRECTOR_MANIFEST}"
rm -f "${ops}"
echo "  refreshed the cloud_provider credential in ${DIRECTOR_MANIFEST}"
