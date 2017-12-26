/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package action

import (
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	"github.com/denverdino/aliyungo/ecs"
	"bosh-alicloud-cpi/alicloud"
	"encoding/json"
	"strings"
	"bosh-alicloud-cpi/registry"
	"fmt"
)

type SpotStrategyType string
type InstanceChargeType string

const (
	PrePaid  = InstanceChargeType("PrePaid")
	PostPaid = InstanceChargeType("PostPaid")
)

//
// Instance properties: ref to docs/bosh/alicloud-cpi.md
// spot https://help.aliyun.com/knowledge_detail/48269.html
// ram profile https://help.aliyun.com/document_detail/54579.html?spm=5176.doc25481.6.797.UVS7aB
type InstanceProps struct {
	EphemeralDisk DiskInfo `json:"ephemeral_disk"`
	SystemDisk    DiskInfo `json:"system_disk"`

	AvailabilityZone string      `json:"availability_zone"`
	InstanceName     string      `json:"instance_name"`
	InstanceType     string      `json:"instance_type"`
	Slbs             []string    `json:"slbs"`
	SlbWeight        json.Number `json:"slb_weight"`
	Password         string      `json:"password"`
	KeyPairName      string      `json:"key_pair_name"`

	ChargeType      string      `json:"charge_type"`
	ChargePeriod    json.Number `json:"charge_period"`
	AutoRenew       string      `json:"auto_renew"`
	AutoRenewPeriod json.Number `json:"auto_renew_period"`
	SpotStrategy 	ecs.SpotStrategyType `json:"spot_strategy"`
	SpotPriceLimit	float64              `json:"spot_price_limit"`
	RamRoleName     string     			 `json:"ram_role_name"`
}

type CreateVMMethod struct {
	CallContext
	stemcells alicloud.StemcellManager
	instances alicloud.InstanceManager
	disks     alicloud.DiskManager
	networks  alicloud.NetworkManager
	registry  registry.Client
}

func NewCreateVMMethod(
	cc CallContext,
	stemcells alicloud.StemcellManager,
	instances alicloud.InstanceManager,
	disks alicloud.DiskManager,
	networks alicloud.NetworkManager,
	registry registry.Client) CreateVMMethod {
	return CreateVMMethod{cc, stemcells, instances, disks, networks, registry}
}

func NewInstanceProps() InstanceProps {
	return InstanceProps{
	}
}

