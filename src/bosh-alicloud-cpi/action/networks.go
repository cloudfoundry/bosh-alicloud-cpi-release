package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"bosh-alicloud-cpi/registry"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/denverdino/aliyungo/common"
	"strings"
)

type Networks struct {
	networks apiv1.Networks
	main apiv1.Network
	props NetworkProps
}


type NetworkProps struct {
	SecurityGroupId string		`json:"security_group_id,omitempty"`
	VSwitchId string			`json:"vswitch_id,omitempty"`
	InternetChargeType string	`json:"internet_charge_type,omitempty"`
}

func NewNetworks(args apiv1.Networks) (Networks, error) {
	a := Networks{
		networks: args,
		main: nil,
		props: NetworkProps{},
	}
	var p *apiv1.Network
	for k, v := range args {
		if strings.Compare("private", k) == 0 {
			p = &v
		}

		if strings.Compare("default", k) == 0 {
			p = &v
		}
	}

	if p == nil {
		return a, bosherr.Errorf("No main network find %v", args)
	}

	a.main = *p
	err := a.main.CloudProps().As(&a.props)
	if err != nil {
		return a, bosherr.WrapErrorf(err, "unmarshal json failed, %v", a.main.CloudProps())
	}

	return a, nil
}

func (a Networks) FillCreateInstanceArgs(args *ecs.CreateInstanceArgs) (error) {
	if !a.main.IsDynamic() {
		args.PrivateIpAddress = a.GetPrivateAddress()
		args.VSwitchId = a.GetVSwitchId()
		args.SecurityGroupId = a.GetSecurityGroupId()
		args.InternetChargeType = a.GetInternetChargeType()
		return nil
	} else {
		return bosherr.Errorf("NOT IMPLEMENTED Dynamic Networks")
	}
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

func (a Networks) GetPrivateAddress() (string) {
	return a.main.IP()
}

func (a Networks) GetVSwitchId() (string) {
	return a.props.VSwitchId
}

func (a Networks) GetSecurityGroupId() (string) {
	return a.props.SecurityGroupId
}

func (a Networks) GetInternetMaxBandwidthIn() (int) {
	return 100
}

func (a Networks) GetInternetMaxBandwidthOut() (int) {
	return 100
}

func (a Networks) GetInternetChargeType() (common.InternetChargeType) {
	return common.InternetChargeType(a.props.InternetChargeType)
}