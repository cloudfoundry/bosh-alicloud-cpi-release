/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package action

import (
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/denverdino/aliyungo/common"
	"bosh-alicloud-cpi/alicloud"
	"encoding/json"
	"strings"
	"bosh-alicloud-cpi/registry"
)

const (
	DefaultPassword = "Cloud12345"
)

type InstanceProps struct {
	ImageId string 				`json:"image_id"`
	AvailabilityZone string		`json:"availability_zone"`
	InstanceName string 		`json:"instance_name"`
	InstanceChargeType string	`json:"instance_charge_type"`
	InstanceType string 		`json:"instance_type"`
	InstanceRole string			`json:"instance_role"`
	KeyPairName string 			`json:"key_pair"`
	Password string 			`json:"password"`
	EphemeralDisk DiskInfo 		`json:"ephemeral_disk"`
	SystemDisk DiskInfo			`json:"system_disk"`
}

type CreateVMMethod struct {
	CallContext
	stemcells alicloud.StemcellManager
	instances alicloud.InstanceManager
	disks alicloud.DiskManager
	networks alicloud.NetworkManager
	registry registry.Client
}

func NewCreateVMMethod(
	cc CallContext,
	stemcells alicloud.StemcellManager,
	instances alicloud.InstanceManager,
	disks alicloud.DiskManager,
	networks alicloud.NetworkManager,
	registry registry.Client) CreateVMMethod {
	return CreateVMMethod{cc, stemcells,instances, disks, networks, registry}
}

