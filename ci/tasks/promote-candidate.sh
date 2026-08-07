#!/usr/bin/env bash

set -e

source bosh-cpi-src/ci/tasks/utils.sh

integer_version=`cut -f1 release-version-semver/number`
echo $integer_version > promoted/integer_version

cp -r bosh-cpi-src promoted/repo

dev_release=$(echo $PWD/bosh-cpi-release/*.tgz)

pushd promoted/repo
  # No config/private.yml is written: config/final.yml uses the local blobstore,
  # so blobs are committed to the repo and finalize-release needs no credentials.

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

  git diff | cat
  git add .

  git config --global user.email kangtai.kang@alibaba-inc.com
  git config --global user.name ali-tkang
  git commit -m "Bump Alibaba Cloud cpi/$integer_version"
popd