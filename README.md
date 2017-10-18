# BOSH Alibaba Cloud CPI

This is a [BOSH](http://bosh.io) release for the BOSH Alibaba Cloud CPI.

**Notice: this is a early version, only for test**

## Usage

## 1. Prepare your `Alibaba Cloud` environment

- Create a vpc with switch and get `vswtich_id`
- Create security group get `security_group_id`
- Create user access key, get `access_key_id/access_key_secret`
- Create a jumpbox vm

## 2. Install bosh in Alibaba Cloud

- clone [bosh-deployment](https://github.com/aliyun/bosh-deployment) repo from github
- checkout alicloud branch

```
git clone https://github.com/aliyun/bosh-deployment.git
cd bosh-deployment
git checkout alicloud
```

use this command, modify following settings

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


## Install with external-ip(not recommended)

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
 -v vswitch_id=vsw-2zemyfytfclbcmgfkzokx \
 -v security_group_id=sg-2zei0mcphxbdxj49qtmz \
 -v access_key_id=... \
 -v access_key_secret=... \
 -v region=cn-beijing \
 -v zone=cn-beijing-e \
 -v external_ip=...
```

## *NOTICE*

- *This version of cpi use a temporary public `registry` for install bosh-director, it not safe, we will fix it in later version*
- *Not support mbus_bootstrap_ssl *
- *Not support ssh-tunnel*


