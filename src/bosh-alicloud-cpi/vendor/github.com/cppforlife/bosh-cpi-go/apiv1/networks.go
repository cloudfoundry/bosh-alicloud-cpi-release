package apiv1

import (
	"encoding/json"
	"fmt"
	gonet "net"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type Networks map[string]Network

type Network interface {
	Type() string

	IP() string
	Netmask() string
	Gateway() string

	DNS() []string
	Default() []string

	SetMAC(string)
	SetDNS([]string)
	SetPreconfigured()

	CloudProps() NetworkCloudProps

	// Misc
	IsDynamic() bool
	IsDefaultFor(string) bool
	IPWithSubnetMask() string

	_final() // interface unimplementable from outside
}

type NetworkCloudProps interface {
	As(interface{}) error
}

type NetworkImpl struct {
	spec networkSpec2

	mac           string
	preconfigured bool
}

type networkSpec2 struct {
	Type string

	IP      string
	Netmask string
	Gateway string

	DNS     []string
	Default []string

	CloudProps CloudPropsImpl `json:"cloud_properties"`
}

type NetworkOpts struct {
	Type string

	IP      string
	Netmask string
	Gateway string

	DNS     []string
	Default []string

	CloudProps CloudPropsImpl `json:"cloud_properties"`
}

func NewNetwork(opts NetworkOpts) Network {
	return &NetworkImpl{
		spec: networkSpec2{
			Type: opts.Type,

			IP:      opts.IP,
			Netmask: opts.Netmask,
			Gateway: opts.Gateway,

			DNS:     opts.DNS,
			Default: opts.Default,
		},
	}
}

func (n NetworkImpl) Type() string { return n.spec.Type }

func (n NetworkImpl) IP() string      { return n.spec.IP }
func (n NetworkImpl) Netmask() string { return n.spec.Netmask }
func (n NetworkImpl) Gateway() string { return n.spec.Gateway }

func (n NetworkImpl) DNS() []string     { return n.spec.DNS }
func (n NetworkImpl) Default() []string { return n.spec.Default }

func (n *NetworkImpl) SetMAC(mac string)           { n.mac = mac }
func (n *NetworkImpl) SetDNS(nameservers []string) { n.spec.DNS = nameservers }
func (n *NetworkImpl) SetPreconfigured()           { n.preconfigured = true }

func (n NetworkImpl) CloudProps() NetworkCloudProps { return n.spec.CloudProps }

func (n *NetworkImpl) _final() {}

func (ns *Networks) UnmarshalJSON(data []byte) error {
	var newNets map[string]networkSpec2

	err := json.Unmarshal(data, &newNets)
	if err != nil {
		return bosherr.WrapErrorf(err, "Unmarshaling networks")
	}

	nets := Networks{}

	for netName, net := range newNets {
		nets[netName] = &NetworkImpl{spec: net}
	}

	*ns = nets

	return nil
}

func (ns Networks) Default() Network {
	var n Network

	for _, n = range ns {
		if n.IsDefaultFor("gateway") {
			break
		}
	}

	if n == nil {
		return &NetworkImpl{}
	}

	return n // returns last network
}

func (ns Networks) BackfillDefaultDNS(nameservers []string) {
	for _, n := range ns {
		if n.IsDefaultFor("dns") {
			if len(n.DNS()) == 0 {
				n.SetDNS(nameservers)
			}
			break
		}
	}
}

func (n NetworkImpl) IsDefaultFor(what string) bool {
	for _, def := range n.Default() {
		if def == what {
			return true
		}
	}

	return false
}

func (n NetworkImpl) IsDynamic() bool {
	return n.Type() == "dynamic"
}

func (n NetworkImpl) IPWithSubnetMask() string {
	netmaskIP := gonet.ParseIP(n.Netmask())
	if v4 := netmaskIP.To4(); v4 != nil {
		netmaskIP = v4
	}
	ones, _ := gonet.IPMask(netmaskIP).Size()
	return fmt.Sprintf("%s/%d", n.IP(), ones)
}
