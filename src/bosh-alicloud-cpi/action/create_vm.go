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
	EphemeralDisk DiskInfo 		`json:"ephemeral_disk"`
	InstanceName string 		`json:"instance_name"`
	InstanceChargeType string	`json:"instance_charge_type"`
	InstanceType string 		`json:"instance_type"`
	InstanceRole string			`json:"instance_role"`
	SystemDisk DiskInfo			`json:"system_disk"`
}

type DiskInfo struct {
	Size json.Number		`json:"size,omitempty"`
	Type string				`json:"type,omitempty"`
}

func (a DiskInfo) GetSize() int {
	n, err := a.Size.Int64()
	if err != nil {
		n = 55
	}
	return int(n)
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

	//
	// convert CloudProps to alicloud dedicated Props
	var args ecs.CreateInstanceArgs

	var instProps InstanceProps
	logger.Info("INPUT", "json %s", cloudProps)
	cloudProps.As(&instProps)
	logger.Info("INPUT", "unmarshal CloudProps<Instance>: %s", instProps)

	logger.Info("INPUT", "unmarshal NetworkProps<Instance>: %v", networkArgs)
	networks, err := NewNetworks(networkArgs)
	if err != nil {
		return cid, bosherr.WrapErrorf(err, "create_vm failed when parse Networks %v", networkArgs)
	}

	networks.FillCreateInstanceArgs(&args)

	args.RegionId = common.Region(a.runner.Config.OpenApi.RegionId)
	args.ZoneId = a.runner.Config.OpenApi.ZoneId //TODO use AZ
	args.ImageId = stemcellCID.AsString()
	args.UserData = a.runner.Config.Registry.ToInstanceUserData()

	args.InstanceType = instProps.InstanceType
	if strings.Compare(args.InstanceType, "") == 0 {
		args.InstanceType = "ecs.n4.xlarge"
	}

	args.InstanceName = instProps.InstanceName
	args.IoOptimized = "optimized"

	disk := instProps.EphemeralDisk
	if disk.Type != "" {
		args.DataDisk = []ecs.DataDiskType{
			{Size: disk.GetSize(), Category: ecs.DiskCategory(disk.Type),},
		}
	}

	disk = instProps.SystemDisk
	if disk.Type != "" {
		args.SystemDisk.Size = disk.GetSize()
		args.SystemDisk.Category = ecs.DiskCategory(disk.Type)
	} else {
		args.SystemDisk.Size = 50
		args.SystemDisk.Category = ecs.DiskCategory("cloud_efficiency")
	}

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

	instid, err := client.CreateInstance(&args)

	if strings.Compare("fake", instProps.InstanceRole) == 0 {
		return apiv1.VMCID{}, bosherr.Errorf("Halt for test instProps=%s\n networkProps=%s\n args=%s\n",
			instProps, networks, args)
	}

	if err != nil {
		return apiv1.VMCID{}, bosherr.WrapErrorf(err, "CreateInstance failed INPUT=%s ", string(req))
	}

	// insert agent re
	agentSettings := registry.AgentSettings {
		AgentID: agentID.AsString(),
		Blobstore: a.runner.Config.Agent.Blobstore.AsRegistrySettings(),
		Disks: registry.DisksSettings {
			System: "/dev/vda",
			Ephemeral: "/dev/vdb",
			Persistent: map[string]interface{} {},
		},
		Env: registry.EnvSettings {
			BoshEnv: registry.BoshEnv {
				Password: "",
				KeepRootPassword: false,
				RemoveDevTools: false,
			},
			PersistentDiskFs: "",
		},
		Mbus: a.runner.Config.Agent.Mbus,
		Networks: networks.AsRegistrySettings(),
		Ntp: a.runner.Config.Agent.Ntp,
		VM: registry.VMSettings {
			Name: instid,
		},
	}

	if strings.Compare("director", instProps.InstanceRole) != 0 {
		err = a.UpdateAgentSettings(instid, agentSettings)
		if err != nil {
			return apiv1.NewVMCID(instid), bosherr.WrapErrorf(err, "UpdateAgentSettings Failed %s", )
		}
	}

	err = a.runner.StartInstance(instid)

	if err != nil {
		return apiv1.NewVMCID(instid), bosherr.WrapErrorf(err, "StartInstance failed instanceid =", instid)
	}

	err = a.runner.WaitForInstanceStatus(instid, ecs.Running)

	if err != nil {
		return apiv1.NewVMCID(instid), bosherr.WrapErrorf(err, "StartInstance failed instanceid=", instid)
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
		//
		// TODO, when first time to create director, this action must failed, how to aviod it??
		// return bosherr.WrapErrorf(err, "UpdateAgentSettings failed %s %s %s", clientOptions, agentSettings, conf)
	}

	return nil
}