func (a CreateVMMethod) CreateVM(
	agentID apiv1.AgentID, stemcellCID apiv1.StemcellCID,
	cloudProps apiv1.VMCloudProps, networkArgs apiv1.Networks,
	associatedDiskCIDs []apiv1.DiskCID, env apiv1.VMEnv) (apiv1.VMCID, error) {

	cid := apiv1.VMCID{}

	//
	// convert registry env
	registryEnv, err := registry.UnmarshalEnvSettings(env)
	if err != nil {
		return cid, a.WrapErrorf(err, "unmarshal EnvSettings failed %v", env)
	}

	//
	// convert CloudProps to alicloud dedicated Props
	instProps := NewInstanceProps()
	err = cloudProps.As(&instProps)
	if err != nil {
		return cid, a.WrapErrorf(err, "unmarshal instance cloud_properties failed %v", cloudProps)
	}

	//
	// parse networks from networkArgs
	networks, err := NewNetworks(networkArgs)
	if err != nil {
		return cid, a.WrapErrorf(err, "parse network cloud_properties %v", networkArgs)
	}

	var args ecs.CreateInstanceArgs
	networks.FillCreateInstanceArgs(&args)

	//
	// assign zone
	if instProps.AvailabilityZone != "" {
		args.ZoneId = instProps.AvailabilityZone
	} else {
		args.ZoneId = a.Config.OpenApi.ZoneId
	}
	if args.ZoneId == "" {
		return cid, a.Errorf("can't get zone from availability_zone or cpi.config")
	}

	//
	// config instance_type
	args.InstanceType = instProps.InstanceType
	if args.InstanceType == "" {
		return cid, a.Errorf("missing instance_type")
	}

	// spot props
	if err := validateSpotProps(instProps); err != nil {
		return apiv1.VMCID{}, a.WrapErrorf(err, "invalid spot properties ")
	}

	//
	// config vm charge type
	if instProps.ChargeType == "PrePaid" {
		args.InstanceChargeType = "PrePaid"
		period, err := instProps.ChargePeriod.Int64()
		if err != nil {
			return cid, a.WrapErrorf(err, "parse charge_period %s failed when charge_type is `PrePaid`", instProps.ChargePeriod.String())
		}
		args.Period = int(period)
		if strings.EqualFold(instProps.AutoRenew, "True") {
			args.AutoRenew = true
			period, err = instProps.AutoRenewPeriod.Int64()
			if err != nil {
				return cid, a.WrapErrorf(err, "parse charge_auto_renew_period %s failed when charge_auto_renew is `True`", instProps.AutoRenewPeriod.String())
			}
			args.AutoRenewPeriod = int(period)
		} else if strings.EqualFold(instProps.AutoRenew, "False") || instProps.AutoRenew == "" {
			args.AutoRenew = false
		} else {
			return cid, a.Errorf("unexpected charge_auto_renew: %s", instProps.AutoRenew)
		}
	} else if instProps.ChargeType == "PostPaid" || instProps.ChargeType == "" {
		args.InstanceChargeType = "PostPaid"
	} else {
		return cid, a.Errorf("unexpected charge type %s", instProps.ChargeType)
	}

	//
	// compare key pair or password
	if len(strings.TrimSpace(instProps.KeyPairName)) > 0 {
		args.KeyPairName = instProps.KeyPairName
	} else if len(strings.TrimSpace(instProps.Password)) > 0 {
		args.Password = instProps.Password
	}

	//
	args.RegionId = a.Config.OpenApi.GetRegion()
	args.ImageId = stemcellCID.AsString()
	args.InstanceName = instProps.InstanceName
	args.IoOptimized = "optimized"
	args.UserData = a.Config.Registry.ToInstanceUserData()
	args.SpotStrategy = instProps.SpotStrategy
	args.SpotPriceLimit = instProps.SpotPriceLimit
	args.RamRoleName = instProps.RamRoleName

	//
	// fill disks
	disks, err := NewDisksWithProps(instProps.SystemDisk, instProps.EphemeralDisk)
	if err != nil {
		return cid, a.WrapErrorf(err, "bad disks format, %v", instProps)
	}
	disks.FillCreateInstanceArgs(&args)
	//
	// do CreateInstance !!!
	instCid, err := a.instances.CreateInstance(args)
	if err != nil {
		req, _ := json.Marshal(args)
		return apiv1.VMCID{}, a.WrapErrorf(err, "create instance failed with input=%s ", string(req))
	}

	//
	// Wait for the instance status to STOPPED
	err = a.instances.ChangeInstanceStatus(instCid, ecs.Stopped, func(status ecs.InstanceStatus) (bool, error) {
		switch status {
		case ecs.Stopped:
			return true, nil
		case ecs.Creating:
			return false, nil
		case ecs.Pending:
			return false, nil
		default:
			return false, fmt.Errorf("unexcepted status %s", status)
		}
	})

	if err != nil {
		return apiv1.VMCID{}, a.WrapErrorf(err, "wait %s to STOPPED failed", instCid)
	}

	agentSettings := registry.AgentSettings{
		AgentID:   agentID.AsString(),
		Blobstore: a.Config.Agent.Blobstore.AsRegistrySettings(),
		Disks:     disks.AsRegistrySettings(),
		Env:       registryEnv,
		Mbus:      a.Config.Agent.Mbus,
		Networks:  networks.AsRegistrySettings(),
		Ntp:       a.Config.Agent.Ntp,
		VM: registry.VMSettings{
			Name: instCid,
		},
	}

	//
	// updateInstance
	err = a.updateInstance(instCid, associatedDiskCIDs, instProps, networks, disks, agentSettings)

	//
	// for every error must free created vm before terminated
	if err != nil {
		err2 := a.instances.ChangeInstanceStatus(instCid, ecs.Deleted, func(status ecs.InstanceStatus) (bool, error) {
			switch status {
			case ecs.Running:
				return false, a.instances.StopInstance(instCid)
			case ecs.Stopped:
				return false, a.instances.DeleteInstance(instCid)
			case ecs.Deleted:
				return true, nil
			case ecs.Pending:
				return false, nil
			case ecs.Creating:
				return false, nil
			case ecs.Starting:
				return false, nil
			case ecs.Stopping:
				return false, nil
			default:
				return false, fmt.Errorf("unexpect %s for ReleaseInstance", status)
			}
		})
		if err2 == nil {
			return apiv1.NewVMCID(instCid), a.WrapErrorf(err, "update %s failed vm deleted", instCid)
		} else {
			return apiv1.NewVMCID(instCid), a.WrapErrorf(err, "update %s failed delete failed %v", instCid, err2)
		}
	}

	return apiv1.NewVMCID(instCid), nil
}

