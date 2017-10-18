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
	runner alicloud.Runner
}

func NewRebootVMMethod(runner alicloud.Runner) RebootVMMethod {
	return RebootVMMethod{runner}
}

func (a RebootVMMethod) RebootVM(cid apiv1.VMCID) error {
	instid := cid.AsString()
	err := a.runner.RebootInstance(instid)
	if err != nil {
		//TODO logging
		return err
	}

	err = a.runner.WaitForInstanceStatus(instid, ecs.Running)

	if err != nil {
		//TODO logging
		return err
	}

	return nil
}
