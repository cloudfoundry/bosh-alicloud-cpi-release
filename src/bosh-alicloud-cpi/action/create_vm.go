package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	"bosh-alicloud-cpi/alicloud"
)

type CreateVMMethod struct {
	runner alicloud.Runner
}

func NewCreateVMMethod(runner alicloud.Runner) CreateVMMethod {
	return CreateVMMethod{runner}
}

func (a CreateVMMethod) CreateVM(
	agentID apiv1.AgentID, stemcellCID apiv1.StemcellCID,
	cloudProps apiv1.VMCloudProps, networks apiv1.Networks,
	associatedDiskCIDs []apiv1.DiskCID, env apiv1.VMEnv) (apiv1.VMCID, error) {

	//client := a.runner.NewClient()
	//
	//
	//var args ecs.CreateInstanceArgs
	//args.RegionId = common.Region(a.runner.Config.OpenApi.RegionId)
	//args.ImageId = stemcellCID.AsString()
	//args.
	//
	//
	//
	//
	//client.CreateInstance()

	//client :=
	//client.CreateIn

	 // TODO
	return apiv1.VMCID{}, bosherr.Error("NOT IMPLEMENTED")
}
