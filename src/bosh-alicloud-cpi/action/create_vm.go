/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package action

import (
	"bosh-alicloud-cpi/alicloud"
	"bosh-alicloud-cpi/registry"

	//"github.com/cloudfoundry-incubator/bosh-alicloud-cpi-release/src/bosh-alicloud-cpi/registry"
	"encoding/json"
	"fmt"
	"strings"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/cppforlife/bosh-cpi-go/apiv1"

	"encoding/base64"

	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
)

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
	EphemeralDisk DiskInfo          `json:"ephemeral_disk"`
	SystemDisk    DiskInfo          `json:"system_disk"`
	Tags          map[string]string `json:"tags"`

	Region               string                    `json:"region"`
	AvailabilityZone     string                    `json:"availability_zone"`
	InstanceName         string                    `json:"instance_name"`
	InstanceType         string                    `json:"instance_type"`
	SlbServerGroupWeight json.Number               `json:"slb_server_group_weight"`
	SlbServerGroupPort   json.Number               `json:"slb_server_group_port"`
	SlbServerGroup       []string                  `json:"slb_server_group"`
	Slbs                 []string                  `json:"slbs"`
	SlbWeight            json.Number               `json:"slb_weight"`
	Password             string                    `json:"password"`
	KeyPairName          string                    `json:"key_pair_name"`
	SecurityGroupIds     []string                  `json:"security_group_ids"`
	ChargeType           string                    `json:"charge_type"`
	ChargePeriod         json.Number               `json:"charge_period"`
	ChargePeriodUnit     string                    `json:"charge_period_unit"`
	AutoRenew            string                    `json:"auto_renew"`
	AutoRenewPeriod      json.Number               `json:"auto_renew_period"`
	SpotStrategy         alicloud.SpotStrategyType `json:"spot_strategy"`
	SpotPriceLimit       float64                   `json:"spot_price_limit"`
	RamRoleName          string                    `json:"ram_role_name"`
	StemcellId           string                    `json:"stemcell_id"`
}

type CreateVMMethod struct {
	CallContext
	stemcells alicloud.StemcellManager
	instances alicloud.InstanceManager
	disks     alicloud.DiskManager
	networks  alicloud.NetworkManager
	registry  alicloud.RegistryManager
}

func NewCreateVMMethod(
	cc CallContext,
	stemcells alicloud.StemcellManager,
	instances alicloud.InstanceManager,
	disks alicloud.DiskManager,
	networks alicloud.NetworkManager,
	registry alicloud.RegistryManager) CreateVMMethod {
	return CreateVMMethod{cc, stemcells, instances, disks, networks, registry}
}

func NewInstanceProps() InstanceProps {
	return InstanceProps{}
}
func (a CreateVMMethod) CreateVM(
	agentID apiv1.AgentID, stemcellCID apiv1.StemcellCID,
	cloudProps apiv1.VMCloudProps, networkArgs apiv1.Networks,
	associatedDiskCIDs []apiv1.DiskCID, env apiv1.VMEnv) (apiv1.VMCID, error) {
	vmCid, _, err := a.createVM(agentID, stemcellCID, cloudProps, networkArgs, associatedDiskCIDs, env)
	return vmCid, err
}

func (a CreateVMMethod) CreateVMV2(
	agentID apiv1.AgentID, stemcellCID apiv1.StemcellCID,
	cloudProps apiv1.VMCloudProps, networkArgs apiv1.Networks,
	associatedDiskCIDs []apiv1.DiskCID, env apiv1.VMEnv) (apiv1.VMCID, apiv1.Networks, error) {
	return a.createVM(agentID, stemcellCID, cloudProps, networkArgs, associatedDiskCIDs, env)
}

