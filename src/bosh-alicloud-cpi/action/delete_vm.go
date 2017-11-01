/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	"bosh-alicloud-cpi/alicloud"
	"github.com/denverdino/aliyungo/ecs"
	"fmt"
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
		if status == ecs.Deleted {
			return nil
		}
		return a.WrapErrorf(err, "delete %s get status failed", instCid)
	}

	err = a.instances.ChangeInstanceStatus(instCid, ecs.Stopped, func(status ecs.InstanceStatus) (bool, error) {
		switch status {
		case ecs.Stopped:
			return true, nil
		case ecs.Running:
			return false, a.instances.StopInstance(instCid)
		case ecs.Stopping:
			return false, nil
		default:
			return false, fmt.Errorf("unexpect %s for StopInstance", status)
		}
	})

	if err != nil {
		return bosherr.WrapErrorf(err, "delete %s failed when stop it", instCid)
	}

	err = a.instances.ChangeInstanceStatus(instCid, ecs.Deleted, func(status ecs.InstanceStatus) (bool, error) {
		switch status {
		case ecs.Stopped:
			return false, a.instances.DeleteInstance(instCid)
		case ecs.Deleted:
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
