package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/denverdino/aliyungo/common"
	"bosh-alicloud-cpi/alicloud"
	"encoding/json"
)

type CreateVMMethod struct {
	runner alicloud.Runner
}

type DiskProps struct {
	EphemeralDisk DiskInfo 	`json:"ephemeral_disk"`
	ImageId string 			`json:"image_id"`
	InstanceName string 	`json:"instance_name"`
	InstanceType string 	`json:"instance_type"`
	SystemDisk DiskInfo		`json:"system_disk"`
}

type DiskInfo struct {
	Size int				`json:"size"`
	Type string 			`json:"cloud_efficiency"`
}

type NetworkProps struct {
	SecurityGroupId string	`json:"security_group_id"`
	VSwitchId string		`json:"vswitch_id"`
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
	if network == nil {
		network = networks["default"]
	}
	var networkProps NetworkProps
	network.CloudProps().As(&networkProps)

	logger.Info("NETWORK", "IP: ", network.IP())
	logger.Info("NETWORK", "NETMASK: ", network.Netmask())

	//
	// TODO stemcellCID verification

	var args ecs.CreateInstanceArgs
	args.RegionId = common.Region(a.runner.Config.OpenApi.RegionId)
	args.ZoneId = a.runner.Config.OpenApi.ZoneId
	args.ImageId = stemcellCID.AsString()
	args.UserData = a.runner.Config.Registry.ToInstanceUserData()

	args.InstanceType = "ecs.mn4.small"
	args.InstanceName = diskProps.InstanceName
	args.IoOptimized = "optimized"

	args.SecurityGroupId = networkProps.SecurityGroupId
	// args.InstanceType = diskProps.InstanceType

	var disk ecs.DataDiskType
	// if diskProps.EphemeralDisk != nil { TODO judge disk type
	disk.Size = diskProps.EphemeralDisk.Size
	args.DataDisk = []ecs.DataDiskType{
		{Size: 50, Category: "cloud_efficiency", },
	}

	args.SystemDisk.Size = 50 // diskProps.SystemDisk.Size
	args.SystemDisk.Category = "cloud_ssd"

	args.VSwitchId = networkProps.VSwitchId
	args.PrivateIpAddress = network.IP()

	args.InstanceChargeType = "PostPaid"		// TODO
	args.AutoRenew = false

	args.Password = "Cloud12345"	// TODO

	req, _ := json.Marshal(args)

	logger.Info("OPENAPI", "Args %s", string(req))


	instid, err := client.CreateInstance(&args)

	if err != nil {
		return apiv1.VMCID{}, bosherr.WrapErrorf(err, "CreateInstance failed INPUT=%s AccessKeyId=%s", string(req), a.runner.Config.OpenApi.AccessKeyId)
	}

	err = a.runner.StartInstance(instid)

	if err != nil {
		return apiv1.NewVMCID(instid), bosherr.WrapErrorf(err, "StartInstance failed instanceid =", err)
	}

	err = a.runner.WaitForInstanceStatus(instid, ecs.Running)

	if err != nil {
		return apiv1.NewVMCID(instid), bosherr.WrapErrorf(err, "StartInstance failed instanceid =", err)
	}

	logger.Info("INFO", "FINISH create_vm %s", args)
	return apiv1.NewVMCID(instid), nil
}


func (* DiskInfo) GetCategory() ecs.DiskCategory {
	//
	// TODO
	return ecs.DiskCategory("cloud_efficiency")
}