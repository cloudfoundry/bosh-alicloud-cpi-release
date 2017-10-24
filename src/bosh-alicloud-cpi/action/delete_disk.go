/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package action

import (
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	"bosh-alicloud-cpi/alicloud"
	"github.com/denverdino/aliyungo/ecs"
)

type DeleteDiskMethod struct {
	CallContext
	disks alicloud.DiskManager
}

func NewDeleteDiskMethod(cc CallContext, disks alicloud.DiskManager) DeleteDiskMethod {
	return DeleteDiskMethod{cc, disks}
}

func (a DeleteDiskMethod) DeleteDisk(diskCID apiv1.DiskCID) error {
	diskCid := diskCID.AsString()

	disk, err := a.disks.GetDisk(diskCid)
	if err != nil {
		return a.WrapErrorf(err, "delete_disk %s get status failed", diskCid)
	}

	if disk == nil {
		return nil
	}

	if disk.Status != ecs.DiskStatusAvailable {
		return a.WrapErrorf(nil,"delete_disk %s unexpected status %v", diskCid, disk.Status)
	}

	err = a.disks.DeleteDisk(diskCid)
	if err != nil {
		return a.WrapErrorf(err, "delete_disk %s failed", diskCid)
	}

	return nil
}
