### Done
+ check slb attach to ecs success
  - setup iaas
  - deploy bosh
  - create release(define Jobs, compile packages, run scripts)
  - upload cloud config
  - deploy release

### TODO
+ spot instance.
  - CPI will support it.
  - refers: https://help.aliyun.com/document_detail/25499.html?spm=5176.product25365.6.826.4G5N3N
+ create ecs by heavy stemcell.
  - Alicloud does not support, we may do it by work around.
+ ram profile attached to ecs
  - Alicloud does not support.
+ raw ephemeral disk.
  - Alicloud does not support.
  - refers: https://bosh.io/docs/aws-instance-storage.html

### Nice to have
+ install dependency offline