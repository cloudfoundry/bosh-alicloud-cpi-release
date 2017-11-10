/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package action

import (
"github.com/cppforlife/bosh-cpi-go/apiv1"
"bosh-alicloud-cpi/alicloud"
)

type SnapshotDiskMethod struct {
	CallContext
	disks alicloud.DiskManager
}

func NewSnapshotDiskMethod(cc CallContext, disks alicloud.DiskManager) SnapshotDiskMethod {
	return SnapshotDiskMethod{cc, disks}
}


func (a SnapshotDiskMethod) SnapshotDisk(diskCID apiv1.DiskCID, meta apiv1.DiskMeta) (apiv1.SnapshotCID, error) {
	diskCid := diskCID.AsString()

	//
	// TODO, get snapshot name from meta
	ssid, err := a.disks.CreateSnapshot(diskCid, "")
	if err != nil {
		return apiv1.SnapshotCID{}, a.WrapErrorf(err, "create snapshot for %s failed", diskCid)
	}

	return apiv1.NewSnapshotCID(ssid), nil
}