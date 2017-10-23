/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package action

import (
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	"bosh-alicloud-cpi/alicloud"
)

type CreateDiskMethod struct {
	CallContext
	disks alicloud.DiskManager
}

func NewCreateDiskMethod(cc CallContext, disks alicloud.DiskManager) CreateDiskMethod {
	return CreateDiskMethod{cc, disks}
}

func (a CreateDiskMethod) CreateDisk(size int, props apiv1.DiskCloudProps, vmCid *apiv1.VMCID) (apiv1.DiskCID, error) {
	a.Logger.Debug("create_disk", "size=%d, cloudProps=%v, vmCid=%v", size, props, vmCid)
	var cid apiv1.DiskCID

	diskInfo, err := NewDiskInfo(size, props)

	if err != nil {
		return cid, a.WrapErrorf(err, "check CreateDisk input failed %n, %v", size, props)
	}

	diskCid, err := a.disks.CreateDisk(diskInfo.GetSizeGB(), diskInfo.GetCategory())

	if err != nil {
		return cid, a.WrapError(err, "CreateDisk failed")
	}

	cid = apiv1.NewDiskCID(diskCid)

	//
	// TODO?
	// vm_cid [String]: Cloud ID of the VM created disk will most likely be attached; it could be used to
	// .optimize disk placement so that disk is located near the VM.
	//
	// with zone
	// NewAttachDiskMethod(a.runner).AttachDisk(*vmcid, diskcid)
	//
	//if vmCid != nil {
	//	return cid, a.WrapError(nil, "NOT_IMPLEMENTED create_disk and attach with vm")
	//
	//}

	return cid, nil
}
