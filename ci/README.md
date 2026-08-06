## BOSH Alicloud CPI Concourse Pipeline
  In order to run the BOSH Alicloud CPI Concourse Pipeline you must have an existing [Concourse](http://concourse.ci/) environment. See [Deploying Concourse on Alicloud Compute Engine](https://yq.aliyun.com/articles/178450?spm=5176.8091938.0.0.5IFmqS) for instructions.

* Target your Concourse CI environment:

```
fly -t alicloud login -c <YOUR CONCOURSE URL>
```

And then:
 * Create the variable file vars-pipeline-develop.yml

 The tasks that drive the test environment take no access key. Each one starts
 from the RAM role attached to the Concourse worker ECS instance and assumes a
 role for the work it needs to do, so only role ARNs, buckets and regions are
 configured for them.

 The release and version objects are a separate matter: Concourse's built-in
 `s3` and `semver` resource types accept only a static credential and cannot
 refresh an assumed one, so `aliyun_cli_*` still holds an access key. Reaching
 zero requires a role-aware OSS resource type.

 Two RAM roles are needed, both trusting `ecs.aliyuncs.com`:

 * one attached to the Concourse worker, holding only `sts:AssumeRole`
 * the role named by `director_role_name`, holding the CPI's ECS/VPC/SLB/NLB/OSS
   permissions; terraform attaches it to the director VM

 `e2e_observer_role_arn` is read-only: the E2E errands only describe resources,
 so they must not be given the provisioning role.

 ```

 alicloud_region__primary:     REGION # cn-beijing
 terraform_role_arn:           acs:ram::<account-id>:role/<provisioning-role>
 e2e_observer_role_arn:        acs:ram::<account-id>:role/<read-only-role>
 director_role_name:           BoshDirectorRole
 terraform_backend_bucket:     OSS_BUCKET_FOR_TERRAFORM_STATE
 terraform_backend_region:     REGION
 aliyun_cli_bucket:            OSS_BUCKET_HOLDING_THE_ALIYUN_CLI
 aliyun_cli_endpoint:          OSS_ENDPOINT
 aliyun_cli_access_key:        ACCESS_KEY_ID
 aliyun_cli_secret_key:        ACCESS_KEY_SECRET
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