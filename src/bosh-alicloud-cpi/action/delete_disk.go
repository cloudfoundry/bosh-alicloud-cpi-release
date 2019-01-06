/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package action

import (
	"bosh-alicloud-cpi/alicloud"

	"github.com/cppforlife/bosh-cpi-go/apiv1"
)

type DeleteDiskMethod struct {
	CallContext
	disks     alicloud.DiskManager
	instances alicloud.InstanceManager
}

func NewDeleteDiskMethod(cc CallContext, disks alicloud.DiskManager, instances alicloud.InstanceManager) DeleteDiskMethod {
	return DeleteDiskMethod{cc, disks, instances}
}

func (a DeleteDiskMethod) DeleteDisk(diskCID apiv1.DiskCID) error {
	diskCid := diskCID.AsString()

	disk, err := a.disks.GetDisk(diskCid)
	if err != nil {
		return a.WrapErrorf(err, "DeleteDisk get disk %s failed", diskCid)
	}

	if disk == nil {
		a.Logger.Warn("delete_disk", "disk not found %s", diskCid)
		return nil
	}

	switch alicloud.DiskStatus(disk.Status) {
	case alicloud.DiskStatusAvailable:
		// nothing to do, delete
	case alicloud.DiskStatusInUse:
		instCid := disk.InstanceId
		if instCid != "" {
			inst, err := a.instances.GetInstance(instCid)
			if err != nil {
				return a.WrapErrorf(err, "DeleteDisk %s get host vm %s failed", diskCid, instCid)
			}
			if inst != nil && alicloud.InstanceStatus(inst.Status) != alicloud.Deleted {
				return a.WrapErrorf(err, "can't delete disk %s with a <%s> vm %s", diskCid, inst.Status, instCid)
			}
		}
		_, err = a.disks.WaitForDiskStatus(diskCid, alicloud.DiskStatusAvailable)
		if err != nil {
			return a.WrapErrorf(err, "DeleteDisk %s waiting status from %s to %s failed",
				diskCid, disk.Status, alicloud.DiskStatusAvailable)
		}
	case alicloud.DiskStatusDetaching:
		_, err = a.disks.WaitForDiskStatus(diskCid, alicloud.DiskStatusAvailable)
		if err != nil {
			return a.WrapErrorf(err, "DeleteDisk %s waiting status from %s to %s failed",
				diskCid, disk.Status, alicloud.DiskStatusAvailable)
		}
	default:
		return a.Errorf("DeleteDisk %s unexpected status %s", diskCid, disk.Status)
	}

	err = a.disks.DeleteDisk(diskCid)
	if err != nil {
		return a.WrapErrorf(err, "delete_disk %s failed", diskCid)
	}

	return nil
}
