/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package action

import (
	"bosh-alicloud-cpi/alicloud"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	"github.com/denverdino/aliyungo/ecs"
	"bosh-alicloud-cpi/registry"
	"fmt"
)

type AttachDiskMethod struct {
	CallContext
	disks alicloud.DiskManager
	registry registry.Client
}

func NewAttachDiskMethod(cc CallContext, disks alicloud.DiskManager, rc registry.Client) AttachDiskMethod {
	return AttachDiskMethod{cc, disks, rc}
}


func (a AttachDiskMethod) AttachDisk(vmCID apiv1.VMCID, diskCID apiv1.DiskCID) error {
	instCid := vmCID.AsString()
	diskCid := diskCID.AsString()
	device := ""

	err := a.disks.ChangeDiskStatus(diskCid, ecs.DiskStatusInUse, func(disk *ecs.DiskItemType) (bool, error) {
		if disk == nil {
			return false, fmt.Errorf("missing disk %s", diskCid)
		}
		switch disk.Status {
		case ecs.DiskStatusInUse:
			device = alicloud.AmendDiskPath(disk.Device, disk.Category)
			return true, nil
		case ecs.DiskStatusAvailable:
			return false, a.disks.AttachDisk(instCid, diskCid)
		case ecs.DiskStatusAttaching:
			return false, nil
		default:
			return false, fmt.Errorf("unexcepted disk %s status %s", diskCid, disk.Status)
		}
	})

	if err != nil {
		return a.WrapErrorf(err, "attach disk %s to %s failed", diskCid, instCid)
	}

	registryClient := a.registry
	agentSettings, _ := registryClient.Fetch(instCid)
	agentSettings.AttachPersistentDisk(diskCid, "", device)

	err = registryClient.Update(instCid, agentSettings)
	if err != nil {
		return a.WrapErrorf(err, "update registry failed %s %s", diskCid, instCid)
	}
	return nil
}

