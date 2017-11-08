/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package action

import (
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	"bosh-alicloud-cpi/alicloud"
)

type DeleteSnapshotMethod struct {
	CallContext
	disks alicloud.DiskManager
}

func NewDeleteSnapshotMethod(cc CallContext, disks alicloud.DiskManager) DeleteSnapshotMethod {
	return DeleteSnapshotMethod{cc, disks}
}

func (a DeleteSnapshotMethod) DeleteSnapshot(snapshotCID apiv1.SnapshotCID) (error) {
	return nil
}