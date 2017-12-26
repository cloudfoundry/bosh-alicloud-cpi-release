# Install cf with `cf-deployment`

## Prepare you Alibaba Cloud Environment

- Select a region get `region`
- Create a vpc
- Select a zone get `zone`
- Create a vswitch in your `zone` and get `vswtich_id`
- Create several security groups get their `security_group_id`
- Create user access key, get `access_key_id/access_key_secret`
- Create a key pair, get `key_pair_name` and download it private key, like bosh.pem
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
- key_pair_name
- private_key


```
export BOSH_ENVIRONMENT=...
```

need to specify `bosh-deployment/alicloud/releases-in-china.yml` in deploying bosh in China

```
bosh create-env bosh-deployment/bosh.yml --state=state.json \
 --vars-store=creds.yml \
 -o bosh-deployment/alicloud/cpi.yml \
 -o bosh-deployment/alicloud/releases-in-china.yml \
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
bosh alias-env my-bosh -e BOSH_ENVIRONMENT --ca-cert ca-cert
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

Upload cf releases (important in China regions, skip it in oversea region)
Install maybe very slow in China regions, so you can upload your releases first.
Get releases from [cf-release-278.zip](http://bosh-alicloud.oss-cn-hangzhou.aliyuncs.com/cf-release-278.zip) and unzip it.


```
bosh upload-release binary-buildpack-release-1.0.14.tgz
bosh upload-release capi-release-1.44.0.tgz
bosh upload-release cf-mysql-release-36.7.0.tgz
bosh upload-release cf-networking-release-1.8.1.tgz
bosh upload-release cf-routing-release-0.166.0.tgz
bosh upload-release cf-smoke-tests-release-39.tgz
bosh upload-release cf-syslog-drain-release-3.tgz
bosh upload-release cflinuxfs2-release-1.165.0.tgz
bosh upload-release consul-release-190.tgz
bosh upload-release diego-release-1.29.1.tgz
bosh upload-release dotnet-core-buildpack-release-1.0.27.tgz
bosh upload-release garden-runc-release-1.9.5.tgz
bosh upload-release go-buildpack-release-1.8.11.tgz
bosh upload-release java-buildpack-release-4.6.tgz
bosh upload-release loggregator-release-99.tgz
bosh upload-release nats-release-22.tgz
bosh upload-release nodejs-buildpack-release-1.6.8.tgz
bosh upload-release php-buildpack-release-4.3.42.tgz
bosh upload-release python-buildpack-release-1.5.26.tgz
bosh upload-release ruby-buildpack-release-1.7.3.tgz
bosh upload-release staticfile-buildpack-release-1.4.16.tgz
bosh upload-release statsd-injector-release-1.0.30.tgz
bosh upload-release uaa-release-52.2.tgz
```

Modify `releases` section in `cf-deployment.yml` (important in China regions, skip it in oversea region)

```yaml
...
releases:
- name: binary-buildpack
  version: latest
- name: capi
  version: latest
- name: cf-mysql
  version: latest
- name: cf-networking
  version: latest
- name: cf-smoke-tests
  version: latest
- name: cf-syslog-drain
  version: latest
- name: cflinuxfs2
  version: latest
- name: consul
  version: latest
- name: diego
  version: latest
- name: dotnet-core-buildpack
  version: latest
- name: garden-runc
  version: latest
- name: go-buildpack
  version: latest
- name: java-buildpack
  version: latest
- name: loggregator
  version: latest
- name: nats
  version: latest
- name: nodejs-buildpack
  version: latest
- name: php-buildpack
  version: latest
- name: python-buildpack
  version: latest
- name: routing
  version: latest
- name: ruby-buildpack
  version: latest
- name: staticfile-buildpack
  version: latest
- name: statsd-injector
  version: latest
- name: uaa
  version: latest
```

Setup Domain, use your domain name

```
export CF_DOMAIN=...
```

Install Cloud Foundry

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

Update buildpacks (important in China regions, skip it in oversea region)
Default cf buildpacks need to download during `cf push`, but is unusable in China region, so download [offline-buildpacks](http://bosh-alicloud.oss-cn-hangzhou.aliyuncs.com/cf-offline-buildpacks_20171107) and update.

```
cf update-buildpack staticfile_buildpack -p staticfile_buildpack-cached-v1.4.18.zip -i 1
cf update-buildpack java_buildpack -p ~/Downloads/java-buildpack-offline-dad0000.zip -i 2
cf update-buildpack ruby_buildpack -p ruby_buildpack-cached-v1.7.5.zip -i 3
cf update-buildpack go_buildpack -p go_buildpack-cached-v1.8.13.zip -i 6
```

Enjoy your Cloud Foundry