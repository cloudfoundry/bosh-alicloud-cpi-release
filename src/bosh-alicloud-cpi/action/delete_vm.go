/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	"bosh-alicloud-cpi/alicloud"
	"github.com/denverdino/aliyungo/ecs"
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
		return bosherr.WrapErrorf(err, "DeleteVM get status failed cid=%s", instCid)
	}

	if status == ecs.Stopped {
		// nothing to do, can delete now
	} else if status == ecs.Running {
		err = a.instances.StopInstance(instCid)

		if err != nil {
			return bosherr.WrapErrorf(err, "DeleteVM when stopping vm failed cid=%s", instCid)
		}

		err = a.instances.WaitForInstanceStatus(instCid, ecs.Stopped)

		if err != nil {
			return bosherr.WrapErrorf(err, "DeleteVM waiting for vm stopped failed cid=%s", instCid)
		}
	} else if status == ecs.Stopping {
		err = a.instances.WaitForInstanceStatus(instCid, ecs.Stopped)

		if err != nil {
			return bosherr.WrapErrorf(err, "DeleteVM waiting for vm stopped failed cid=%s", instCid)
		}
	} else {
		return bosherr.WrapErrorf(err, "DeleteVM unexcepted status: %s cid=%s", status, instCid)
	}

	a.Logger.Info("CLOUD", "try DeleteVM:", cid.AsString())

	err = a.instances.DeleteInstance(cid.AsString())

	if err != nil {
		return bosherr.WrapErrorf(err, "Deleting vm '%s' failed %s", cid, status)
	}

	err = a.instances.WaitForInstanceStatus(cid.AsString(), ecs.Deleted)

	if err != nil {
		return bosherr.WrapErrorf(err, "Deleting vm waitForStatus '%s", cid)
	}

	return nil
}
