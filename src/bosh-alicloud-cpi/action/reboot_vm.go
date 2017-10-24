/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package action

import (
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	"bosh-alicloud-cpi/alicloud"
	"github.com/denverdino/aliyungo/ecs"
)

type RebootVMMethod struct{
	CallContext
	instances alicloud.InstanceManager
}

func NewRebootVMMethod(cc CallContext, instances alicloud.InstanceManager) RebootVMMethod {
	return RebootVMMethod{cc, instances}
}

func (a RebootVMMethod) RebootVM(cid apiv1.VMCID) error {
	instCid := cid.AsString()
	err := a.instances.RebootInstance(instCid)

	if err != nil {
		return a.WrapErrorf(err, "RebootInstance failed cid=%s", instCid)
	}

	err = a.instances.WaitForInstanceStatus(instCid, ecs.Running)

	if err != nil {
		return a.WrapErrorf(err, "WaitForInstanceStatus failed cid=%s", instCid)
		return err
	}

	return nil
}
