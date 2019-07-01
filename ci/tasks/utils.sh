#!/usr/bin/env bash

# Oportunistically configure bosh for use
configure_bosh_cli() {
  local bosh_input="$(realpath bosh-cli/*bosh-cli-* 2>/dev/null || true)"
  if [[ -n "${bosh_input}" ]]; then
    export bosh_cli="/usr/local/bin/bosh"
    cp "${bosh_input}" "${bosh_cli}"
    chmod +x "${bosh_cli}"
  fi
}
configure_bosh_cli

# Oportunistically configure aliyun cli for use
configure_aliyun_cli() {
  local cli_input="$(realpath aliyun-cli/aliyun-cli-* 2>/dev/null || true)"
  if [[ -n "${cli_input}" ]]; then
    tar -xzf aliyun-cli/aliyun-cli-linux-*.tgz -C /usr/bin
  fi
}
configure_aliyun_cli

check_param() {
  local name=$1
  local value=$(eval echo '$'$name)
  if [ "$value" == 'replace-me' ]; then
    echo "environment variable $name must be set"
    exit 1
  fi
}

print_git_state() {
  echo "--> last commit..."
  TERM=xterm-256color git log -1
  echo "---"
  echo "--> local changes (e.g., from 'fly execute')..."
  TERM=xterm-256color git status --verbose
  echo "---"
}

declare -a on_exit_items
on_exit_items=()

function on_exit {
  echo "Running ${#on_exit_items[@]} on_exit items..."
  for i in "${on_exit_items[@]}"
  do
    for try in $(seq 0 9); do
      sleep $try
      echo "Running cleanup command $i (try: ${try})"
        eval $i || continue
      break
    done
  done
}

function add_on_exit {
  local n=${#on_exit_items[@]}
  on_exit_items=("${on_exit_items[@]}" "$*")
  if [[ $n -eq 0 ]]; then
    trap on_exit EXIT
  fi
}


function check_go_version {
  local cpi_release=$1
  local release_go_version="$(cat "$cpi_release/packages/golang/spec" | \
   grep linux | awk '{print $2}' | sed "s/golang\/go\(.*\)\.linux-amd64.tar.gz/\1/")"

  local current=$(go version)
  if [[ "$current" != *"$release_go_version"* ]]; then
    echo "Go version is incorrect. Current version: $current, Required version: $release_go_version"
    # Go version is incorrect. Current version: go version go1.8.1 linux/amd64, Required version: go1.8.1.linux-amd64.tar.gz
    #exit 1
  fi
}

# cidr to mask
# example:
#   cdr2mask 18 => 255.255.192
function cdr2mask {
   # Number of args to shift, 255..255, first non-255 byte, zeroes
   set -- $(( 5 - ($1 / 8) )) 255 255 255 255 $(( (255 << (8 - ($1 % 8))) & 255 )) 0 0 0
   [ $1 -gt 1 ] && shift $1 || shift
   echo ${1-0}.${2-0}.${3-0}.${4-0}
}

function getCidrNotation {
 echo $1 | awk -F / '{print $2}'
}
