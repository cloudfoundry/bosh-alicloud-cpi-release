/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package action

import (
	"bosh-alicloud-cpi/alicloud"
	"fmt"

	"github.com/cppforlife/bosh-cpi-go/apiv1"
)

type RebootVMMethod struct {
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

	err = a.instances.ChangeInstanceStatus(instCid, alicloud.Running, func(status alicloud.InstanceStatus) (bool, error) {
		switch status {
		case alicloud.Running:
			return true, nil
		case alicloud.Stopping:
			return false, nil
		case alicloud.Starting:
			return false, nil
		default:
			return false, fmt.Errorf("unexpect status %s", status)
		}
	})

	if err != nil {
		return a.WrapErrorf(err, "reboot %s failed", instCid)
	}

	return nil
}
