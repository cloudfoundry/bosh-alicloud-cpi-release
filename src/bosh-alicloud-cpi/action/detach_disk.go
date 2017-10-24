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

	err := a.disks.DetachDisk(instCid, diskCid)
	if err != nil {
		return a.WrapErrorf(err, "Detach disk '%s' to VM '%s' failed", diskCid, instCid)
	}

	_, err = a.disks.WaitForDiskStatus(diskCid, ecs.DiskStatusAvailable)
	if err != nil {
		return a.WrapErrorf(err, "Detach disk '%s' to VM '%s' wait failed", diskCid, instCid)
	}

	registryClient := a.registry
	agentSettings, _ := registryClient.Fetch(instCid)
	agentSettings.DetachPersistentDisk(diskCid)
	err = registryClient.Update(vmCID.AsString(), agentSettings)
	if err != nil {
		return a.WrapErrorf(err, "UpdateRegistry failed %s", diskCid)
	}

	return err
}
