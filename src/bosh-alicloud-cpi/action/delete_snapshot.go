/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package action

import (
	"bosh-alicloud-cpi/alicloud"

	"github.com/cloudfoundry/bosh-cpi-go/apiv1"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
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