func (a CreateVMMethod) updateInstance(instCid string, associatedDiskCIDs []apiv1.DiskCID, instProps InstanceProps, networks Networks, disks Disks, agentSettings registry.AgentSettings) (error) {
	//
	// join instance to multiple security groups
	if len(networks.privateProps.SecurityGroupIds) > 0 {
		inst, err := a.instances.GetInstance(instCid)
		if err != nil {
			return a.WrapErrorf(err, "GetInstance instanceId=%s falied.", instCid)
		}
		for _, group := range networks.privateProps.SecurityGroupIds {
			find := false
			for _, gr := range inst.SecurityGroupIds.SecurityGroupId {
				if group == gr {
					find = true
					break
				}
			}
			if !find {
				if err := a.networks.JoinSecurityGroup(instCid, group); err != nil {
					return a.WrapErrorf(err, "Instance %s JoinSecurityGroup %s failed.", instCid, group)
				}
			}
		}
	}

	// associate persistent disks,
	// TODO: use ChangeDiskStatus to avoid failed
	for _, diskCid := range associatedDiskCIDs {
		err := a.disks.AttachDisk(instCid, diskCid.AsString())
		if err != nil {
			return a.WrapErrorf(err, "associate Persistent Disk error diskCid=%s", diskCid)
		}

		path, err := a.disks.WaitForDiskStatus(diskCid.AsString(), ecs.DiskStatusInUse)
		if err != nil {
			return a.WrapErrorf(err, "associate and WaitForDiskStatus Failed diskCid=%s", diskCid)
		}
		disks.AssociatePersistentDisk(diskCid.AsString(), path)
	}

	//
	// put agent settings.json to registry
	err := a.UpdateAgentSettings(instCid, agentSettings)
	if err != nil {
		return a.WrapErrorf(err, "UpdateAgentSettings Failed %s", )
	}

	//
	// wait for instance to start
	err = a.instances.ChangeInstanceStatus(instCid, ecs.Running, func(status ecs.InstanceStatus) (bool, error) {
		switch status {
		case ecs.Stopped:
			return false, a.instances.StartInstance(instCid)
		case ecs.Pending:
			return false, nil
		case ecs.Starting:
			return false, nil
		case ecs.Running:
			return true, nil
		default:
			return false, fmt.Errorf("unexcepted status %s", status)
		}
	})

	if err != nil {
		return a.WrapErrorf(err, "change %s to Running failed", instCid)
	}

	for _, eip := range networks.GetVips() {
		err := a.networks.BindEip(instCid, eip)
		if err != nil {
			return a.WrapErrorf(err, "bind eip %s to %s failed", eip, instCid)
		}
	}

	slbWeight, err := instProps.SlbWeight.Int64()
	if err != nil {
		slbWeight = alicloud.DefaultSlbWeight
	} else if slbWeight == 0 {
		slbWeight = alicloud.DefaultSlbWeight
	}

	for _, slb := range instProps.Slbs {
		err := a.networks.BindSLB(instCid, slb, int(slbWeight))
		if err != nil {
			return a.WrapErrorf(err, "bind %s to slb %s failed ", instCid, slb)
		}
	}
	return nil
}

func validateSpotProps(p InstanceProps) error {
	strategy := string(p.SpotStrategy)
	strategyArr := []string{string(ecs.NoSpot), string(ecs.SpotWithPriceLimit), string(ecs.SpotAsPriceGo)}
	limitPrice := float64(p.SpotPriceLimit)

	if limitPrice == 0 && strategy == "" {
		return nil
	}

	if p.ChargeType == string(PrePaid) {
		return fmt.Errorf("the spot strategy only support 'PostPaid' instance charge type")
	}

	if err := validAllowedStringValues(strategy, strategyArr); err != nil {
		return err
	}

	if limitPrice != 0 && strategy != string(ecs.SpotWithPriceLimit) {
		return fmt.Errorf("spot limit price only support 'SpotWithPriceLimit' strategy")
	}
	return nil
}

func (a CreateVMMethod) UpdateAgentSettings(instId string, agentSettings registry.AgentSettings) error {
	client := a.registry
	err := client.Update(instId, agentSettings)

	if err != nil {
		txt, _ := json.Marshal(agentSettings)
		a.Logger.Error("create_vm", "UpdateAgentSettings to registry failed %s json:%s", txt)
		return a.WrapErrorf(err, "UpdateAgentSettings failed %v %s", client, txt)
	}

	return nil
}
