# Change Log

All releases of the BOSH CPI for Alibaba Cloud will be documented in this file.
## 55.0.0 (Unreleased)
## 54.0.0 (July 01, 2025)

- update go version to a specified version ([#182](https://github.com/cloudfoundry/bosh-alicloud-cpi-release/pull/182))
- go: download go1.24 for linux/amd64: toolchain not available ([#180](https://github.com/cloudfoundry/bosh-alicloud-cpi-release/pull/180))
   
## 53.0.0 (April 24, 2025)

- upgrade go verison to 1.24.2 ([#179](https://github.com/cloudfoundry/bosh-alicloud-cpi-release/pull/179))
- improves the ci test ([#178](https://github.com/cloudfoundry/bosh-alicloud-cpi-release/pull/178))

## 52.0.0 (February 11, 2025)

- throw running error when binding nlb server groups ([#177](https://github.com/cloudfoundry/bosh-alicloud-cpi-release/pull/177))
- os agnostic bash path ([#176](https://github.com/cloudfoundry/bosh-alicloud-cpi-release/pull/176))

## 51.0.0 (September 05, 2024)

- supports binding vm to nlb server groups by setting nlb_server_groups ([#173](https://github.com/cloudfoundry/bosh-alicloud-cpi-release/pull/173))
- upgrades the alibaba cloud sdk verison and go version to 1.21 ([#172](https://github.com/cloudfoundry/bosh-alicloud-cpi-release/pull/172))


## 50.0.0 (June 14, 2024)

- Revert "supports to encrypt system disk" ([#171](https://github.com/cloudfoundry/bosh-alicloud-cpi-release/pull/171))

## 49.0.0 (April 10, 2024)

- supports to encrypt system disk ([#169](https://github.com/cloudfoundry/bosh-alicloud-cpi-release/pull/169))

## 48.0.0 (December 26, 2023)

- Update Documentation ([#166](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/166))
- update dependencies ([#165](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/165))

## 47.0.0 (October 25, 2023)

IMPROVEMENTS:

- adds new instance properties: nlb_server_group_weight, nlb_server_group_port and nlb_server_group_ids to support binding vm with nlb ([#163](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/163))
- upgrade go version to 1.20; upgrade bosh-utils version to v0.0.407 ([#164](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/164))

## 46.0.0 (May 22, 2023)

IMPROVEMENTS:

- Fixes the disk device diff bug after attching it to instance ([#161](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/161))
- fixes the empty tag value bug when creating vm instance ([#160](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/160))
- fixes the parsing terraform output error before prepare-director ([#159](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/159))
- fixes the terraform destroy failed ([#158](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/158))
- ci: setting provider version ([#157](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/157))
- ci: Improves the terraform template ([#156](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/156))
- ci: upgrade terraform version ([#155](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/155))

## 45.0.0 (March 30, 2022)

IMPROVEMENTS:

- Add NTP field to BoshEnv struct ([#152](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/152))

## 44.0.0 (January 19, 2021)

IMPROVEMENTS:

- fix the groupTagValue bug and format the code ([#151](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/151))
- Update for issue #149 ([#150](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/150)) 
- Make agent credentials optional for accessing blobstore with signed URLs enabled ([#148](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/148))

## 43.0.0 (November 01, 2021)

IMPROVEMENTS:

- use env bosh group as cr required tag ([#147](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/147))

## 42.0.0 (October 22, 2021)

IMPROVEMENTS:

- Add CR required bosh tag when create vm ([#146](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/146))

## 41.0.1 (October 21, 2021)

Publish a new tag same as the version 41.0.0

## 41.0.0 (October 18, 2021)

IMPROVEMENTS:

- support to setting tags in which the manifest file ([#145](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/145))
- Supports getting tags from registry ([#143](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/143))

## 40.0.0 (August 01, 2021)

BUG FIXES：

Fixes the reboot instance bug when attrching disk with cpi verison v2 ([#141](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/141))

## 39.0.0 (June 24, 2021)

IMPROVEMENTS:

- Changing the adding tags action to Create ([#135](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/135))
- specify terraform provider source ([#134](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/134))
- Add blobstores field to BoshEnv struct ([#133](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/133))

## 38.0.0 (May 25, 2021)

IMPROVEMENTS:

- add support for cpi api v2 ([#131](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/131))

## 37.0.0 (May 20, 2021)

BUG FIXES：

- Fixes the deleteing full stemcell failed bug when setting the encrypted is true ([#130](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/130))

## 36.0.0 (May 8, 2021)

BUG FIXES：

- Fixes the addTags bug when copy image ([#129](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/129))

## 35.0.0 (April 28, 2021)

IMPROVEMENTS:

- Delete ram image when there has copied image ([#128](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/128))
- Fixes the InvalidImageName.Duplicated issue when copyImage ([#127](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/127))

## 34.0.0 (April 18, 2021)

IMPROVEMENTS:

- Supports to encrypt the system disk by copyImage when creating stemcell ([#125](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/125))
- improves the create_stemcell to avoid opening ipv6 setting when building bionic stemcell ([#124](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/124))

## 33.0.0 (January 10, 2021)

IMPROVEMENTS:

- Modify the data type and restore the default value ([#120](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/120))
- improve ci environment for alicloud_slb ([#118](https://github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/pull/118))

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

