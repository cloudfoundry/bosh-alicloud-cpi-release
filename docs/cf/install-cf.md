# Install cf with `cf-deployment`

## Prepare you Alibaba Cloud Environment

- Select a region get `region`
- Create a vpc
- Select a zone get `zone`
- Create a vswitch in your `zone` and get `vswtich_id`, `internal_cidr`
- Create security group get `security_group_id`
- Create user access key, get `access_key_id/access_key_secret`
- Create a jumpbox vm or elastic eip

## Install Bosh 

Install [Bosh Cli](http://bosh.io/docs/cli-v2.html#install) first

Get `bosh-deployment` from github

```
$ git clone https://github.com/aliyun/bosh-deployment.git
$ cd bosh-deployment
$ git checkout alicloud
```

```
export BOSH_ENVIRONMENT=...
```

Use this command, modify the parameters

- vswitch_id
- security_group_id
- access_key_id
- access_key_secret
- access_endpoint
- region
- zone

export your BOSH_ENVIRONMENT ip address

Select `access_endpoint` near your region, you can use the following list:

- Tokyo(ap-northeast-1): ap-northeast-1.aliyuncs.com
- Sydney(ap-southeast-2): ap-southeast-2.aliyuncs.com
- Dubai(me-east-1): me-east-1.aliyuncs.com
- Frankfurt(eu-central-1): eu-central-1.aliyuncs.com

You can use `aliyuncs.com` for most region, check https://github.com/denverdino/aliyungo/blob/master/common/endpoints.xml for more information.

```
bosh create-env bosh-deployment/bosh.yml --state=state.json \
 --vars-store=creds.yml \
 -o bosh-deployment/alicloud/cpi.yml \
 -o bosh-deployment/jumpbox-user.yml \
 -o bosh-deployment/misc/powerdns.yml \
 -v dns_recursor_ip=8.8.8.8 \
 -v director_name=my-bosh \
 -v internal_cidr=192.168.0.0/24 \
 -v internal_gw=192.168.0.1 \
 -v internal_ip=$BOSH_ENVIRONMENT \
 -v vswitch_id=... \
 -v security_group_id=... \
 -v acccess_endpoint=eu-central-1.aliyuncs.com \
 -v access_key_id=... \
 -v access_key_secret=... \
 -v region=... \
 -v zone=...
```

## Login to Bosh

```
bosh int ./creds.yml --path /director_ssl/ca > ca-cert
bosh alias-env my-bosh -e $BOSH_ENVIRONMENT --ca-cert ca-cert
export BOSH_CLIENT=admin
export BOSH_CLIENT_SECRET=`bosh int ./creds.yml --path /admin_password`
export BOSH_CA_CERT=`bosh int ./creds.yml --path /director_ssl/ca`
bosh -e my-bosh login
```

You can use jumpbox.key to ssh into bosh-director

```
bosh int creds.yml --path /jumpbox_ssh/private_key > jumpbox.key
chmod 600 jumpbox.key
ssh jumpbox@$BOSH_ENVIRONMENT -i jumpbox.key
```

## Prepare Cloud Foundry Environment

- Select 3 availability zones
- Create vswitch in each zone get `vswitch_id`, `zone_id`, `internal_cidr`, `internal_gateway`
- Config VPC SNAT with each vswitch to enable vm internet access
- Create a Http SLB and get `http_slb_id`
- Create a domain name wild bind to slb ip. Example: config *.hello-cf.cc to 47.47.47.47
    - You can use 47.47.47.47.xip.io instead custom DNS, but it's not very stable.
- create a Tcp slb get `tcp_slb_id` [optional]


Base your previous settings, modify `bosh-deployment/alicloud/cloud-config.yml` in `bosh-deployment/alicloud/cloud-config.yml`, and update-cloud-config

```
bosh -e my-bosh update-cloud-config bosh-deployment/alicloud/cloud-config.yml
```

## Install Cloud Foundry

Get `cf-deployment`

```
$ git clone https://github.com/cloudfoundry/cf-deployment.git
```

Upload stemcell, you can download stemcell from [here](http://bosh-alicloud.oss-cn-hangzhou.aliyuncs.com/light-bosh-stemcell-1010-alicloud-kvm-ubuntu-trusty-go_agent.tgz)

```
bosh -e my-bosh upload-stemcell light-bosh-stemcell-1009-alicloud-kvm-ubuntu-trusty-go_agent.tgz
```

modify `stemcells` section in `cf-deployment.yml`

```yaml
...
stemcells:
- alias: default
  name: bosh-alicloud-kvm-ubuntu-trusty-go_agent
  version: 1010
```

Setup Domain, use your domain name

```
export CF_DOMAIN=...
```

Install Cloud Foundry...

```
bosh -e my-bosh -d cf deploy cf-deployment.yml \
  --vars-store cf-vars.yml \
  -v system_domain=$CF_DOMAIN
```

Login Cloud Foundry

```
cf login -a http://api.$CF_DOMAIN --skip-ssl-validation -u admin -p `bosh int ./cf-vars.yml --path /cf_admin_password`
```

Enjoy your Cloud Foundry