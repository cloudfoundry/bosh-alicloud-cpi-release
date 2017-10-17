package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/denverdino/aliyungo/common"
	"bosh-alicloud-cpi/alicloud"
	"encoding/json"
	"strings"
	"bosh-alicloud-cpi/registry"
)

type CreateVMMethod struct {
	runner alicloud.Runner
}

type InstanceProps struct {
	ImageId string 				`json:"image_id"`
	AvailabilityZone string		`json:"availability_zone"`
	InstanceName string 		`json:"instance_name"`
	InstanceChargeType string	`json:"instance_charge_type"`
	InstanceType string 		`json:"instance_type"`
	InstanceRole string			`json:"instance_role"`
	EphemeralDisk DiskInfo 		`json:"ephemeral_disk"`
	SystemDisk DiskInfo			`json:"system_disk"`
}

func NewCreateVMMethod(runner alicloud.Runner) CreateVMMethod {
	return CreateVMMethod{runner}
}

func (a CreateVMMethod) CreateVM(
	agentID apiv1.AgentID, stemcellCID apiv1.StemcellCID,
	cloudProps apiv1.VMCloudProps, networkArgs apiv1.Networks,
	associatedDiskCIDs []apiv1.DiskCID, env apiv1.VMEnv) (apiv1.VMCID, error) {

	cid := apiv1.VMCID{}

	logger := a.runner.Logger
	client := a.runner.NewClient()

	logger.Info("ENV", "inv: %v", env)
	env2, err := registry.UnmarshalEnvSettings(env)
	if err != nil {
		return cid, bosherr.WrapErrorf(err,"UnmarshalEnvSettings failed %v", env)
	}
	logger.Info("ENV", "inv: %v", env2)

	//
	// convert CloudProps to alicloud dedicated Props
	var args ecs.CreateInstanceArgs

	var instProps InstanceProps
	logger.Info("INPUT", "json %s", cloudProps)
	err = cloudProps.As(&instProps)
	if err != nil {
		return cid, bosherr.WrapErrorf(err, "unmarshal CloudProps failed %v", cloudProps)
	}

	logger.Info("INPUT", "unmarshal CloudProps<Instance>: %s", instProps)

	logger.Info("INPUT", "unmarshal NetworkProps<Instance>: %v", networkArgs)
	networks, err := NewNetworks(networkArgs)
	if err != nil {
		return cid, bosherr.WrapErrorf(err, "create_vm failed when parse Networks %v", networkArgs)
	}

	networks.FillCreateInstanceArgs(&args)

	args.RegionId = common.Region(a.runner.Config.OpenApi.RegionId)
	args.ZoneId = instProps.AvailabilityZone
	args.ImageId = stemcellCID.AsString()
	args.UserData = a.runner.Config.Registry.ToInstanceUserData()

	args.InstanceType = instProps.InstanceType
	args.InstanceChargeType = common.PostPaid
	if strings.Compare(args.InstanceType, "") == 0 {
		args.InstanceType = "ecs.n4.xlarge"
	}

	args.InstanceName = instProps.InstanceName
	args.IoOptimized = "optimized"

	disks, err := NewDisks(instProps.SystemDisk, []DiskInfo{instProps.EphemeralDisk})
	if err != nil {
		return cid, bosherr.WrapErrorf(err, "bad disks format, %v", instProps)
	}
	disks.FillCreateInstanceArgs(&args)

	//
	//args.SecurityGroupId = networks.GetSecurityGroupId()
	//args.VSwitchId = networks.GetVSwitchId()
	//args.PrivateIpAddress = networks.GetPrivateAddress()
	//args.InternetMaxBandwidthIn = networks.GetInternetMaxBandwidthIn()
	//args.InternetMaxBandwidthOut = networks.GetInternetMaxBandwidthOut()
	//args.InternetChargeType = networks.GetInternetChargeType()
	//args.AutoRenew = false
	args.Password = "Cloud12345"

	req, _ := json.Marshal(args)

	logger.Info("OPENAPI", "Args %s", string(req))

	//
	// insert agent re
	agentSettings := registry.AgentSettings {
		AgentID: agentID.AsString(),
		Blobstore: a.runner.Config.Agent.Blobstore.AsRegistrySettings(),
		Disks: disks.AsRegistrySettings(),
		Env: env2,
		Mbus: a.runner.Config.Agent.Mbus,
		Networks: networks.AsRegistrySettings(),
		Ntp: a.runner.Config.Agent.Ntp,
		VM: registry.VMSettings {
			Name: "",
		},
	}

	if strings.Compare("fake", instProps.InstanceRole) == 0 {
		j1, _ := json.Marshal(args)
		j2, _ := json.Marshal(agentSettings)
		return apiv1.VMCID{}, bosherr.Errorf("haltForTest instProps=%v\n networkProps=%v\n args=%s\n registry=%s\n",
			instProps, networks, j1, j2)
	}

	//
	// do create instance
	instid, err := client.CreateInstance(&args)
	if err != nil {
		return apiv1.VMCID{}, bosherr.WrapErrorf(err, "CreateInstance failed INPUT=%s ", string(req))
	}

	agentSettings.VM.Name = instid

	//
	//
	for _, diskCid := range associatedDiskCIDs {
		err := a.runner.AttachDisk(instid, diskCid.AsString())
		if err != nil {
			return cid, bosherr.WrapErrorf(err, "associate Persistent Disk error diskCid=%s", diskCid)
		}

		path, err := a.runner.WaitForDiskStatus(diskCid.AsString(), ecs.DiskStatusInUse)
		if err != nil {
			return cid, bosherr.WrapErrorf(err, "associate and WaitForDiskStatus Failed diskCid=%s", diskCid)
		}
		disks.AssociatePersistentDisk(diskCid.AsString(), path)
	}

	err = a.UpdateAgentSettings(instid, agentSettings)
	if err != nil {
		return apiv1.NewVMCID(instid), bosherr.WrapErrorf(err, "UpdateAgentSettings Failed %s", )
	}

	err = a.runner.StartInstance(instid)
	if err != nil {
		return apiv1.NewVMCID(instid), bosherr.WrapErrorf(err, "StartInstance failed cid =", instid)
	}

	err = a.runner.WaitForInstanceStatus(instid, ecs.Running)
	if err != nil {
		return apiv1.NewVMCID(instid), bosherr.WrapErrorf(err, "StartInstance failed cid=", instid)
	}

	err = networks.BindInstanceEip(client, instid, args.RegionId)
	if err != nil {
		return apiv1.NewVMCID(instid), bosherr.WrapErrorf(err, "StartInstance failed cid=")
	}

	err = networks.BindInstanceEip(client, instid, args.RegionId)
	if err != nil {
		return apiv1.NewVMCID(instid), bosherr.WrapErrorf(err, "StartInstance failed instanceid=")
	}

	//
	// TODO: every error must free created vm before terminated
	logger.Info("INFO", "FINISH create_vm %s", args)
	return apiv1.NewVMCID(instid), nil
}

func (a CreateVMMethod) UpdateAgentSettings(instId string, agentSettings registry.AgentSettings) error {
	client := a.runner.GetHttpRegistryClient()
	err := client.Update(instId, agentSettings)

	if err != nil {
		json, _ := json.Marshal(agentSettings)
		a.runner.Logger.Error("create_vm", "UpdateAgentSettings to registery failed %s json:%s", json)
		return bosherr.WrapErrorf(err, "UpdateAgentSettings failed %v %s", client, json)
	}

	return nil
}

