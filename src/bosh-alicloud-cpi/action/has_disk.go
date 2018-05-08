/*
 * Copyright (C) 2017-2018 Alibaba Group Holding Limited
 */
package action

import (
	"github.com/cppforlife/bosh-cpi-go/apiv1"

	"bosh-alicloud-cpi/alicloud"
)

type HasDiskMethod struct {
	CallContext
	disks alicloud.DiskManager
}

func NewHasDiskMethod(cc CallContext, disks alicloud.DiskManager) HasDiskMethod {
	return HasDiskMethod{cc, disks}
}

func (a HasDiskMethod) HasDisk(diskCID apiv1.DiskCID) (bool, error) {

	diskCid := diskCID.AsString()
	disk, err := a.disks.GetDisk(diskCid)

	if err != nil {
		return false, a.WrapErrorf(err, "DescribeDisks failed cid=%s", diskCid)
	}

	if disk == nil {
		return false, nil
	} else {
		return true, nil
	}
}
