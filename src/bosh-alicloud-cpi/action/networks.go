/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package action

import (
	"bosh-alicloud-cpi/registry"
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
)

const (
	NetworkTypeManual  = "manual"
	NetworkTypeDynamic = "dynamic"
	NetworkTypeVip     = "vip"
)

type Networks struct {
	networks     apiv1.Networks
	private      apiv1.Network
	privateProps NetworkProps
	vips         []apiv1.Network
}

type NetworkProps struct {
	SecurityGroupId    string   `json:"security_group_id"`
	SecurityGroupIds   []string `json:"security_group_ids"`
	VSwitchId          string   `json:"vswitch_id"`
	InternetChargeType string   `json:"internet_charge_type,omitempty"`
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

func (a Networks) HasVip() bool {
	return len(a.vips) > 0
}

func (a Networks) FillCreateInstanceArgs(args *ecs.CreateInstanceRequest) error {
	props := a.privateProps

	if props.VSwitchId == "" {
		return fmt.Errorf("unexpected empty vswitch_id")
	}
	if props.SecurityGroupId == "" && len(props.SecurityGroupIds) < 1 {
		return fmt.Errorf("unexpected empty security_group_ids")
	}

	if props.SecurityGroupId != "" {
		args.SecurityGroupId = props.SecurityGroupId
	} else {
		args.SecurityGroupId = props.SecurityGroupIds[0]
	}
	args.VSwitchId = props.VSwitchId
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

func (a Networks) GetVips() []string {
	result := []string{}
	for _, net := range a.networks {
		if net.Type() == "vip" {
			result = append(result, net.IP())
		}
	}
	return result
}

func (a Networks) AsRegistrySettings() registry.NetworksSettings {
	r := map[string]registry.NetworkSettings{}

	for k, v := range a.networks {
		r[k] = registry.NetworkSettings{
			Type:    v.Type(),
			IP:      v.IP(),
			Netmask: v.Netmask(),
			Gateway: v.Gateway(),
			DNS:     v.DNS(),
			DHCP:    false,
			Default: v.Default(),
			// CloudProperties: v.CloudProps()
		}
	}
	return r
}
