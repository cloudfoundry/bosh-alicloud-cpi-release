package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	"alibaba/bosh-alicloud-cpi/src/bosh-alicloud-cpi/alicloud"
)

type DeleteVMMethod struct {
	config alicloud.AlicloudConfig
}

func NewDeleteVMMethod(config alicloud.AlicloudConfig) DeleteVMMethod {
	return DeleteVMMethod{config}
}

func (a DeleteVMMethod) DeleteVM(cid apiv1.VMCID) error {
	//vm, _, err := a.vmFinder.Find(cid)
	//if err != nil {
	//	return bosherr.WrapErrorf(err, "Finding vm '%s'", cid)
	//}
	//
	//err = vm.Delete()
	//if err != nil {
	//	return bosherr.WrapErrorf(err, "Deleting vm '%s'", cid)
	//}

	return bosherr.Error("NOT IMPLEMENTED")
}
