package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	"alibaba/bosh-alicloud-cpi/src/bosh-alicloud-cpi/alicloud"
)

type HasVMMethod struct {
	config alicloud.AlicloudConfig
}

func NewHasVMMethod(config alicloud.AlicloudConfig) HasVMMethod {
	return HasVMMethod{config}
}

func (a HasVMMethod) HasVM(cid apiv1.VMCID) (bool, error) {
	//_, found, err := a.vmFinder.Find(cid)
	//if err != nil {
	//	return false, bosherr.WrapErrorf(err, "Finding VM '%s'", cid)
	//}

	return false, bosherr.Error("NOT IMPLEMENTED")
}
