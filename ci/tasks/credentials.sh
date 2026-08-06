#!/usr/bin/env bash
#
# Short-lived credential bootstrap for the CPI pipeline.
#
# No pipeline task receives a long-lived access key. Every task starts from the
# RAM role attached to the Concourse worker ECS instance, assumes the role for
# the work it is about to do, and exports only the resulting short-lived STS
# credential. Cleanup tasks assume their role independently, so they still work
# when the task that created the environment failed.
#
# Source this file; do not execute it.

# ecs_metadata_base is the link-local instance metadata service. Overridable
# only so the pipeline can be exercised against a stub.
: ${ecs_metadata_root:=http://100.100.100.200}
: ${ecs_metadata_base:=${ecs_metadata_root}/latest/meta-data}
: ${ecs_metadata_token_ttl:=60}

# assume_role_duration_seconds covers a full create/test/destroy leg. The role's
# MaxSessionDuration must be at least this large.
: ${assume_role_duration_seconds:=3600}

# json_field extracts a top-level-ish string field from a JSON document without
# requiring jq, which is not present in every task image.
json_field() {
  local field="$1" document="$2"

  if command -v jq >/dev/null 2>&1; then
    printf '%s' "${document}" | jq -r --arg f "${field}" '..|objects|select(has($f))|.[$f]' 2>/dev/null | head -1
    return 0
  fi

  printf '%s' "${document}" |
    tr ',{}' '\n\n\n' |
    sed -n "s/.*\"${field}\"[[:space:]]*:[[:space:]]*\"\([^\"]*\)\".*/\1/p" |
    head -1
}

# require_credential_value fails when a credential field came back empty. The
# value itself is never printed.
require_credential_value() {
  local name="$1" value="$2"
  if [[ -z "${value}" || "${value}" == "null" ]]; then
    echo "Credential bootstrap failed: ${name} is empty" >&2
    return 1
  fi
}

# metadata_token fetches an IMDSv2 token. An instance configured to require
# IMDSv2 rejects unauthenticated reads, so the token is obtained first. An
# instance that still allows IMDSv1 does not serve this endpoint, in which case
# this prints nothing and the caller reads metadata unauthenticated.
metadata_token() {
  curl -s --max-time 5 -X PUT \
    -H "X-aliyun-ecs-metadata-token-ttl-seconds: ${ecs_metadata_token_ttl}" \
    "${ecs_metadata_root}/latest/api/token" 2>/dev/null || true
}

# read_metadata GETs a metadata path, presenting an IMDSv2 token when one is
# available.
read_metadata() {
  local path="$1" token
  token=$(metadata_token)

  if [[ -n "${token}" ]]; then
    curl -s --max-time 5 -H "X-aliyun-ecs-metadata-token: ${token}" "${ecs_metadata_base}/${path}"
  else
    curl -s --max-time 5 "${ecs_metadata_base}/${path}"
  fi
}

# worker_role_name returns the RAM role attached to this ECS instance.
worker_role_name() {
  local role_name
  role_name=$(read_metadata "ram/security-credentials/" | head -1 | tr -d '[:space:]')

  if [[ -z "${role_name}" ]]; then
    cat >&2 <<'EOF'
Credential bootstrap failed: no RAM role is attached to this Concourse worker.

Attach the bootstrap role to the worker ECS instance and make sure task
containers can reach 100.100.100.200.
EOF
    return 1
  fi

  printf '%s' "${role_name}"
}

# load_worker_role_credentials puts the worker role's STS credential into the
# environment so the aliyun CLI can call AssumeRole with it.
load_worker_role_credentials() {
  local role_name document
  role_name=$(worker_role_name) || return 1

  echo "Using Concourse worker RAM role: ${role_name}"

  document=$(read_metadata "ram/security-credentials/${role_name}")

  local code
  code=$(json_field Code "${document}")
  if [[ "${code}" != "Success" ]]; then
    # The document holds a usable credential, so report only the status field.
    echo "Credential bootstrap failed: instance metadata returned Code='${code:-<none>}'" >&2
    return 1
  fi

  export ALIBABA_CLOUD_ACCESS_KEY_ID="$(json_field AccessKeyId "${document}")"
  export ALIBABA_CLOUD_ACCESS_KEY_SECRET="$(json_field AccessKeySecret "${document}")"
  export ALIBABA_CLOUD_SECURITY_TOKEN="$(json_field SecurityToken "${document}")"

  require_credential_value "worker AccessKeyId" "${ALIBABA_CLOUD_ACCESS_KEY_ID}" || return 1
  require_credential_value "worker AccessKeySecret" "${ALIBABA_CLOUD_ACCESS_KEY_SECRET}" || return 1
  require_credential_value "worker SecurityToken" "${ALIBABA_CLOUD_SECURITY_TOKEN}" || return 1
}

# build_session_name derives an ActionTrail session name from the Concourse
# build so a cloud API call can be traced back to the build that made it.
build_session_name() {
  local prefix="$1"
  local raw="${prefix}-${BUILD_PIPELINE_NAME:-local}-${BUILD_JOB_NAME:-task}-${BUILD_ID:-0}"

  # Session names allow [a-zA-Z0-9.@-_] and at most 64 characters.
  printf '%s' "${raw}" | tr -c 'a-zA-Z0-9.@_-' '-' | cut -c1-64
}

# ensure_aliyun_cli installs the CLI when the task image does not ship it. The
# credential bootstrap needs it to call AssumeRole, and it is fetched from
# Alibaba Cloud's public CDN so no credential is required to get it.
ensure_aliyun_cli() {
  if command -v aliyun >/dev/null 2>&1; then
    return 0
  fi
  local dir="$(mktemp -d)"
  wget -qO "${dir}/aliyun-cli.tgz" https://aliyuncli.alicdn.com/aliyun-cli-linux-latest-amd64.tgz
  tar -zxf "${dir}/aliyun-cli.tgz" -C /usr/bin
  rm -rf "${dir}"
}

# assume_pipeline_role exchanges the worker role credential for the role that
# owns the work about to be done, and exports the result under every variable
# name the aliyun CLI and the Terraform alicloud provider understand.
#
# Usage: assume_pipeline_role <role-arn> <session-prefix>
assume_pipeline_role() {
  local role_arn="$1" session_prefix="$2"

  if [[ -z "${role_arn}" ]]; then
    echo "assume_pipeline_role needs a role ARN" >&2
    return 1
  fi

  clear_cloud_credentials
  load_worker_role_credentials || return 1

  local session_name
  session_name=$(build_session_name "${session_prefix}")

  echo "Assuming ${role_arn} as session ${session_name}"

  local response
  # The response carries a live credential, so it is captured and never echoed.
  response=$(aliyun sts AssumeRole \
    --RoleArn "${role_arn}" \
    --RoleSessionName "${session_name}" \
    --DurationSeconds "${assume_role_duration_seconds}" \
    --region "${region:-${ALIBABA_CLOUD_REGION_ID}}" 2>&1) || {
    echo "AssumeRole ${role_arn} failed. Check the role's trust policy and MaxSessionDuration." >&2
    # Surface only the error code; the body may contain credential material.
    printf '%s\n' "${response}" | sed -n 's/.*\(ErrorCode: [A-Za-z.]*\).*/\1/p' >&2
    return 1
  }

  local access_key_id access_key_secret security_token
  access_key_id=$(json_field AccessKeyId "${response}")
  access_key_secret=$(json_field AccessKeySecret "${response}")
  security_token=$(json_field SecurityToken "${response}")

  require_credential_value "assumed AccessKeyId" "${access_key_id}" || return 1
  require_credential_value "assumed AccessKeySecret" "${access_key_secret}" || return 1
  require_credential_value "assumed SecurityToken" "${security_token}" || return 1

  # aliyun CLI and recent versions of the Terraform provider.
  export ALIBABA_CLOUD_ACCESS_KEY_ID="${access_key_id}"
  export ALIBABA_CLOUD_ACCESS_KEY_SECRET="${access_key_secret}"
  export ALIBABA_CLOUD_SECURITY_TOKEN="${security_token}"

  # Terraform alicloud provider and its OSS backend.
  export ALICLOUD_ACCESS_KEY="${access_key_id}"
  export ALICLOUD_SECRET_KEY="${access_key_secret}"
  export ALICLOUD_SECURITY_TOKEN="${security_token}"

  echo "Obtained short-lived credentials valid for ${assume_role_duration_seconds}s"
}

# clear_cloud_credentials removes every credential variable from the
# environment. Call it before handing control to a step that must not inherit
# credentials, such as a CPI running under its own instance role.
clear_cloud_credentials() {
  unset ALIBABA_CLOUD_ACCESS_KEY_ID ALIBABA_CLOUD_ACCESS_KEY_SECRET ALIBABA_CLOUD_SECURITY_TOKEN
  unset ALIBABACLOUD_ACCESS_KEY_ID ALIBABACLOUD_ACCESS_KEY_SECRET ALIBABACLOUD_SECURITY_TOKEN
  unset ALICLOUD_ACCESS_KEY ALICLOUD_SECRET_KEY ALICLOUD_SECURITY_TOKEN
  unset ALICLOUD_ACCESS_KEY_ID ALICLOUD_ACCESS_KEY_SECRET
  unset ACCESS_KEY_ID ACCESS_KEY_SECRET SECURITY_TOKEN
}

# assert_no_credentials_in_file fails when a file that gets published as a task
# output or stored remotely contains credential material.
assert_no_credentials_in_file() {
  local path="$1"
  [[ -f "${path}" ]] || return 0

  if grep -qE '(LTAI[A-Za-z0-9]{6,}|"?[Aa]ccess[Kk]ey[Ss]ecret"?[[:space:]]*[:=]|"?[Ss]ecurity[Tt]oken"?[[:space:]]*[:=])' "${path}"; then
    echo "Refusing to continue: ${path} contains credential material" >&2
    return 1
  fi
}
