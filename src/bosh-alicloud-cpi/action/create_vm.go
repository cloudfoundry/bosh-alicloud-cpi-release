package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	"alibaba/bosh-alicloud-cpi/src/bosh-alicloud-cpi/alicloud"
)

type CreateVMMethod struct {
	config alicloud.AlicloudConfig
}

func NewCreateVMMethod(config alicloud.AlicloudConfig) CreateVMMethod {
	return CreateVMMethod{config: config}
}

func (a CreateVMMethod) CreateVM(
	agentID apiv1.AgentID, stemcellCID apiv1.StemcellCID,
	cloudProps apiv1.VMCloudProps, networks apiv1.Networks,
	associatedDiskCIDs []apiv1.DiskCID, env apiv1.VMEnv) (apiv1.VMCID, error) {

	client:=a.config.NewClient();
	client.CreateIn

	 // TODO
	return apiv1.VMCID{}, bosherr.Error("NOT IMPLEMENTED")
}
