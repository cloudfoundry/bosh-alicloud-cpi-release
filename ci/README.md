## BOSH Alicloud CPI Concourse Pipeline
  In order to run the BOSH Alicloud CPI Concourse Pipeline you must have an existing [Concourse](http://concourse.ci/) environment. See [Deploying Concourse on Alicloud Compute Engine](https://yq.aliyun.com/articles/178450?spm=5176.8091938.0.0.5IFmqS) for instructions.

* Target your Concourse CI environment:

```
fly -t alicloud login -c <YOUR CONCOURSE URL>
```

And then:
 * Create the variable file vars-pipeline-develop.yml

 The pipeline holds no long-lived Alibaba Cloud access key. Every task starts from
 the RAM role attached to the Concourse worker ECS instance and assumes
 `terraform_role_arn` for the work it needs to do, so only the role ARN, bucket
 and region are configured here. See `docs/design/zero-long-lived-ak-spec.md`.

 ```

 alicloud_region__primary:     REGION # cn-beijing
 terraform_role_arn:           acs:ram::<account-id>:role/<provisioning-role>
 terraform_backend_bucket:     OSS_BUCKET_FOR_TERRAFORM_STATE
 terraform_backend_region:     REGION
 PUBLIC_KEY:                   "ssh-ed25519 ... # must match the private key below"
 github_user_email:            YOUR_GITHUB_ACCOUNT_EMAIL
 github_user_name:             YOUR_GITHUB_ACCOUNT_NAME
 github_user_id:               YOUR_GITHUB_ACCOUNT_ID
 github_user_password:         YOUR_GITHUB_ACCOUNT_PASSWORD
 gitlab_bosh-alicloud-cpi-release_private-key: |
   -----BEGIN RSA PRIVATE KEY-----
   YOUR_LOCAL_PRIVATE_KEY
   -----END RSA PRIVATE KEY-----
 alicloud_director_vars_file: |
   private_key: |
     -----BEGIN RSA PRIVATE KEY-----
     YOUR_BOSH_DIRECTOR_PRIVATE_KEY
     -----END RSA PRIVATE KEY-----
 ```

* Set the BOSH Alicloud CPI pipeline:

```
fly -t alicloud set-pipeline -p bosh-alicloud-cpi -c pipeline-develop.yml --load-vars-from vars-pipeline-develop.yml
```

* Unpause the BOSH Google CPI pipeline:

```
fly -t alicloud unpause-pipeline -p bosh-alicloud-cpi