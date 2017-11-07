/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"bosh-alicloud-cpi/registry"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/denverdino/aliyungo/common"
	"fmt"
)

const (
	NetworkTypeManual = "manual"
	NetworkTypeDynamic = "dynamic"
	NetworkTypeVip = "vip"
)

type Networks struct {
	networks	apiv1.Networks
	private   	apiv1.Network
	privateProps NetworkProps
	vips		[]apiv1.Network
}

type NetworkProps struct {
	SecurityGroupId string		`json:"security_group_id,omitempty"`
	VSwitchId string			`json:"vswitch_id,omitempty"`
	InternetChargeType string	`json:"internet_charge_type,omitempty"`
}

func NewNetworks(args apiv1.Networks) (Networks, error) {
	r := Networks{networks: args}

	for _, v := range args {
		switch v.Type() {
		case NetworkTypeManual:
			if r.private == nil {
				r.private = v
			} else {
				return r, fmt.Errorf("support only one private Network")
			}
		case NetworkTypeDynamic:
			if r.private == nil {
				r.private = v
			} else {
				return r, fmt.Errorf("support only one private Network")
			}
		case NetworkTypeVip:
			r.vips = append(r.vips, v)
		default:
			return r, fmt.Errorf("unknown Network type: %s", v.Type())
		}
	}

	if r.private == nil {
		return r, fmt.Errorf("no private Networks")
	}

	err := r.private.CloudProps().As(&r.privateProps)
	if err != nil {
		return r, bosherr.WrapErrorf(err, "unmarshal cloudProps failed, %v", r.private.CloudProps())
	}

	return r, nil
}

func (a Networks) HasVip() (bool) {
	return len(a.vips) > 0
}

func (a Networks) FillCreateInstanceArgs(args *ecs.CreateInstanceArgs) (error) {
	props := a.privateProps

	if props.VSwitchId == "" {
		return fmt.Errorf("unexpected empty VSwitchId")
	}
	if props.VSwitchId == "" {
		return fmt.Errorf("unexpected empty SecurityGroupId")
	}

	args.VSwitchId = props.VSwitchId
	args.SecurityGroupId = props.SecurityGroupId
	// args.InternetChargeType = common.InternetChargeType(props.InternetChargeType)

	// TODO no need to add
	// args.InternetMaxBandwidthIn = a.GetInternetMaxBandwidthIn()
	// args.InternetMaxBandwidthOut = a.GetInternetMaxBandwidthOut()

	private := a.private
	if private.Type() == NetworkTypeManual {
		if private.IsDynamic() {
			return fmt.Errorf("confilct! manual Network IsDynamic")
		}
		args.PrivateIpAddress = a.private.IP()
	} else if private.Type() == NetworkTypeDynamic {
		if !private.IsDynamic() {
			return fmt.Errorf("confilct! dynamic Network IsDynamic=false")
		}
		// Nothing to do for dynamic Network?
	}
	return nil
}

//
// TODO move to alicloud.NetworkManager
func (a Networks) BindInstanceEip(client *ecs.Client, instanceId string, regionId common.Region) (error) {
	net := a.VipNetwork()

	if net == nil{
		return nil
	}

	var allocationId string
	var arg ecs.DescribeEipAddressesArgs
	arg.EipAddress = net.IP()
	arg.RegionId = regionId

	eipAddresses, _, _:=client.DescribeEipAddresses(&arg)
	for _, address := range eipAddresses {
		allocationId = address.AllocationId
	}

	if err := client.AssociateEipAddress(allocationId, instanceId); err != nil {
		return err
	}

	if err := client.WaitForEip(regionId, allocationId, ecs.EipStatusInUse, 60); err != nil {
		return fmt.Errorf("error waitting for EIP allocated: %#v", err)
	}
	return nil
}

func (a Networks) VipNetwork() (apiv1.Network) {
	for _, net := range a.networks {
		if net.Type() == "vip" {
			return net
		}
	}
	return nil
}

func (a Networks) AsRegistrySettings() (registry.NetworksSettings) {
	r := map[string]registry.NetworkSettings{}

	for k, v := range a.networks {
		r[k] = registry.NetworkSettings{
			Type: v.Type(),
			IP: v.IP(),
			Netmask: v.Netmask(),
			Gateway: v.Gateway(),
			DNS: v.DNS(),
			DHCP: false,
			Default: v.Default(),
			// CloudProperties: v.CloudProps()
		}
	}
	//r := map[string]registry.NetworkSettings {
	//	"private": { Type: "vip"},
	//	"public" : {Type: "vip"},
	//}
	return r
}