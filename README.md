# BOSH Alibaba Cloud CPI

This is a [BOSH](http://bosh.io) release for the BOSH Alibaba Cloud CPI.

**Notice: this is a early version, only for test**

## Usage

### Prepare your `Alibaba Cloud` environment

- Create a vpc with switch and get `vswtich_id`
- Create security group get `security_group_id`
- Create user access key, get `access_key_id/access_key_secret`
- Create a jumpbox vm

### Make release

- Clone this repo
- Install golang and export $GOROOT
- Install bosh-cli
- Download `go1.9.linux-amd64.tar.gz` from https://storage.googleapis.com/golang/go1.9.linux-amd64.tar.gz
- Add blob and create release
- Configure golang env and make
- Create bosh release

```
$ git clone https://github.com/aliyun/bosh-alicloud-cpi-release.git
$ cd bosh-alicloud-cpi-release
$ mkdir blobs
$ bosh add-blob ~/Downloads/go1.9.linux-amd64.tar.gz go1.9.linux-amd64.tar.gz
$ source .envrc
$ make
$ bosh create-release --force --tarball=../bosh-alicloud-cpi.tgz
```
*Binary download is not provided now, so make it by your self*

### Install bosh in Alibaba Cloud

- Clone [bosh-deployment](https://github.com/aliyun/bosh-deployment) repo from github
- Checkout alicloud branch

```
$ git clone https://github.com/aliyun/bosh-deployment.git
$ cd bosh-deployment
$ git checkout alicloud
```

use this command, modify the parameters

- vswitch_id
- security_group_id
- access_key_id
- access_key_secret
- region
- zone

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
 -v zone=cn-beijing-a
```

### Install with external-ip(not recommended)

- Create a `Elastic IP` in `Alibaba Cloud Console`, get an `external_ip`

```
bosh create-env bosh.yml \
 --state=state.json \
 --vars-store=creds.yml \
 -o alicloud/cpi.yml \
 -o external-ip-not-recommended.yml \
 -v director_name=my-bosh \
 -v internal_cidr=192.168.0.0/24 \
 -v internal_gw=192.168.0.1 \
 -v internal_ip=192.168.0.6 \
 -v vswitch_id=... \
 -v security_group_id=... \
 -v access_key_id=... \
 -v access_key_secret=... \
 -v region=cn-beijing \
 -v zone=cn-beijing-e \
 -v external_ip=...
```

### *NOTICE*

- *This version of cpi use a temporary public `registry` for install bosh-director, it not safe, we will fix it in later version*
- *Not support mbus_bootstrap_ssl*
- *Not support ssh-tunnel*

### Run Unit Test

Run following commands

```
$ source .envrc
$ make testdeps
$ ginkgo -r -skipPackage integration src/bosh-alicloud-cpi
```

### Run Integration Test

Prepare your `Alibaba Cloud` environment, and export follow variables

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
```

Go to source code path, run follow commands

```
$ make testdeps
$ ginkgo -r src/bosh-alicloud-cpi/integration
```
