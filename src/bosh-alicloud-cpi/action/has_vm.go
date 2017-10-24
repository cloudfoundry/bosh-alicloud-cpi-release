/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package action

import (
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	"bosh-alicloud-cpi/alicloud"
)

type HasVMMethod struct {
	CallContext
	instances alicloud.InstanceManager
}

func NewHasVMMethod(cc CallContext, instances alicloud.InstanceManager) HasVMMethod {
	return HasVMMethod{cc, instances}
}

func (a HasVMMethod) HasVM(cid apiv1.VMCID) (bool, error) {
	//
	//
	instCid := cid.AsString()
	inst, err := a.instances.GetInstance(instCid)

	if err != nil {
		return false, a.WrapErrorf(err, "Finding VM Failed '%s'", cid)
	}

	if inst != nil {
		return true, nil
	} else {
		return false, nil
	}
}
