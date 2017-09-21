package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	"bosh-alicloud-cpi/alicloud"
	"github.com/denverdino/aliyungo/ecs"
)

type DeleteVMMethod struct {
	runner alicloud.Runner
}

func NewDeleteVMMethod(runner alicloud.Runner) DeleteVMMethod {
	return DeleteVMMethod{runner}
}

func (a DeleteVMMethod) DeleteVM(cid apiv1.VMCID) error {
	client := a.runner.NewClient()

	instid := cid.AsString()
	status, err := a.runner.GetInstanceStatus(instid)

	if err != nil {
		return bosherr.WrapErrorf(err, "DeleteVM get status failed cid=%s", instid)
	}

	if status == ecs.Running {
		err = a.runner.StopInstance(instid)

		if err != nil {
			return bosherr.WrapErrorf(err, "DeleteVM when stopping vm failed cid=%s", instid)
		}

		err = a.runner.WaitForInstanceStatus(instid, ecs.Stopped)

		if err != nil {
			return bosherr.WrapErrorf(err, "DeleteVM waiting for vm stopped failed cid=%s", instid)
		}
	} else if status == ecs.Stopping {
		err = a.runner.WaitForInstanceStatus(instid, ecs.Stopped)

		if err != nil {
			return bosherr.WrapErrorf(err, "DeleteVM waiting for vm stopped failed cid=%s", instid)
		}
	} else {
		return bosherr.WrapErrorf(err, "DeleteVM unexcepted status: %s cid=%s", status, instid)
	}

	a.runner.Logger.Info("CLOUD", "try DeleteVM:", cid.AsString())

	err = client.DeleteInstance(cid.AsString())

	if err != nil {
		return bosherr.WrapErrorf(err, "Deleting vm '%s'", cid)
	}

	return nil
}
