# BOSH Aliyun Cloud Provider Interface

This is the BOSH cpi for Aliyun.

## stemcell

Alicloud only allows creation of images from running VMs on Alicloud, which means CPI must use already imported OS image and that’s where Alicloud light stemcells come in.
Light stemcell tarballs include additional details about already imported OS images in the ```cloud_properties``` section. For example light stemcells for Alicloud have ```image-id``` in the ```cloud_properties``` section, that contains region-to-image-id mappings. When Alicloud CPI’s create_stemcell call is made, it will return matching ```image-id``` without doing any IaaS API calls.

```
name: cf-alicloud-kvm-hubuntu-trusty-go_agent
version: "1001"
operating_system: ubuntu-trusty
sha1: e2a346c8ce067270f397cc08fb010051450d0175
bosh_protocol: "1"
cloud_properties:
  name: bosh-alicloud-kvm-hubuntu-trusty-go_agent
  version: "1001"
  infrastructure: alicloud
  architecture: x86_64
  hypervisor: kvm
  image_id:
    cn-beijing: m-78hyduj
    cn-hangzhou: m-j89jnsd
    cn-shanghai: m-8ju9kng
    
```

### CPI Methods

- create_stemcell -- Return matching image id using light-stemcell.
- delete_stemcel  -- Delete a stemcell.

## Test

```
bundle exec rake
```

## Pull Request

Please follow these steps to make a contribution.

1. Fork the repository

2. Create a feature branch (`git checkout -b feature/your_feature_name`)

   - Run the tests to ensure that your local environment is working `bundle && bundle exec rake` (this may take a while).

3. Make changes on the branch:

   - Add a feature

       a. Add tests for the new feature
       b. Make the tests pass

   - Fixing a bug

       a. Add a test/tests which exercises the bug
       b. Fix the bug, making the tests pass

   - Refactoring existing functionality

       a. Change the implementation
       b. Ensure that tests still pass
           - If you find yourself changing tests after a refactor, consider refactoring the tests first.

4. Push to your fork (`git push origin feature/your_feature_name`) and submit a pull request selecting `develop` as the target branch.
