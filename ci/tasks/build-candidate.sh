#!/usr/bin/env bash

set -e


: ${GIT_USER_EMAIL:?}
: ${GIT_USER_NAME:?}
: ${GIT_USER_ID:?}
: ${GIT_USER_PASSWORD:?}
: ${BOSH_REPO_HOST:?}
: ${BOSH_REPO_BRANCH:?}

source bosh-cpi-src/ci/tasks/utils.sh
source /etc/profile.d/chruby-with-ruby-2.1.2.sh

cpi_release_name="bosh-alicloud-cpi"
semver=`cat version-semver/number`

# install bosh
echo "Installing Bosh CLI..."
curl -O https://s3.amazonaws.com/bosh-cli-artifacts/bosh-cli-2.0.28-linux-amd64
chmod +x ./bosh-cli-*
mv ./bosh-cli-* /usr/local/bin/bosh2

cp -r bosh-cpi-src candidate/repo

pushd candidate/repo
  echo "******** using eof *********"
  echo git fetch https://${GIT_USER_ID}@${BOSH_REPO_HOST} ${BOSH_REPO_BRANCH} \<\< EOF > bosh_login.sh
  echo ${GIT_USER_PASSWORD} >> bosh_login.sh
  echo EOF >> bosh_login.sh
  chmod +x bosh_login.sh
  ./bosh_login.sh > log1

  cat log1

  echo "******** git install expect ********"
  sudo apt-get install expect -y

  echo "******** git pull by https ********"
  echo "#!/usr/bin/expect" > git_install.sh
  echo "spawn git fetch https://${GIT_USER_ID}@${BOSH_REPO_HOST} ${BOSH_REPO_BRANCH}" >> git_install.sh
  echo "expect \"Password for 'https://${GIT_USER_ID}@gitlab.com': \"" >> git_install.sh
  echo "send \"${GIT_USER_PASSWORD}\"" >> git_install.sh
  echo "expect eof" >> git_install.sh
  echo exit >> git_install.sh
  chmod +x git_install.sh
  ./git_install.sh > log2
  cat log2
  rm -rf ./git_install.sh

  echo $'\n'
  echo "****** git merge ******"
  git merge FETCH_HEAD

  echo "using bosh CLI version..."
  bosh2 -v

  cpi_release_name="bosh-alicloud-cpi"

  # fix cannot find package "bosh-alicloud-cpi/action
  source .envrc
  #echo $GOPATH

  make
  
  # add go cpi blob
  #ls ../../go-cpi-blobs
  bosh2 add-blob ../../go-cpi-blobs/go1.8.1.linux-amd64.tar.gz go1.8.1.linux-amd64.tar.gz

  # export TERM=msys
  echo "git status..."
  git status

  # todo: get email and user from params
  git config --global user.email ${GIT_USER_EMAIL}
  git config --global user.name ${GIT_USER_NAME}

  git add .
  git commit -m 'do nothing'
  #git pull

  ls -al

  echo "building CPI release..."
  # refers: https://bosh.io/docs/cli-v2#create-release
  bosh2 create-release --name $cpi_release_name --version $semver --tarball $cpi_release_name-$semver.tgz

  ls -al dev-release-artifacts
  mkdir dev-release-artifacts
  rm -rf dev-release-artifacts/*
  mv $cpi_release_name-$semver.tgz dev-release-artifacts/
  ls -al dev-release-artifacts

  echo "git status..."
  git status

  git add .
  git commit -m 'create cpi release'
popd