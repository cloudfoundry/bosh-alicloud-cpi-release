#!/usr/bin/env bash

set -e

source bosh-cpi-src/ci/tasks/utils.sh

cpi_release_name="bosh-alicloud-cpi"
semver=`cat version-semver/number`
image_path=$PWD/bosh-cpi-src/${cpi_release_name}-${semver}.tgz

pushd bosh-cpi-src
  echo "Using Bosh CLI version..."
  bosh -v

  source .envrc

  echo "Finding Golang Release and downloading it..."
  cat config/blobs.yml | while read LINE
    do
        if [[ ${LINE//:/} =~ ^go[0-9.]+linux-[a-z0-9]+.tar.gz$ ]]; then
            gorelease=${LINE//:/}
            echo "Downloading ${gorelease}..."
            wget https://go.dev/dl/${gorelease}
            echo "Adding ${gorelease} to blob..."
            bosh add-blob ./${gorelease} ${gorelease}
            break
        fi
    done

  echo "Exposing release semver to bosh-alicloud-cpi"
  echo ${semver} > "src/bosh-alicloud-cpi/release"

  # We have to use the --force flag because we just added the `src/bosh-alicloud-cpi/release` file
  echo "Creating CPI BOSH Release..."
  bosh create-release --name=${cpi_release_name} --version=${semver} --tarball=${image_path} --force
popd

echo -n $(sha1sum $image_path | awk '{print $1}') > $image_path.sha1

mv ${image_path} candidate/
mv ${image_path}.sha1 candidate/