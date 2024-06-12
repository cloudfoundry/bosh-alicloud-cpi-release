# BOSH Alibaba Cloud CPI

This is a [BOSH](http://bosh.io) release for the BOSH Alibaba Cloud CPI.

**[Cloud Foundry Install Guide](docs/cf/install-cf.md)**

### How to make release

- Clone this repo
- Install golang and export $GOROOT
- Install bosh-cli
- Download `go1.20.10.linux-amd64.tar.gz` from https://storage.googleapis.com/golang/go1.20.10.linux-amd64.tar.gz
- Add blob and create release
- Configure golang env and make
- Create bosh release

```
$ git clone https://github.com/aliyun/bosh-alicloud-cpi-release.git
$ cd bosh-alicloud-cpi-release
$ mkdir blobs
$ bosh add-blob ~/Downloads/go1.20.10.linux-amd64.tar.gz go1.20.10.linux-amd64.tar.gz
$ source .envrc
$ make
$ bosh create-release --force --tarball=../bosh-alicloud-cpi.tgz
```
### Latest Version (v18)

At present, alibaba cloud the lastet cpi release is stored in the OSS. You can reference this release in your deployment manifest from the releases section:
```
- name: bosh-alicloud-cpi
  version: 18
  url: http://bosh-((region)).oss-((region)).aliyuncs.com/bosh-alicloud-cpi-release-r18.tgz
  sha1: 2e36b60f99e51fe2dd4e7ea021bf3ecce883487c
```
Or upload it to your director with the upload-release command:
```
bosh upload-release --sha1 2e36b60f99e51fe2dd4e7ea021bf3ecce883487c \
  http://bosh-((region)).oss-((region)).aliyuncs.com/bosh-alicloud-cpi-release-r18.tgz
```

**-> Note:** If region is in China mainland, its download url is unified as: `http://bosh.oss-cn-hangzhou.aliyuncs.com/bosh-alicloud-cpi-release-r18.tgz`.

## Usage

### Prepare your `Alibaba Cloud` environment

- Create a vpc with switch and get `vswtich_id`
- Create security group get `security_group_id`
- Create user access key, get `access_key_id/access_key_secret`
- Create a key pair, get `key_pair_name` and download it private key, like bosh.pem
- Create a jumpbox vm

### Install bosh in Alibaba Cloud

- Clone [bosh-deployment](https://github.com/aliyun/bosh-deployment) repo from github

```
$ git clone https://github.com/aliyun/bosh-deployment.git
$ cd bosh-deployment
```

use this command, modify the parameters

- vswitch_id
- security_group_id
- access_key_id
- access_key_secret
- region
- zone
- key_pair_name
- private_key

```
bosh create-env bosh-deployment/bosh.yml --state=state.json \
 --vars-store=creds.yml \
 -o bosh-deployment/alicloud/cpi.yml \
 -v director_name=my-bosh \
 -v internal_cidr=192.168.0.0/24 \
 -v internal_gw=192.168.0.1 \
 -v internal_ip=192.168.0.2 \
 -v vswitch_id=... \
 -v security_group_id=... \
 -v access_key_id=... \
 -v access_key_secret=... \
 -v region=cn-beijing \
 -v zone=cn-beijing-a \
 -v key_pair_name=... \
 -v private_key=bosh.pem
```

### Install with jumpbox-user

```
$ bosh create-env bosh-deployment/bosh.yml \
 --state=state.json \
 --vars-store=creds.yml \
 -o bosh-deployment/alicloud/cpi.yml \
 -o bosh-deployment/jumpbox-user.yml \
 -v director_name=my-bosh \
 -v internal_cidr=192.168.0.0/24 \
 -v internal_gw=192.168.0.1 \
 -v internal_ip=192.168.0.7 \
 -v vswitch_id=... \
 -v security_group_id=... \
 -v access_key_id=... \
 -v access_key_secret=... \
 -v region=cn-beijing \
 -v zone=cn-beijing-e
 -v key_pair_name=... \
 -v private_key=bosh.pem
```

```
$ bosh int creds.yml --path /jumpbox_ssh/private_key > jumpbox.key
$ chmod 600 jumpbox.key
$ ssh jumpbox@<external-or-internal-ip> -i jumpbox.key
```

### Install with external-ip

- Create a `Elastic IP` in `Alibaba Cloud Console`, get an `external_ip`

```
bosh create-env bosh-deployment/bosh.yml \
 --state=state.json \
 --vars-store=creds.yml \
 -o bosh-deployment/alicloud/cpi.yml \
 -o bosh-deployment/external-ip-not-recommended.yml \
 -o bosh-deployment/jumpbox-user.yml \
 -o bosh-deployment/misc/powerdns.yml \
 -v director_name=my-bosh  \
 -v internal_cidr=192.168.0.0/24 \
 -v internal_gw=192.168.0.1 \
 -v internal_ip=192.168.0.7 \
 -v vswitch_id=... \
 -v security_group_id=... \
 -v access_key_id=... \
 -v access_key_secret=... \
 -v region=cn-beijing \
 -v zone=cn-beijing-e \
 -v key_pair_name=... \
 -v private_key=bosh.pem \
 -v external_ip=... \
 -v dns_recursor_ip=8.8.8.8
```

### Running ERB job templates unit tests

The ERB templates rendered by the jobs of this Bosh Release have unit tests
using Ruby. The required Ruby version is specified in `.ruby-version` as per
convention with `chruby` or similar tools. A script will help you to install
the correct Ruby version if necessary and run the ERB unit tests:

```
./scripts/test-unit-erb
```

### Run Unit Test

Run following commands

```
$ source .envrc
$ make testdeps
$ ginkgo -r -skipPackage integration src/bosh-alicloud-cpi
```

### Run Integration Test

Prepare your `Alibaba Cloud` environment, and export follow variables

- Create a SLB (Load balance) get `slb_id`

```
export CPI_REGION=cn-beijing
export CPI_ZONE=cn-beijing-e
export CPI_ACCESS_KEY_ID=...
export CPI_ACCESS_KEY_SECRET=...
export CPI_SECURITY_GROUP_ID=...
export CPI_VSWITCH_ID=...
export CPI_STEMCELL_ID=...
export CPI_INTERNAL_CIDR=192.168.0.0/24/
export CPI_INTERNAL_NETMASK=255.255.255.0
export CPI_INTERNAL_IP=192.168.0.2
export CPI_INTERNAL_GW=192.168.0.1
export CPI_EXTERNAL_IP=47.47.47.47
export CPI_SLB_ID=...
```

Go to source code path, run follow commands

```
$ make testdeps
$ ginkgo -r src/bosh-alicloud-cpi/integration
```
