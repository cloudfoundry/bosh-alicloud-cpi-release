/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package action

import (
	"bosh-alicloud-cpi/alicloud"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
)

type DeleteVMMethod struct {
	CallContext
	instances alicloud.InstanceManager
}

func NewDeleteVMMethod(cc CallContext, instances alicloud.InstanceManager) DeleteVMMethod {
	return DeleteVMMethod{cc, instances}
}

func (a DeleteVMMethod) DeleteVM(cid apiv1.VMCID) error {
	instCid := cid.AsString()
	status, err := a.instances.GetInstanceStatus(instCid)

	if err != nil {
		return bosherr.WrapErrorf(err, "delete %s get status failed", instCid)
	}

	if status == alicloud.Deleted {
		return nil
	}

	eniIds := a.instances.GetAttachedNetworkInterfaceIds(instCid)

	err = a.instances.ChangeInstanceStatus(instCid, alicloud.Deleted, func(status alicloud.InstanceStatus) (bool, error) {
		switch status {
		case alicloud.Running, alicloud.Stopped:
			return false, a.instances.DeleteInstance(instCid)
		case alicloud.Deleted:
			return true, a.instances.CleanupInstanceNetworkInterfaces(instCid, eniIds)
		default:
			return false, nil
		}
	})

	if err != nil {
		return bosherr.WrapErrorf(err, "delete %s failed", instCid)
	}

	return nil
}
