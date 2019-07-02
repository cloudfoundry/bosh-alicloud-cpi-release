# Change Log

All releases of the BOSH CPI for Alibaba Cloud will be documented in this file.

## 24.0.0 (Unreleased)

## 23.0.0 (July 02, 2019)

BUG FIXES

- Add retry to fix OperationConflict error when JoinSecurityGroup ([#82](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/82))

## [r22]

IMPROVEMENTS

- Remove the needless code and files in the CI ([#70](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/70))
- Improve bosh cpi CI including assets, e2e, candidate, integration and so on ([#68](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/68))
- Improve integration testcase and mock method ([#67](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/67))
- Upgrade SDK and support HTTP PROXY ([#66](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/66))

BUG FIXES

- fix IdempotentProcessing error when create disk ([#69](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/69))

## [r21]

IMPROVEMENTS

- Remove useless codes when waiting for attaching disk is ok ([#62](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/62))
- Enlarge the timeout of waitting for image ready ([#61](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/61))
- Update go version to 1.11.5 ([#60](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/60))
- Modify oss bucket acl to private when uploading stemcell ([#59](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/59))
- Modify copyright to 2017-2019 ([#55](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/55))
- Support self-define service endpoint and upgrade sdk alibaba-cloud-sdk-go ([#43](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/43))
- Update sdk alibaba-cloud-sdk-go/sdk to improve fetching endpoint ([#41](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/41))

BUG FIXES

- fix device name (xvdc -> vdc) for disk category cloud_ssd ([#58](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/58))

## [r20]

IMPROVEMENTS

- Modify copyright to 2017-2019 ([#55](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/55))
- Support self-define service endpoint and upgrade sdk alibaba-cloud-sdk-go ([#43](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/43))
- Update sdk alibaba-cloud-sdk-go/sdk to improve fetching endpoint ([#41](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/41))

BUG FIXES

- Fix uploading stemcell image faild error by using multipart ([#54](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/54))
- Fix deleting stemcell faild error ([#53](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/53))
- Fix creating vm IdempotentProcessing and add sync lock for client ([#52](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/52))
- Fix full stemcell create flow ([#50](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/50))
- Fix bosh CPI to use https based communication instead of http ([#48](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/48))
- Fix DeleteWithInstance not working bug ([#46](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/46))
- Fix upload stemcell bug and upgrade oss sdk version ([#42](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/42))

## [r19]

- Support security_group_ids in vm_type and vm_extensions. ([#40](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/40))
- New instance field region to support crossing region deployment. ([#39](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/39))
- Deprecate global parameter region_id and use region instead. ([#39](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/39))
- Deprecate global parameter zone_id and use availability_zone instead. ([#39](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/39))
- Add global parameter 'encrypted' to support disk encryption. ([#39](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/39))
- Add global parameter 'security_token' to support STS token. ([#39](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/39))

## [r18]

- Support instance new field charge_period_unit. ([#33](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/33))
- Change go sdk to alibaba cloud official go sdk. ([#33](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/33))
- Fix [issue Ephemeral Disk Encryption](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/issues/33).  ([#33](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/33))
- Format go file. ([#33](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/33))

## [r17]

A temporary release for fixing a go sdk invalid timestamp error. (Changed in local)

## [r16]

- Second GA version
- Fix registry issue
- Support join instance into multiple security groups

## [r15]

- First GA version
- Use instance.UserData instead registry to install bosh 

## [r13]
- Remove `access_endpoint` for regions and using latest SDK to get endpoint automatically.

## [r12]
- Support RAM.

## [r11]
- Support 'SpotStrategy' and 'SpotPriceLimit'.

## [r10]
- Support full cpiv1 methods, add following methods:
    - snapshot_disk
    - delete_snapshot
    - set_disk_metadata
    - resize_disk
- Add ClientToken to keep idempotent when retry, for methods:
    - create_vm
    - create_disk
    - create_snapshot,
- More cloud properties for create_vm methods.
- Configurable `access_endpoint` for regions.

## [r9]

- Support bind instance to SLB
- Stronger InstanceManager with retry and status change
- Improve unit test and integration
- Tolerate small disk size

## [r8]

- For compatible downgrade golang to 1.8.1
- Refactoring code to support unit test and integration test
- Add Concourse ci
- Add docs

## [r7]

- Run with bosh.263.3 and formal manifest v2, but with some issues and bugs.

## [r2 - r6]

- Run with bosh.255 with clipped manifest.v1, can install cf-release-215

## [r1]

- Rewrite with golang, by 4 Alibaba engineers with great excitement, energy and coffee.

## [r0]

- Early ruby version.

