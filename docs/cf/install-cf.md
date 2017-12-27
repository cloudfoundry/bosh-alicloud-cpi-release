# Install cf with `cf-deployment`

## Prepare you Alibaba Cloud Environment

- Select a region get `region`
- Create a vpc
- Select a zone get `zone`
- Create a vswitch in your `zone` and get `vswtich_id`, `internal_cidr`
- Create several security groups get their `security_group_id`
- Create user access key, get `access_key_id/access_key_secret`
- Create a key pair, get `key_pair_name` and download it private key, like bosh.pem
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
- region
- zone
- key_pair_name
- private_key

export your BOSH_ENVIRONMENT ip address

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
 -v access_key_id=... \
 -v access_key_secret=... \
 -v region=... \
 -v zone=... \
 -v key_pair_name=... \
 -v private_key=bosh.pem
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

## Install Cloud Foundry

Get `cf-deployment`

```
$ git clone https://github.com/cloudfoundry/cf-deployment.git
```

Base your previous settings, modify `cf-deployment/iaas-support/alicloud/cloud-config.yml`, and update cloud-config

```
bosh -e my-bosh update-cloud-config cf-deployment/iaas-support/alicloud/cloud-config.yml
```

Upload stemcell

```
bosh -e my-bosh upload-stemcell http://bosh.oss-cn-hangzhou.aliyuncs.com/light-bosh-stemcell-1016-alicloud-kvm-ubuntu-trusty-go_agent.tgz
```

Upload CF release, For more CF version refer to [cf-release](https://bosh.io/releases/github.com/cloudfoundry/cf-release?all=1)

```
bosh upload-release https://bosh.io/d/github.com/cloudfoundry/cf-release?v=278 --sha1 7e05e98a9333b187807501ab4252e52058859a2c
```

modify `stemcells` section in `cf-deployment.yml`

```yaml
...
stemcells:
- alias: default
  name: bosh-alicloud-kvm-ubuntu-trusty-go_agent
  version: 1016
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

Install [CF Cli](https://docs.cloudfoundry.org/cf-cli/install-go-cli.html) before using Cloud Foundry.

Login Cloud Foundry

```
cf login -a http://api.$CF_DOMAIN --skip-ssl-validation -u admin -p `bosh int ./cf-vars.yml --path /cf_admin_password`
```

Enjoy your Cloud Foundry