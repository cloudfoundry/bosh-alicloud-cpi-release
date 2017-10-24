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
		return a.WrapErrorf(err, "delete_disk/getStatus cid=%s", diskCid)
	}

	if disk == nil {
		return nil
	}

	switch disk.Status {
	case ecs.DiskStatusDetaching:
		_, err := a.disks.WaitForDiskStatus(diskCid, ecs.DiskStatusAvailable)
		if err != nil {
			return a.WrapErrorf(err, "delete_disk/wait for detaching failed=%s", diskCid)
		}
	case ecs.DiskStatusAvailable:
		break
	default:
		return a.WrapErrorf(nil,"delete_disk/unexcepted_DiskStatus: %s-%v", diskCid, disk.Status)
	}

	err = a.disks.DeleteDisk(diskCid)
	if err != nil {
		return a.WrapErrorf(err, "DeleteDisk %s failed", diskCid)
	}

	return nil
}
