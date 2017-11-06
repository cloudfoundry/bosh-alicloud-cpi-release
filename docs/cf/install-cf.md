# Install cf with `cf-deployment`

## Prepare you Alibaba Cloud Environment

- Select a region get `region`
- Create a vpc
- Select a zone get `zone`
- Create a vswitch in your `zone` and get `vswtich_id`
- Create security group get `security_group_id`
- Create user access key, get `access_key_id/access_key_secret`
- Create a jumpbox vm

## Install Bosh

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


export your BOSH_DIRECTOR_IP

```
export BOSH_DIRECTOR_IP=...
```

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
 -v internal_ip=$BOSH_DIRECTOR_IP \
 -v vswitch_id=... \
 -v security_group_id=... \
 -v access_key_id=... \
 -v access_key_secret=... \
 -v region=... \
 -v zone=...
```

## Login to Bosh


```
bosh int ./creds.yml --path /director_ssl/ca > ca-cert
bosh alias-env my-bosh -e $BOSH_DIRECTOR_IP --ca-cert ca-cert
export BOSH_ENVIRONMENT=$BOSH_DIRECTOR_IP
export BOSH_CLIENT=admin
export BOSH_CLIENT_SECRET=`bosh int ./creds.yml --path /admin_password`
export BOSH_CA_CERT=`bosh int ./creds.yml --path /director_ssl/ca`
bosh -e my-bosh login
```

You can use jumpbox.key to ssh into bosh-director

```
bosh int creds.yml --path /jumpbox_ssh/private_key > jumpbox.key
chmod 600 jumpbox.key
ssh jumpbox@${BOSH_DIRECTOR_IP} -i jumpbox.key
```

## Prepare Cloud Foundry Environment

- Select 3 availability zones
- create vswitch in each zone get `vswitch_id`, `zone_id`, `internal_cidr`, `internal_gateway`
- create a Http SLB get `http_slb_id`
- create a domain name wild bind to slb ip. Example: config *.hello-cf.cc to 47.47.47.47
    - You can use 47.47.47.47.xip.io instead custom DNS, but it's not very stable.
- create a Tcp slb get `tcp_slb_id` [optional]


Base your previous settings, modify `bosh-deployment/alicloud/cloud-config.yml` in `bosh-deployment/alicloud/cloud-config.yml`, and update-cloud-config

```
bosh -e my-bosh update-cloud-config bosh-deployment/alicloud/cloud-config.yml
```


## Install Cloud Foundry

Get `cf-deployment`

```
$ git clone https://github.com/cloudfoundry/cf-deployment.yml
```

Setup Domain

```
export CF_DOMAIN=hello-cf.cc
```

Install Cloud Foundry

```
bosh -e my-bosh -d cf deploy cf-deployment.yml \
  --vars-store cf-vars.yml \
  -v system_domain=hello-cf.cc
```

Login Cloud Foundry

```
cf login -a http://api.cloudfoundry.cc --skip-ssl-validation -u admin -p `bosh int ./cf-vars.yml --path /cf_admin_password`
```