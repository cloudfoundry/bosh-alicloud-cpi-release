/*
 * Copyright (C) 2017-2018 Alibaba Group Holding Limited
 */
package action

import (
	"bosh-alicloud-cpi/alicloud"

	"github.com/cppforlife/bosh-cpi-go/apiv1"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
)

type CreateDiskMethod struct {
	CallContext
	disks     alicloud.DiskManager
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
		return cid, a.WrapErrorf(err, "create_disk GetInstance failed %s", vmCid.AsString())
	}

	if inst == nil {
		return cid, a.Errorf("create_disk missing instance id=%s", vmCid.AsString())
	}

	disk, err := NewDiskInfoWithSize(size, props)

	if err != nil {
		return cid, a.WrapErrorf(err, "create_disk check input failed %n, %v", size, props)
	}

	args := ecs.CreateCreateDiskRequest()
	args.ZoneId = inst.ZoneId
	args.Size = requests.NewInteger(disk.GetSizeGB())
	args.DiskCategory = string(disk.GetCategory())
	encrypt := disk.Encrypted
	if encrypt == nil {
		encrypt = a.Config.OpenApi.Encrypted
	}
	args.Encrypted = requests.NewBoolean(*encrypt)
	diskCid, err := a.disks.CreateDisk(inst.RegionId, args)

	if err != nil {
		return cid, a.WrapError(err, "create_disk failed")
	}

	cid = apiv1.NewDiskCID(diskCid)
	return cid, nil
}
