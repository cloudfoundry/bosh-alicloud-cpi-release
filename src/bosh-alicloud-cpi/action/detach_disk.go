/*
 * Copyright (C) 2017-2018 Alibaba Group Holding Limited
 */
package action

import (
	"bosh-alicloud-cpi/alicloud"
	"bosh-alicloud-cpi/registry"
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
)

type DetachDiskMethod struct {
	CallContext
	disks    alicloud.DiskManager
	registry registry.Client
}

func NewDetachDiskMethod(cc CallContext, disks alicloud.DiskManager, rc registry.Client) DetachDiskMethod {
	return DetachDiskMethod{cc, disks, rc}
}

func (a DetachDiskMethod) DetachDisk(vmCID apiv1.VMCID, diskCID apiv1.DiskCID) error {
	instCid := vmCID.AsString()
	diskCid := diskCID.AsString()

	err := a.disks.ChangeDiskStatus(diskCid, alicloud.DiskStatusAvailable, func(disk *ecs.Disk) (bool, error) {
		if disk == nil {
			return false, fmt.Errorf("missing disk %s", diskCid)
		}
		switch alicloud.DiskStatus(disk.Status) {
		case alicloud.DiskStatusInUse:
			return false, a.disks.DetachDisk(instCid, diskCid)
		case alicloud.DiskStatusAvailable:
			return true, nil
		case alicloud.DiskStatusDetaching:
			return false, nil
		default:
			return false, fmt.Errorf("unexpect disk %s status %s", diskCid, disk.Status)
		}
	})

	if err != nil {
		return a.WrapErrorf(err, "detach disk %s from %s failed", diskCid, instCid)
	}

	registryClient := a.registry
	agentSettings, _ := registryClient.Fetch(instCid)
	agentSettings.DetachPersistentDisk(diskCid)
	err = registryClient.Update(instCid, agentSettings)
	if err != nil {
		return a.WrapErrorf(err, "DetachDisk update registry failed %s", diskCid)
	}

	return err
}
