package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/denverdino/aliyungo/common"
	"bosh-alicloud-cpi/alicloud"
)

type CreateVMMethod struct {
	runner alicloud.Runner
}

type DiskProps struct {
	EphemeralDisk DiskInfo `json:"ephemeral_disk"`
	ImageId       string   `json:"image_id"`
	InstanceName  string   `json:"instance_name"`
	InstanceType  string   `json:"instance_type"`
	SystemDisk    DiskInfo `json:"system_disk"`
}

type DiskInfo struct {
	Size int    `json:"size"`
	Type string `json:"cloud_efficiency"`
}

type NetworkProps struct {
	SecurityGroupId string `json:"SecurityGroupId"`
	VSwitchId       string `json:"VSwitchId"`
}

func NewCreateVMMethod(runner alicloud.Runner) CreateVMMethod {
	return CreateVMMethod{runner}
}

func (a CreateVMMethod) CreateVM(
	agentID apiv1.AgentID, stemcellCID apiv1.StemcellCID,
	cloudProps apiv1.VMCloudProps, networks apiv1.Networks,
	associatedDiskCIDs []apiv1.DiskCID, env apiv1.VMEnv) (apiv1.VMCID, error) {

	logger := a.runner.Logger
	client := a.runner.NewClient()

	//
	// convert CloudProps to alicloud dedicated Props
	var diskProps DiskProps
	cloudProps.As(&diskProps)

	network := networks["private"]
	var networkProps NetworkProps
	network.CloudProps().As(&networkProps)

	//
	// TODO stemcellCID verification

	var args ecs.CreateInstanceArgs
	args.RegionId = common.Region(a.runner.Config.OpenApi.RegionId)
	args.ImageId = stemcellCID.AsString()
	args.UserData = a.runner.Config.Actions.Registry.ToInstanceUserData()
	args.ImageId = stemcellCID.AsString()

	args.SecurityGroupId = networkProps.SecurityGroupId
	args.InstanceType = diskProps.InstanceType

	var disk ecs.DataDiskType
	// if diskProps.EphemeralDisk != nil { TODO judge disk type
	disk.Size = diskProps.EphemeralDisk.Size

	args.SystemDisk.Size = diskProps.SystemDisk.Size
	args.SystemDisk.Category = diskProps.SystemDisk.GetCategory()

	args.VSwitchId = networkProps.VSwitchId
	args.PrivateIpAddress = network.IP()

	args.InstanceChargeType = "PostPaid" // TODO
	args.SpotStrategy = "NoSpot"         // TODO
	args.AutoRenew = false

	instid, err := client.CreateInstance(&args)

	if err != nil {
		return apiv1.VMCID{}, bosherr.WrapErrorf(err, "CreateInstance failed %s", args)
	}

	logger.Info("INFO", "FINISHE createvm %s", args)

	// TODO registry settings
	var settings alicloud.RegistrySettings
	settings.Vm.Name = instid
	settings.AgentId = agentID.AsString()
	// TODO
	//settings.Networks=
	if diskProps.SystemDisk.GetCategory() == "cloud" {
		settings.Disks.System = "/dev/xvda"
		settings.Disks.Ephemeral = "/dev/xvdb"
	} else {
		settings.Disks.System = "/dev/vda"
		settings.Disks.Ephemeral = "/dev/vdb"
	}
	// TODO
	//settings.Mbus
	//settings.Blobstore


	logger.Info("INFO", "FINISH update registry settings for VM: %s", instid)

	return apiv1.NewVMCID(instid), nil
}

func (*DiskInfo) GetCategory() ecs.DiskCategory {
	//
	// TODO
	return ecs.DiskCategory("cloud_efficiency")
}
