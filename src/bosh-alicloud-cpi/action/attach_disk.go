/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package action

import (
	"bosh-alicloud-cpi/alicloud"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	"github.com/denverdino/aliyungo/ecs"
	"bosh-alicloud-cpi/registry"
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

	err := a.disks.AttachDisk(instCid, diskCid)

	if err != nil {
		return a.WrapErrorf(err, "Attaching disk '%s' to VM '%s'", diskCid, instCid)
	}

	device, err := a.disks.WaitForDiskStatus(diskCid, ecs.DiskStatusInUse)

	// client.DescribeDisks()
	registryClient := a.registry
	agentSettings, _ := registryClient.Fetch(instCid)

	agentSettings.AttachPersistentDisk(diskCid, "", device)
	err = registryClient.Update(instCid, agentSettings)
	if err != nil {
		return a.WrapErrorf(err, "UpdateRegistry failed %s %s", diskCid, instCid)
	}


	if err != nil {
		return a.WrapErrorf(err, "WaitForDiskStatus failed %s", diskCid)
	}
	return nil
}

