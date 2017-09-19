package action

import (
	"github.com/cppforlife/bosh-cpi-go/apiv1"
)

type RebootVMMethod struct{}

func NewRebootVMMethod() RebootVMMethod {
	return RebootVMMethod{}
}

func (a RebootVMMethod) RebootVM(_ apiv1.VMCID) error {
	return nil
}
