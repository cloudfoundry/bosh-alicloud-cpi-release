# Change Log

All releases of the BOSH CPI for Alibaba Cloud will be documented in this file.

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

