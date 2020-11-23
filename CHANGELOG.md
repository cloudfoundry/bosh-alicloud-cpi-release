# Change Log

All releases of the BOSH CPI for Alibaba Cloud will be documented in this file.
## 33.0.0 (Unreleased)
## 32.0.0 (August 04, 2020)

IMPROVEMENTS:

- add tags for disk ([#116](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/116))
- fix port and weight conversion issues ([#115](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/115))
- Slbservergroup cancel weight and port defaults ([#114](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/114))

## 31.0.0 (June 29, 2020)

IMPROVEMENTS

- add a given vm to a custom VServerGroup ([#112](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/112))
- support setting tags for bosh vms ([#111](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/111))

BUG FIXES

- Update ClientToken to fix IdempotentFailed (issue: #105) ([#106](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/104))

## 30.0.0 (December 24, 2019)

BUG FIXES

- Update ClientToken to fix IdempotentFailed (issue: #105) ([#82](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/82))

## 29.0.0 (October 08, 2019)

IMPROVEMENTS

- deprecate self-defined wraperror and use bosherr's instead ([#104](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/104))
- remove the updating clienttoken when creating vm failed ([#103](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/103))

## 28.0.0 (September 12, 2019)

IMPROVEMENTS

- fix InvalidIPAddress.AlreadyUsed issue by updating clienttoken ([#101](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/101))

## 27.0.0 (September 09, 2019)

IMPROVEMENTS

- enlarge client timeout and add more error message ([#100](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/100))
- improve pipeline and download light stemcell from bosh.io ([#99](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/99))

## 26.0.0 (September 01, 2019)

IMPROVEMENTS

- update changelog automatically when bump a new release ([#98](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/98))
- force delete instance and cleanup the network interface when deleting vm ([#97](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/97)) 
- update alibaba cloud go sdk to 1.60.132 and add the missing mock methods ([#96](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/96))
- cleanup the left network interface after deleting the instance ([#95](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/95))

## 25.0.0 (August 29, 2019)

IMPROVEMENTS

- enlarge retry times for InvalidIPAddress.AlreadyUsed error ([#93](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/93))
- use light stemcell to test bats and e2e; change downloading cli ways ([#92](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/92))

## 24.0.0 (July 25, 2019)

IMPROVEMENTS

- change to official master and final release version ([#90](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/90))
- Add retry strategy for delete vm when ChangeInstanceStatus ([#89](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/89))
- Fix CI failed bug and improve the pipeline ([#88](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/88))
- Remove the useless input in the ci/task/pub-environment ([#87](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/87))
- Improve CI terraform and to support remote state ([#86](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/86))
- support darwin golang; allow any golang version ([#81](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/81))

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

