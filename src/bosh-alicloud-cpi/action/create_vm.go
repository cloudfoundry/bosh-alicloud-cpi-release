package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/denverdino/aliyungo/common"
	"bosh-alicloud-cpi/alicloud"
	"encoding/json"
	"strings"
	"strconv"
	"bosh-alicloud-cpi/registry"
)

type CreateVMMethod struct {
	runner alicloud.Runner
}

type InstanceProps struct {
	ImageId string 				`json:"image_id,omitempty"`
	AvailabilityZone string		`json:"availability_zone,omitempty"`
	EphemeralDisk DiskInfo 		`json:"ephemeral_disk,omitempty"`
	InstanceName string 		`json:"instance_name,omitempty"`
	InstanceChargeType string	`json:"instance_charge_type,omitempty"`
	InstanceType string 		`json:"instance_type,omitempty"`
	SystemDisk DiskInfo			`json:"system_disk,omitempty"`
	HaltMark string	 			`json:"halt_mark,omitempty"`
}

type DiskInfo struct {
	Size interface{}		`json:"size,omitempty"`
	Type string 			`json:"type,omitempty"`
}

type NetworkProps struct {
	SecurityGroupId string		`json:"security_group_id,omitempty"`
	VSwitchId string			`json:"vswitch_id,omitempty"`
	InternetChargeType string	`json:"internet_charge_type,omitempty"`
}

func (a DiskInfo) GetSize() int {
	switch a.Size.(type) {
	case int:
		return a.Size.(int)
	case string:
		s, _ := a.Size.(string)
		n, err := strconv.Atoi(s)
		if err != nil {
			return 55
		}
		return n
	default:
		return 50
	}
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
	logger.Info("INPUT", "json %s", cloudProps)
	cloudProps.As(&instProps)
	logger.Info("INPUT", "unmarshall CloudProps<Instance>: %s", instProps)

	network := networks["private"]
	networkName := "private"
	if network == nil {
		network = networks["default"]
		networkName = "default"
	}
	_ = networkName

	var networkProps NetworkProps
	network.CloudProps().As(&networkProps)

	logger.Info("INPUT", "unmarshall NetworkProps<Instance>: %s", networkProps)

	var args ecs.CreateInstanceArgs
	args.RegionId = common.Region(a.runner.Config.OpenApi.RegionId)
	args.ZoneId = a.runner.Config.OpenApi.ZoneId
	args.ImageId = stemcellCID.AsString()
	args.UserData = a.runner.Config.Registry.ToInstanceUserData()

	args.InstanceType = instProps.InstanceType
	if strings.Compare(args.InstanceType, "") == 0 {
		args.InstanceType = "ecs.n4.xlarge"
	}

	args.InstanceName = instProps.InstanceName
	args.IoOptimized = "optimized"

	args.SecurityGroupId = networkProps.SecurityGroupId

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

	args.VSwitchId = networkProps.VSwitchId
	args.PrivateIpAddress = network.IP()
	args.InternetMaxBandwidthIn = 5
	args.InternetMaxBandwidthOut = 5
	args.InternetChargeType = common.InternetChargeType(networkProps.InternetChargeType)
	args.AutoRenew = false
	args.Password = "Cloud12345"	// TODO

	req, _ := json.Marshal(args)

	logger.Info("OPENAPI", "Args %s", string(req))

	instid, err := client.CreateInstance(&args)

	if strings.Compare("true", instProps.HaltMark) == 0 {
		return apiv1.VMCID{}, bosherr.Errorf("Halt for test instProps=%s\n networkProps=%s\n args=%s\n",
			instProps, networkProps, args)
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
		Networks: map[string]registry.NetworkSettings {
			"private": {Type: "vip"},
			"public": {Type: "vip"},		// ? why
			//networkName: {
			//	Type: network.Type(),
			//	IP: network.IP(),
			//	Netmask: network.Netmask(),
			//	Gateway: network.Gateway(),
			//	DNS: network.DNS(),
			//},
		},
		Ntp: a.runner.Config.Agent.Ntp,
		VM: registry.VMSettings {
			Name: instid,
		},
	}


	err = a.UpdateAgentSettings(instid, agentSettings)
	if err != nil {
		return apiv1.NewVMCID(instid), bosherr.WrapErrorf(err, "UpdateAgentSettings Failed %s", )
	}

	err = a.runner.StartInstance(instid)

	if err != nil {
		return apiv1.NewVMCID(instid), bosherr.WrapErrorf(err, "StartInstance failed instanceid =", instid)
	}

	err = a.runner.WaitForInstanceStatus(instid, ecs.Running)

	if err != nil {
		return apiv1.NewVMCID(instid), bosherr.WrapErrorf(err, "StartInstance failed instanceid=", instid)
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

