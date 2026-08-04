### Done
+ check slb attach to ecs success
  - setup iaas
  - deploy bosh
  - create release(define Jobs, compile packages, run scripts)
  - upload cloud config
  - deploy release
+ validate the customer upgrade journey for NVMe instance types
  - start on a legacy instance with an encrypted stemcell and cloud-efficiency disks
  - preserve persistent data while upgrading to a 9th-generation instance and ESSD
  - preserve data while shrinking through the Director copy/migrate fallback
  - recreate on NVMe, then move compute back to the legacy family with the ESSD attached

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
