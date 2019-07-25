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

  echo "finalizing CPI release..."
  bosh finalize-release ${dev_release} --version $integer_version --force

  rm config/private.yml

  git diff | cat
  git add .

  git config --global user.email guimin.hgm@alibaba-inc.com
  git config --global user.name xiaozhu36
  git commit -m "New final release v $integer_version"
popd