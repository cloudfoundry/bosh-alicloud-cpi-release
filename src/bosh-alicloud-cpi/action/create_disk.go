/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package action

import (
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	"bosh-alicloud-cpi/alicloud"
	"github.com/denverdino/aliyungo/ecs"
)

type CreateDiskMethod struct {
	CallContext
	disks alicloud.DiskManager
	instances alicloud.InstanceManager
}

func NewCreateDiskMethod(cc CallContext, disks alicloud.DiskManager, instances alicloud.InstanceManager) CreateDiskMethod {
	return CreateDiskMethod{cc, disks, instances}
}

func (a CreateDiskMethod) CreateDisk(size int, props apiv1.DiskCloudProps, vmCid *apiv1.VMCID) (apiv1.DiskCID, error) {
	var cid apiv1.DiskCID

	//
	// vm_cid [String]: Cloud ID of the VM created disk will most likely be attached;
	// it could be used to .optimize disk placement so that disk is located near the VM.
	//
	if vmCid == nil {
		return cid, a.Errorf("create_disk must provide vmCid")
	}

	inst, err := a.instances.GetInstance(vmCid.AsString())
	if err != nil {
		return cid, a.WrapErrorf(err,"create_disk GetInstance failed %s", vmCid.AsString())
	}

	if inst == nil {
		return cid, a.Errorf("create_disk missing instance id=%s", vmCid.AsString())
	}

	disk, err := NewDiskInfoWithSize(size, props)

	if err != nil {
		return cid, a.WrapErrorf(err, "create_disk check input failed %n, %v", size, props)
	}

	var args ecs.CreateDiskArgs
	args.ZoneId = inst.ZoneId
	args.Size = disk.GetSizeGB()
	args.DiskCategory = disk.GetCategory()
	args.Encrypted = disk.Encrypted
	diskCid, err := a.disks.CreateDisk(args)

	if err != nil {
		return cid, a.WrapError(err, "create_disk failed")
	}

	cid = apiv1.NewDiskCID(diskCid)
	return cid, nil
}
