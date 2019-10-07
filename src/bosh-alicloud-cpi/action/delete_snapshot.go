/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package action

import (
	"bosh-alicloud-cpi/alicloud"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
)

type DeleteSnapshotMethod struct {
	CallContext
	disks alicloud.DiskManager
}

func NewDeleteSnapshotMethod(cc CallContext, disks alicloud.DiskManager) DeleteSnapshotMethod {
	return DeleteSnapshotMethod{cc, disks}
}

func (a DeleteSnapshotMethod) DeleteSnapshot(snapshotCID apiv1.SnapshotCID) error {
	cid := snapshotCID.AsString()

	err := a.disks.DeleteSnapshot(cid)

	if err != nil {
		return bosherr.WrapErrorf(err, "delete snapshot %s failed", cid)
	}

	return nil
}
