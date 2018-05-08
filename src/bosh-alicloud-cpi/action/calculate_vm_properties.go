package action

import (
	"github.com/cppforlife/bosh-cpi-go/apiv1"
)

type CalculateVMCloudPropertiesMethod struct {
	CallContext
}

func NewCalculateVMCloudPropertiesMethod(cc CallContext) CalculateVMCloudPropertiesMethod {
	return CalculateVMCloudPropertiesMethod{cc}
}

func (a CalculateVMCloudPropertiesMethod) CalculateVMCloudProperties(resources apiv1.VMResources) (apiv1.VMCloudProps, error) {
	return apiv1.NewVMCloudPropsFromMap(nil), nil

}
