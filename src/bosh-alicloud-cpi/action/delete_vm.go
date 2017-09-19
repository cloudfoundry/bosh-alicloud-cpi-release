package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	"alibaba/bosh-alicloud-cpi/src/bosh-alicloud-cpi/alicloud"
)

type DeleteVMMethod struct {
	runner alicloud.Runner
}

func NewDeleteVMMethod(runner alicloud.Runner) DeleteVMMethod {
	return DeleteVMMethod{runner}
}

func (a DeleteVMMethod) DeleteVM(cid apiv1.VMCID) error {
	client := a.runner.NewClient()

	//
	// TODO try detect VM is exists
	//vm, _, err := a.vmFinder.Find(cid)
	//if err != nil {
	//	return bosherr.WrapErrorf(err, "Finding vm '%s'", cid)
	//}

	a.runner.Logger.Info("CLOUD", "try DeleteVM:", cid.AsString())

	err := client.DeleteInstance(cid.AsString())

	if err != nil {
		return bosherr.WrapErrorf(err, "Deleting vm '%s'", cid)
	}

	return bosherr.Error("NOT IMPLEMENTED")
}