func (a CreateVMMethod) CreateVM(
	agentID apiv1.AgentID, stemcellCID apiv1.StemcellCID,
	cloudProps apiv1.VMCloudProps, networkArgs apiv1.Networks,
	associatedDiskCIDs []apiv1.DiskCID, env apiv1.VMEnv) (apiv1.VMCID, error) {

	cid := apiv1.VMCID{}

	logger := a.Logger

	logger.Info("ENV", "inv: %v", env)
	env2, err := registry.UnmarshalEnvSettings(env)
	if err != nil {
		return cid, a.WrapErrorf(err, "UnmarshalEnvSettings failed %v", env)
	}
	logger.Info("ENV", "inv: %v", env2)

	//
	// convert CloudProps to alicloud dedicated Props
	var args ecs.CreateInstanceArgs

	var instProps InstanceProps
	logger.Info("INPUT", "json %s", cloudProps)
	err = cloudProps.As(&instProps)
	if err != nil {
		return cid, a.WrapErrorf(err, "unmarshal CloudProps failed %v", cloudProps)
	}

	logger.Info("INPUT", "unmarshal CloudProps<Instance>: %s", instProps)

	logger.Info("INPUT", "unmarshal NetworkProps<Instance>: %v", networkArgs)
	networks, err := NewNetworks(networkArgs)
	if err != nil {
		return cid, a.WrapErrorf(err, "create_vm failed when parse Networks %v", networkArgs)
	}

	networks.FillCreateInstanceArgs(&args)

	args.RegionId = common.Region(a.Config.OpenApi.RegionId)
	if instProps.AvailabilityZone != "" {
		args.ZoneId = instProps.AvailabilityZone
	} else {
		args.ZoneId = a.Config.OpenApi.ZoneId
	}
	args.ImageId = stemcellCID.AsString()
	args.UserData = a.Config.Registry.ToInstanceUserData()

	if len(strings.TrimSpace(instProps.KeyPairName)) > 0 {
		args.KeyPairName = instProps.KeyPairName
	} else if len(strings.TrimSpace(instProps.Password)) > 0 {
		args.Password = instProps.Password
	} else {
		args.Password = DefaultPassword
	}

	args.InstanceType = instProps.InstanceType
	args.InstanceChargeType = common.PostPaid
	if strings.Compare(args.InstanceType, "") == 0 {
		args.InstanceType = "ecs.n4.xlarge"
	}

	args.InstanceName = instProps.InstanceName
	args.IoOptimized = "optimized"

	disks, err := NewDisks(instProps.SystemDisk, instProps.EphemeralDisk)
	if err != nil {
		return cid, a.WrapErrorf(err, "bad disks format, %v", instProps)
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

	req, _ := json.Marshal(args)

	logger.Info("OPENAPI", "Args %s", string(req))

	//
	// insert agent re
	agentSettings := registry.AgentSettings{
		AgentID:   agentID.AsString(),
		Blobstore: a.Config.Agent.Blobstore.AsRegistrySettings(),
		Disks:     disks.AsRegistrySettings(),
		Env:       env2,
		Mbus:      a.Config.Agent.Mbus,
		Networks:  networks.AsRegistrySettings(),
		Ntp:       a.Config.Agent.Ntp,
		VM: registry.VMSettings{
			Name: "",
		},
	}

	//if strings.Compare("fake", instProps.InstanceRole) == 0 {
	//	j1, _ := json.Marshal(args)
	//	j2, _ := json.Marshal(agentSettings)
	//	return apiv1.VMCID{}, a。Wr("haltForTest instProps=%v\n networkProps=%v\n args=%s\n registry=%s\n",
	//		instProps, networks, j1, j2)
	//}

	//
	// do create instance
	instCid, err := a.instances.CreateInstance(args)
	if err != nil {
		return apiv1.VMCID{}, a.WrapErrorf(err, "CreateInstance failed INPUT=%s ", string(req))
	}

	agentSettings.VM.Name = instCid

	//
	//
	for _, diskCid := range associatedDiskCIDs {
		err := a.disks.AttachDisk(instCid, diskCid.AsString())
		if err != nil {
			return cid, a.WrapErrorf(err, "associate Persistent Disk error diskCid=%s", diskCid)
		}

		path, err := a.disks.WaitForDiskStatus(diskCid.AsString(), ecs.DiskStatusInUse)
		if err != nil {
			return cid, a.WrapErrorf(err, "associate and WaitForDiskStatus Failed diskCid=%s", diskCid)
		}
		disks.AssociatePersistentDisk(diskCid.AsString(), path)
	}

	err = a.UpdateAgentSettings(instCid, agentSettings)
	if err != nil {
		return apiv1.NewVMCID(instCid), a.WrapErrorf(err, "UpdateAgentSettings Failed %s", )
	}

	err = a.instances.StartInstance(instCid)
	if err != nil {
		return apiv1.NewVMCID(instCid), a.WrapErrorf(err, "StartInstance failed cid =", instCid)
	}

	err = a.instances.WaitForInstanceStatus(instCid, ecs.Running)
	if err != nil {
		return apiv1.NewVMCID(instCid), a.WrapErrorf(err, "StartInstance failed cid=", instCid)
	}

	if networks.HasVip() {
		err = networks.BindInstanceEip(a.Config.NewEcsClient(), instCid, args.RegionId)
		if err != nil {
			return apiv1.NewVMCID(instCid), a.WrapErrorf(err, "StartInstance failed cid=")
		}
	}

	//
	// TODO: every error must free created vm before terminated
	logger.Info("INFO", "FINISH create_vm %s", args)
	return apiv1.NewVMCID(instCid), nil
}

func (a CreateVMMethod) UpdateAgentSettings(instId string, agentSettings registry.AgentSettings) error {
	client := a.registry
	err := client.Update(instId, agentSettings)

	if err != nil {
		json, _ := json.Marshal(agentSettings)
		a.Logger.Error("create_vm", "UpdateAgentSettings to registery failed %s json:%s", json)
		return a.WrapErrorf(err, "UpdateAgentSettings failed %v %s", client, json)
	}

	return nil
}



//
//
//
//func TestCloudProps(t *testing.T) {
//	var cloudProps InstanceProps
//	json.Unmarshal(cloudPropsJson, &cloudProps)
//
//	t.Log(cloudProps)
//	t.Log(cloudProps.EphemeralDisk.GetSizeGB())
//
//	var prop2 InstanceProps
//	json.Unmarshal(cloudPropsJson2, &prop2)
//	t.Log(prop2)
//	t.Log(prop2.EphemeralDisk.GetSizeGB())
//}
