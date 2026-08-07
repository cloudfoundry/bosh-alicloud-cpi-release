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

 The release and version objects are a separate matter. Concourse's built-in
 `s3` and `semver` resource types accept only a static credential and cannot
 refresh an assumed one, so four resources still need an access key:
 `bosh-cpi-dev-artifacts`, `version-semver` and `release-version-semver` through
 `bosh_cpi_release_*`, and `aliyun-cli` through `aliyun_cli_*`. Reaching zero
 requires a role-aware OSS resource type.

 Four RAM roles are needed, and the trust relationship differs between them.

 Attached to an instance, so they trust `ecs.aliyuncs.com`:

 * the worker role on the Concourse worker instance. It needs only
   `sts:AssumeRole`: every task bootstraps from it and assumes one of the roles
   below.
 * the role named by `director_role_name`, holding the CPI's ECS/VPC/SLB/NLB/OSS
   permissions. Terraform attaches it to the director VM.

 Assumed by the worker role, so they trust it rather than a service. Name the
 worker role, `acs:ram::<account-id>:role/<worker-role>`, so only sessions of
 that role can assume them. `acs:ram::<account-id>:root` also works but trusts
 every identity in the account, so prefer the narrower form:

 * `terraform_role_arn`, which creates and destroys the test environment
 * `e2e_observer_role_arn`, read-only. The E2E errands only describe resources,
   so they must not be given the provisioning role.

 A role that trusts only `ecs.aliyuncs.com` cannot be assumed, so giving these
 two the same trust policy as the first two makes every task fail at AssumeRole.

 `terraform_role_arn` needs `ram:PassRole` to attach the director role to the
 director VM, and its `Resource` must name that role:

 ```json
 {
   "Effect": "Allow",
   "Action": "ram:PassRole",
   "Resource": "acs:ram::<account-id>:role/BoshDirectorRole"
 }
 ```

 An unrestricted `ram:PassRole` would let the pipeline attach any role in the
 account, including a more privileged one, to an instance it creates and then
 read that role's credentials from the instance's metadata. Withholding
 `ram:CreatePolicy` and `ram:AttachPolicyToRole` does not prevent that on its
 own.

 ```

 alicloud_region__primary:     REGION # cn-beijing
 terraform_role_arn:           acs:ram::<account-id>:role/<provisioning-role>
 e2e_observer_role_arn:        acs:ram::<account-id>:role/<read-only-role>
 director_role_name:           BoshDirectorRole
 terraform_backend_bucket:     OSS_BUCKET_FOR_TERRAFORM_STATE
 terraform_backend_region:     REGION
 bosh_cpi_release_bucket:      OSS_BUCKET_FOR_RELEASE_ARTIFACTS
 bosh_cpi_release_region:      REGION
 bosh_cpi_release_access_key:  ACCESS_KEY_ID
 bosh_cpi_release_secret_key:  ACCESS_KEY_SECRET
 aliyun_cli_bucket:            OSS_BUCKET_HOLDING_THE_ALIYUN_CLI
 aliyun_cli_region:            REGION # the endpoint is derived from this
 aliyun_cli_access_key:        ACCESS_KEY_ID
 aliyun_cli_secret_key:        ACCESS_KEY_SECRET
 PUBLIC_KEY:                   "ssh-ed25519 ... # must match the private key below"
 github_user_email:            YOUR_GITHUB_ACCOUNT_EMAIL
 github_user_name:             YOUR_GITHUB_ACCOUNT_NAME
 github_user_id:               YOUR_GITHUB_ACCOUNT_ID
 github_user_password:         YOUR_GITHUB_ACCOUNT_PASSWORD
 github_bosh-alicloud-cpi-release_private-key: |
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