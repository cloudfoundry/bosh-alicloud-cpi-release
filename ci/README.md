## BOSH Alicloud CPI Concourse Pipeline
  In order to run the BOSH Alicloud CPI Concourse Pipeline you must have an existing [Concourse](http://concourse.ci/) environment. See [Deploying Concourse on Alicloud Compute Engine](https://yq.aliyun.com/articles/178450?spm=5176.8091938.0.0.5IFmqS) for instructions.

* Target your Concourse CI environment:

```
fly -t alicloud login -c <YOUR CONCOURSE URL>
```

And then:
 * Create the variable file vars-pipeline-develop.yml

 ```

 alicloud_access_key__primary: YOUR_ALICLOUD_ACCESS_KEY
 alicloud_secret_key__primary: YOUR_ALICLOUD_SECRET_KEY
 alicloud_region__primary:     REGION # cn-beijing
 alicloud__cpi_stemcell_id:    CPI_STEMCELL_ID # m-2ze1cneefoj075diqyeh
 alicloud__cpi_internal_ip:    CPI_INTERNAL_IP # 172.16.0.2
 alicloud__cpi_internal_gw:    CPI_INTERNAL_GATEWAY # 172.16.0.1
 github_user_email:            YOUR_GITHUB_ACCOUNT_EMAIL
 github_user_name:             YOUR_GITHUB_ACCOUNT_NAME
 github_user_id:               YOUR_GITHUB_ACCOUNT_ID
 github_user_password:         YOUR_GITHUB_ACCOUNT_PASSWORD
 gitlab_bosh-alicloud-cpi-release_private-key: |
   -----BEGIN RSA PRIVATE KEY-----
   YOUR_LOCAL_PRIVATE_KEY
   -----END RSA PRIVATE KEY-----
 alicloud_director_vars_file: |
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