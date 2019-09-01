#!/usr/bin/env bash

set -e

: ${ALICLOUD_ACCESS_KEY_ID:?}
: ${ALICLOUD_ACCESS_KEY_SECRET:?}

source bosh-cpi-src/ci/tasks/utils.sh
source /etc/profile.d/chruby.sh
chruby 2.4.4

integer_version=`cut -f1 release-version-semver/number`
echo $integer_version > promoted/integer_version

cp -r bosh-cpi-src promoted/repo

dev_release=$(echo $PWD/bosh-cpi-release/*.tgz)

pushd promoted/repo
  echo "Creating config/private.yml with blobstore secrets"
  set +x
  cat > config/private.yml << EOF
---
blobstore:
  s3:
    access_key_id: $ALICLOUD_ACCESS_KEY_ID
    secret_access_key: $ALICLOUD_ACCESS_KEY_SECRET
EOF

  # update the changelog
  changelog="CHANGELOG.md"
  line_num=`cat -n $changelog |grep "(Unreleased)"|awk '{print $1}'`
  bump_date=`env LANG=en_US.UTF-8 date '+%B %d, %Y'`
  sed -i -e "s/(Unreleased)/($bump_date)/" $changelog

  arr=(${integer_version//./ })
  next_version="$((${arr[0]} + 1)).0.0"
  sed -i -e "${line_num}i \#\# $next_version (Unreleased)" $changelog

  echo "finalizing CPI release..."
  bosh finalize-release ${dev_release} --version $integer_version --force

  rm config/private.yml

  git diff | cat
  git add .

  git config --global user.email guimin.hgm@alibaba-inc.com
  git config --global user.name xiaozhu36
  git commit -m "Bump Alibaba Cloud cpi/$integer_version"
popd