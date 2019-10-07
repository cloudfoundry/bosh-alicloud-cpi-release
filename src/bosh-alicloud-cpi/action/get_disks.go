/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package action

import (
	"bosh-alicloud-cpi/alicloud"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
)

type GetDisksMethod struct {
	CallContext
	disks alicloud.DiskManager
}

func NewGetDisksMethod(cc CallContext, disks alicloud.DiskManager) GetDisksMethod {
	return GetDisksMethod{cc, disks}
}

func (a GetDisksMethod) GetDisks(cid apiv1.VMCID) ([]apiv1.DiskCID, error) {

	instCid := cid.AsString()

	disks, err := a.disks.GetDisks(instCid)

	if err != nil {
		return nil, bosherr.WrapErrorf(err, "DescribeDisks failed cid=%s", instCid)
	}

	var results []apiv1.DiskCID
	for _, v := range disks {
		results = append(results, apiv1.NewDiskCID(v.DiskId))
	}

	return results, nil
}
