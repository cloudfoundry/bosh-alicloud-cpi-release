/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package action

import (
	"bosh-alicloud-cpi/alicloud"
	"fmt"

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
		return a.WrapErrorf(err, "delete %s get status failed", instCid)
	}

	if status == alicloud.Deleted {
		return nil
	}

	err = a.instances.ChangeInstanceStatus(instCid, alicloud.Stopped, func(status alicloud.InstanceStatus) (bool, error) {
		switch status {
		case alicloud.Stopped:
			return true, nil
		case alicloud.Running:
			return false, a.instances.StopInstance(instCid)
		case alicloud.Stopping:
			return false, nil
		default:
			return false, fmt.Errorf("unexpect %s for StopInstance", status)
		}
	})

	if err != nil {
		return bosherr.WrapErrorf(err, "delete %s failed when stop it", instCid)
	}

	err = a.instances.ChangeInstanceStatus(instCid, alicloud.Deleted, func(status alicloud.InstanceStatus) (bool, error) {
		switch status {
		case alicloud.Stopped:
			return false, a.instances.DeleteInstance(instCid)
		case alicloud.Deleted:
			return true, nil
		default:
			return false, fmt.Errorf("unexpect %s status %s for delete", instCid, status)
		}
	})

	if err != nil {
		return a.WrapErrorf(err, "delete %s failed", instCid)
	}

	return nil
}
