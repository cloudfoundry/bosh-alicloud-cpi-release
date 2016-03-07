# BOSH Aliyun Cloud Provider Interface

This is the BOSH cpi for Aliyun.

## CPI Methods

- create_stemcell -- Create a Aliyun OS template using stemcell image.
- delete_stemcel  -- Delete a stemcell and the accompanying snapshots

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