func (a CreateVMMethod) createVM(
	agentID apiv1.AgentID, stemcellCID apiv1.StemcellCID,
	cloudProps apiv1.VMCloudProps, networkArgs apiv1.Networks,
	associatedDiskCIDs []apiv1.DiskCID, env apiv1.VMEnv) (apiv1.VMCID, apiv1.Networks, error) {

	cid := apiv1.VMCID{}

	// convert registry env
	registryEnv, err := registry.UnmarshalEnvSettings(env)
	if err != nil {
		return cid, nil, bosherr.WrapErrorf(err, "unmarshal EnvSettings failed %v", env)
	}

	// convert CloudProps to alicloud dedicated Props
	instProps := NewInstanceProps()
	err = cloudProps.As(&instProps)
	if err != nil {
		return cid, nil, bosherr.WrapErrorf(err, "unmarshal instance cloud_properties failed %v", cloudProps)
	}

	// if cross region deployment, bosh region and stemcell id must be specified
	if strings.TrimSpace(instProps.Region) != "" && strings.TrimSpace(instProps.StemcellId) == "" {
		return cid, nil, bosherr.Error("'stemcell_id' must be specified when self-defined region is specified.")
	}
	// parse networks from networkArgs
	networks, err := NewNetworks(networkArgs)
	if err != nil {
		return cid, nil, bosherr.WrapErrorf(err, "parse network cloud_properties %v", networkArgs)
	}

	//Security groups can be specified as follows, ordered by greatest precedence: vm_types, followed by networks.
	if len(instProps.SecurityGroupIds) > 0 {
		networks.privateProps.SecurityGroupIds = instProps.SecurityGroupIds
	}

	args := ecs.CreateCreateInstanceRequest()
	if err := networks.FillCreateInstanceArgs(args); err != nil {
		return cid, nil, bosherr.WrapErrorf(err, "fill instance network args failed and args: %v", networks.privateProps)
	}

	if instProps.AvailabilityZone != "" {
		args.ZoneId = instProps.AvailabilityZone
	} else {
		args.ZoneId = a.Config.OpenApi.GetAvailabilityZone()
	}
	if args.ZoneId == "" {
		return cid, nil, bosherr.Errorf("can't get zone from availability_zone or cpi.config")
	}

	// config instance_type
	args.InstanceType = instProps.InstanceType
	if args.InstanceType == "" {
		return cid, nil, bosherr.Errorf("missing instance_type")
	}

	// spot props
	if err := validateSpotProps(instProps); err != nil {
		return apiv1.VMCID{}, nil, bosherr.WrapError(err, "invalid spot properties ")
	}

	// config vm charge type
	if instProps.ChargeType == "PrePaid" {
		args.InstanceChargeType = "PrePaid"
		period, err := instProps.ChargePeriod.Int64()
		if err != nil {
			return cid, nil, bosherr.WrapErrorf(err, "parse charge_period %s failed when charge_type is `PrePaid`", instProps.ChargePeriod.String())
		}
		args.Period = requests.NewInteger64(period)
		args.PeriodUnit = instProps.ChargePeriodUnit
		if strings.EqualFold(instProps.AutoRenew, "True") {
			args.AutoRenew = requests.NewBoolean(true)
			period, err = instProps.AutoRenewPeriod.Int64()
			if err != nil {
				return cid, nil, bosherr.WrapErrorf(err, "parse charge_auto_renew_period %s failed when charge_auto_renew is `True`", instProps.AutoRenewPeriod.String())
			}
			args.AutoRenewPeriod = requests.NewInteger64(period)
		} else if strings.EqualFold(instProps.AutoRenew, "False") || instProps.AutoRenew == "" {
			args.AutoRenew = requests.NewBoolean(false)
		} else {
			return cid, nil, bosherr.Errorf("unexpected charge_auto_renew: %s", instProps.AutoRenew)
		}
	} else if instProps.ChargeType == "PostPaid" || instProps.ChargeType == "" {
		args.InstanceChargeType = "PostPaid"
	} else {
		return cid, nil, bosherr.Errorf("unexpected charge type %s", instProps.ChargeType)
	}

	// compare key pair or password
	if len(strings.TrimSpace(instProps.KeyPairName)) > 0 {
		args.KeyPairName = instProps.KeyPairName
	} else if len(strings.TrimSpace(instProps.Password)) > 0 {
		args.Password = instProps.Password
	}

	args.ImageId = stemcellCID.AsString()
	if instProps.StemcellId != "" {
		args.ImageId = instProps.StemcellId
	}
	args.InstanceName = instProps.InstanceName
	args.IoOptimized = "optimized"
	if a.Config.Registry.ToInstanceUserData() != "" {
		args.UserData = base64.StdEncoding.EncodeToString([]byte(a.Config.Registry.ToInstanceUserData()))
	}
	args.SpotStrategy = string(instProps.SpotStrategy)
	args.SpotPriceLimit = requests.NewFloat(instProps.SpotPriceLimit)
	args.RamRoleName = instProps.RamRoleName

	// fill disks
	disks, err := NewDisksWithProps(instProps.SystemDisk, instProps.EphemeralDisk)
	if err != nil {
		return cid, nil, bosherr.WrapErrorf(err, "bad disks format, %v", instProps)
	}

	disks.FillCreateInstanceArgs(a.Config.OpenApi.Encrypted, args)

	//打标签
	tags := make([]ecs.CreateInstanceTag, 0)
	// 首先获取registry中的tag
	for k, v := range registryEnv.Bosh.Tags {
		tags = append(tags, ecs.CreateInstanceTag{
			Key:   k,
			Value: fmt.Sprint(v),
		})
	}
	// 为了支持CR，tag中添加创建时获取的env.bosh.group参数
	if registryEnv.Bosh.Group != "" {
		tags = append(tags, ecs.CreateInstanceTag{
			Key: "environment_bosh_group_for_using_capacity_reservation_in_alicloud_iaas",
			Value: registryEnv.Bosh.Group[:120],
		})
	}
	// 接下来获取manifest中的tag
	for k, v := range instProps.Tags {
		tags = append(tags, ecs.CreateInstanceTag{
			Key:   k,
			Value: v,
		})
	}
	args.Tag = &tags

	// do CreateInstance !!!
	instCid, err := a.instances.CreateInstance(instProps.Region, args)
	if err != nil {
		req, _ := json.Marshal(args)
		return apiv1.VMCID{}, nil, bosherr.WrapErrorf(err, "create instance failed with input=%s ", string(req))
	}

	// Wait for the instance status to STOPPED
	err = a.instances.ChangeInstanceStatus(instCid, alicloud.Stopped, func(status alicloud.InstanceStatus) (bool, error) {
		switch status {
		case alicloud.Stopped:
			return true, nil
		case alicloud.Running:
			return false, a.instances.StopInstance(instCid)
		default:
			return false, nil
		}
	})

	if err != nil {
		eniIds := a.instances.GetAttachedNetworkInterfaceIds(instCid)
		var err2 error
		for retry := 0; retry < 10; retry++ {
			err2 = a.instances.ChangeInstanceStatus(instCid, alicloud.Deleted, func(status alicloud.InstanceStatus) (bool, error) {
				switch status {
				case alicloud.Running, alicloud.Stopped:
					return false, a.instances.DeleteInstance(instCid)
				case alicloud.Deleted:
					return true, a.instances.CleanupInstanceNetworkInterfaces(instCid, eniIds)
				default:
					return false, nil
				}
			})
			if err2 == nil {
				return apiv1.NewVMCID(instCid), nil, bosherr.WrapErrorf(err, "wait %s to STOPPED failed and the vm has been deleted.", instCid)
			}
			time.Sleep(5 * time.Second)
		}
		return apiv1.VMCID{}, nil, bosherr.WrapErrorf(err, "wait %s to STOPPED failed and then delete it timeout: %v", instCid, err2)
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

	// for every error must free created vm before terminated
	if err != nil {
		eniIds := a.instances.GetAttachedNetworkInterfaceIds(instCid)
		var err2 error
		for retry := 0; retry < 10; retry++ {
			err2 = a.instances.ChangeInstanceStatus(instCid, alicloud.Deleted, func(status alicloud.InstanceStatus) (bool, error) {
				switch status {
				case alicloud.Running, alicloud.Stopped:
					return false, a.instances.DeleteInstance(instCid)
				case alicloud.Deleted:
					return true, a.instances.CleanupInstanceNetworkInterfaces(instCid, eniIds)
				default:
					return false, nil
				}
			})
			if err2 == nil {
				return apiv1.NewVMCID(instCid), nil, bosherr.WrapErrorf(err, "update %s failed and the vm has been deleted.", instCid)
			}
			time.Sleep(5 * time.Second)
		}
		return apiv1.NewVMCID(instCid), nil, bosherr.WrapErrorf(err, "update %s failed and then delete it timeout: %v", instCid, err2)
	}

	return apiv1.NewVMCID(instCid), networkArgs, nil
}

func (a CreateVMMethod) updateInstance(instCid string, associatedDiskCIDs []apiv1.DiskCID, instProps InstanceProps, networks Networks, disks Disks, agentSettings registry.AgentSettings) error {
	//
	// join instance to multiple security groups
	if len(networks.privateProps.SecurityGroupIds) > 0 {
		inst, err := a.instances.GetInstance(instCid)
		if err != nil {
			return bosherr.WrapErrorf(err, "GetInstance instanceId=%s falied.", instCid)
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
				if err := a.networks.JoinSecurityGroup(instProps.Region, instCid, group); err != nil {
					return bosherr.WrapErrorf(err, "Instance %s JoinSecurityGroup %s failed.", instCid, group)
				}
			}
		}
	}

	// associate persistent disks,
	// TODO: use ChangeDiskStatus to avoid failed
	for _, diskCid := range associatedDiskCIDs {
		err := a.disks.AttachDisk(instCid, diskCid.AsString())
		if err != nil {
			return bosherr.WrapErrorf(err, "associate Persistent Disk error diskCid=%s", diskCid)
		}

		path, err := a.disks.WaitForDiskStatus(diskCid.AsString(), alicloud.DiskStatusInUse)
		if err != nil {
			return bosherr.WrapErrorf(err, "associate and WaitForDiskStatus Failed diskCid=%s", diskCid)
		}
		disks.AssociatePersistentDisk(diskCid.AsString(), path)
	}

	// put agent settings.json to registry
	err := a.UpdateAgentSettings(instCid, agentSettings)
	if err != nil {
		return bosherr.WrapError(err, "UpdateAgentSettings Failed.")
	}

	// wait for instance to start
	err = a.instances.ChangeInstanceStatus(instCid, alicloud.Running, func(status alicloud.InstanceStatus) (bool, error) {
		switch status {
		case alicloud.Stopped:
			return false, a.instances.StartInstance(instCid)
		case alicloud.Running:
			return true, nil
		default:
			return false, nil
		}
	})

	if err != nil {
		return bosherr.WrapErrorf(err, "change %s to Running failed", instCid)
	}

	for _, eip := range networks.GetVips() {
		err := a.networks.BindEip(instProps.Region, instCid, eip)
		if err != nil {
			return bosherr.WrapErrorf(err, "bind eip %s to %s failed", eip, instCid)
		}
	}
	slbWeight, err := instProps.SlbWeight.Int64()
	if err != nil {
		slbWeight = alicloud.DefaultSlbWeight
	} else if slbWeight == 0 {
		slbWeight = alicloud.DefaultSlbWeight
	}
	for _, slb := range instProps.Slbs {
		err := a.networks.BindSLB(instProps.Region, instCid, slb, int(slbWeight))
		if err != nil {
			return bosherr.WrapErrorf(err, "bind %s to slb %s failed ", instCid, slb)
		}
	}
	slbServerGroupPort, err := instProps.SlbServerGroupPort.Int64()
	if err != nil {
		slbServerGroupPort = alicloud.DefaultSlbServerGroupPort
	} else if slbServerGroupPort == 0 {
		slbServerGroupPort = alicloud.DefaultSlbServerGroupPort
	}
	slbServerGroupWeight, err := instProps.SlbServerGroupWeight.Int64()
	if err != nil {
		slbServerGroupWeight = alicloud.DefaultSlbServerGroupWeight
	} else if slbServerGroupWeight == 0 {
		slbServerGroupWeight = alicloud.DefaultSlbServerGroupWeight
	}
	for _, slbServerGroup := range instProps.SlbServerGroup {
		err := a.networks.BindSlbServerGroup(instProps.Region, instCid, slbServerGroup, int(slbServerGroupWeight), int(slbServerGroupPort))
		if err != nil {
			return bosherr.WrapErrorf(err, "bind %s to slbServerGroup %s failed,weight:%d,port:%d ", instCid, slbServerGroup, slbServerGroupWeight, slbServerGroupPort)
		}
	}
	return nil
}

func validateSpotProps(p InstanceProps) error {
	strategy := string(p.SpotStrategy)
	strategyArr := []string{string(alicloud.NoSpot), string(alicloud.SpotWithPriceLimit), string(alicloud.SpotAsPriceGo)}
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

	if limitPrice != 0 && strategy != string(alicloud.SpotWithPriceLimit) {
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
		return bosherr.WrapErrorf(err, "UpdateAgentSettings failed %v %s", client, txt)
	}

	return nil
}
