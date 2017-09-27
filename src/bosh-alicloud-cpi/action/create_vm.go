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

type InstanceProps struct {
	ImageId string 				`json:"image_id"`
	EphemeralDisk DiskInfo 		`json:"ephemeral_disk"`
	InstanceName string 		`json:"instance_name"`
	InstanceChargeType string	`json:"instance_charge_type"`
	InstanceType string 		`json:"instance_type"`
	SystemDisk DiskInfo			`json:"system_disk"`
	AvailabilityZone string		`json:"availability_zone"`
}

type DiskInfo struct {
	Size int				`json:"size"`
	Type string 			`json:"type"`
}

type NetworkProps struct {
	SecurityGroupId string		`json:"security_group_id"`
	VSwitchId string			`json:"vswitch_id"`
	InternetChargeType string	`json:"internet_charge_type"`
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
	var instProps InstanceProps
	cloudProps.As(&instProps)

	network := networks["private"]
	if network == nil {
		network = networks["default"]
	}

	var networkProps NetworkProps
	network.CloudProps().As(&networkProps)

	var args ecs.CreateInstanceArgs
	args.RegionId = common.Region(a.runner.Config.OpenApi.RegionId)
	args.ZoneId = a.runner.Config.OpenApi.ZoneId
	args.ImageId = stemcellCID.AsString()
	args.UserData = a.runner.Config.Registry.ToInstanceUserData()

	args.InstanceType = instProps.InstanceType
	args.InstanceName = instProps.InstanceName
	args.IoOptimized = "optimized"

	args.SecurityGroupId = networkProps.SecurityGroupId

	disk := instProps.EphemeralDisk
	if disk.Type != "" {
		args.DataDisk = []ecs.DataDiskType{
			{Size: disk.Size, Category: ecs.DiskCategory(disk.Type),},
		}
	}

	disk = instProps.SystemDisk
	if disk.Type != "" {
		args.SystemDisk.Size = disk.Size
		args.SystemDisk.Category = ecs.DiskCategory(disk.Type)
	} else {
		args.SystemDisk.Size = 50
		args.SystemDisk.Category = ecs.DiskCategory("cloud_efficiency")
	}

	args.VSwitchId = networkProps.VSwitchId
	args.PrivateIpAddress = network.IP()
	args.InternetMaxBandwidthIn = 5
	args.InternetMaxBandwidthOut = 5
	args.InternetChargeType = common.InternetChargeType(networkProps.InternetChargeType)

	args.InstanceChargeType = common.InstanceChargeType(instProps.InstanceChargeType)
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

