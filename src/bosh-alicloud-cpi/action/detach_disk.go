/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package action

import (
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	"bosh-alicloud-cpi/alicloud"
	"bosh-alicloud-cpi/registry"
	"github.com/denverdino/aliyungo/ecs"
)

type DetachDiskMethod struct {
	CallContext
	disks alicloud.DiskManager
	registry registry.Client
}

func NewDetachDiskMethod(cc CallContext, disks alicloud.DiskManager, rc registry.Client) DetachDiskMethod {
	return DetachDiskMethod{cc, disks, rc}
}

func (a DetachDiskMethod) DetachDisk(vmCID apiv1.VMCID, diskCID apiv1.DiskCID) error {
	instCid := vmCID.AsString()
	diskCid := diskCID.AsString()

	disk, err := a.disks.GetDisk(diskCid)
	if err != nil {
		return a.WrapErrorf(err, "DetachDisk get disk failed %s", diskCid)
	}

	if disk == nil {
		return a.WrapErrorf(err, "DetachDisk disk missing %s", diskCid)
	}

	if disk.Status != ecs.DiskStatusInUse {
		return a.WrapErrorf(err, "DetachDisk disk %s status expected `InUse` get %s", diskCid, disk.Status)
	}

	err = a.disks.DetachDisk(instCid, diskCid)
	if err != nil {
		return a.WrapErrorf(err, "DetachDisk '%s' from VM '%s' failed", diskCid, instCid)
	}

	_, err = a.disks.WaitForDiskStatus(diskCid, ecs.DiskStatusAvailable)
	if err != nil {
		return a.WrapErrorf(err, "DetachDisk '%s' to VM '%s' wait failed", diskCid, instCid)
	}

	registryClient := a.registry
	agentSettings, _ := registryClient.Fetch(instCid)
	agentSettings.DetachPersistentDisk(diskCid)
	err = registryClient.Update(vmCID.AsString(), agentSettings)
	if err != nil {
		return a.WrapErrorf(err, "DetachDisk update registry failed %s", diskCid)
	}

	return err
}
