package action

import (
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	"alibaba/bosh-alicloud-cpi/src/bosh-alicloud-cpi/alicloud"
)

type GetDisksMethod struct {
	config alicloud.AlicloudConfig
}

func NewGetDisksMethod(config alicloud.AlicloudConfig) GetDisksMethod {
	return GetDisksMethod{config}
}

func (a GetDisksMethod) GetDisks(cid apiv1.VMCID) ([]apiv1.DiskCID, error) {
	// todo implement
	return nil, nil
}
