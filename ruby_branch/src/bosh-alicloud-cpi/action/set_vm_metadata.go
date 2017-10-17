package action

import (
	"github.com/cppforlife/bosh-cpi-go/apiv1"
)

type SetVMMetadataMethod struct{}

func NewSetVMMetadataMethod() SetVMMetadataMethod {
	return SetVMMetadataMethod{}
}

func (a SetVMMetadataMethod) SetVMMetadata(_ apiv1.VMCID, meta apiv1.VMMeta) error {
	// todo can properties be set on the container
	return nil
}
