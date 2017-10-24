/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package action

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"encoding/json"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
)

type NetworksInputEntity struct {
	networks apiv1.Networks
}

func parseNetworks(in string) (Networks, error) {
	var networks apiv1.Networks

	err := json.Unmarshal([]byte(in), &networks)
	if err != nil {
		return Networks{}, err
	}
	return NewNetworks(networks)
}

var _ = Describe("Networks", func() {
	It("can work with private network", func() {
		networks, err := parseNetworks(`{
			"private": {
				"type": "manual",
				"ip": "10.0.16.109",
				"netmask": "255.255.240.0",
				"cloud_properties": {
					"security_group_id": "sg-2ze2ct08gslmnwyv8c1k",
					"vswitch_id": "vpc-2ze3owai4kbkv2yf6nivg"
				},
				"default": [
					"dns",
					"gateway"
				],
				"dns": [
					"10.0.16.2"
				],
				"gateway": "10.0.16.1"
			}
		}`)
		Expect(err).NotTo(HaveOccurred())
		_ = networks
	})

	It("can work with default network", func() {
		networks, err := parseNetworks(`{
			"default": {
				"type": "manual",
				"ip": "10.0.16.109",
				"netmask": "255.255.240.0",
				"cloud_properties": {
					"security_group_id": "sg-2ze2ct08gslmnwyv8c1k",
					"vswitch_id": "vpc-2ze3owai4kbkv2yf6nivg"
				},
				"default": [
					"dns",
					"gateway"
				],
				"dns": [
					"10.0.16.2"
				],
				"gateway": "10.0.16.1"
			}
		}`)
		Expect(err).NotTo(HaveOccurred())
		_ = networks
	})


	It("can process networks with external ip", func() {
		networks, err := parseNetworks(`{
			"default": {
				"type": "manual",
				"ip": "10.0.16.109",
				"netmask": "255.255.240.0",
				"cloud_properties": {
					"security_group_id": "sg-2ze2ct08gslmnwyv8c1k",
					"vswitch_id": "vpc-2ze3owai4kbkv2yf6nivg"
				},
				"default": [
					"dns",
					"gateway"
				],
				"dns": [
					"10.0.16.2"
				],
				"gateway": "10.0.16.1"
			},
			"public": {
				"type": "vip",
				"ip": "47.100.100.100"
			}
		}`)
		Expect(err).NotTo(HaveOccurred())
		_ = networks
	})


	It("can process network with other name", func() {
		networks, err := parseNetworks(`{
			"cf1": {
				"type": "manual",
				"ip": "10.0.16.109",
				"netmask": "255.255.240.0",
				"cloud_properties": {
					"security_group_id": "sg-2ze2ct08gslmnwyv8c1k",
					"vswitch_id": "vpc-2ze3owai4kbkv2yf6nivg"
				},
				"default": [
					"dns",
					"gateway"
				],
				"dns": [
					"10.0.16.2"
				],
				"gateway": "10.0.16.1"
			}
		}`)
		Expect(err).NotTo(HaveOccurred())
		_ = networks
	})

})