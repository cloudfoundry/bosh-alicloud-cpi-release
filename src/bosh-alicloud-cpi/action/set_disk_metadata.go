/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package action

import (
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	"bosh-alicloud-cpi/alicloud"
)

type SetDiskMetadataMethod struct {
	CallContext
	disks alicloud.DiskManager
}

func NewSetDiskMetadataMethod(cc CallContext, disks alicloud.DiskManager) SetDiskMetadataMethod {
	return SetDiskMetadataMethod{cc, disks}
}


func (a SetDiskMetadataMethod) SetDiskMetadata(diskCID apiv1.DiskCID, meta apiv1.DiskMeta) error {
	return nil
}

