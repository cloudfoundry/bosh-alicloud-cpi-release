/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package integration

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"bosh-alicloud-cpi/mock"
	"time"
)

var _ = Describe("integration:vm", func() {
	It("can run the vm lifecycle", func() {
		By("create vm")
		in := mock.NewBuilder(`{
			"method": "create_vm",
			"arguments": [
				"be387a69-c5d5-4b94-86c2-978581354b50",
				"${STEMCELL_ID}", {
					"ephemeral_disk": {
						"size": "40_960",
						"category": "cloud_efficiency"
					},
					"instance_name": "test-cc",
					"instance_type": "ecs.n4.small",
					"system_disk": {
						"size": "61_440",
						"category": "cloud_efficiency"
					}
				}, {
					"private": {
						"type": "manual",
						"ip": "${INTERNAL_IP}",
						"netmask": "${INTERNAL_NETMASK}",
						"cloud_properties": {
							"security_group_id": "${SECURITY_GROUP_ID}",
							"vswitch_id": "${VSWITCH_ID}"
						},
						"default": [
							"dns",
							"gateway"
						],
						"dns": [
							"8.8.8.8"
						],
						"gateway": "${INTERNAL_GW}"
					}
				},
				[],
				{}
			],
			"context": {
				"director_uuid": "911133bb-7d44-4811-bf8a-b215608bf084"
			}
		}`).
			P("STEMCELL_ID", stemcellId).
			P("SECURITY_GROUP_ID", securityGroupId).
			P("VSWITCH_ID", vswitchId).
			P("INTERNAL_IP", internalIp).
			P("INTERNAL_NETMASK", internalNetmask).
			P("INTERNAL_GW", internalGw).
			ToBytes()

		r := caller.Run(in)
		Expect(r.GetError()).NotTo(HaveOccurred())
		cid := r.GetResultString()

		By("sleep for awhile")
		time.Sleep(time.Duration(90) * time.Second)

		By("delete vm")
		_, err := caller.Call("delete_vm", cid)
		Expect(err).NotTo(HaveOccurred())

		By("vm should not exists")
		exists, err := caller.CallGeneric("has_vm", cid)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(BeFalse())
	})

	It("can run the vm lifecycle with slb", func() {
		if slbId == "" {
			Skip("no CPI_SLB_IP settings")
		}
		By("create vm with slb")
		in := mock.NewBuilder(`{
			"method": "create_vm",
			"arguments": [
				"be387a69-c5d5-4b94-86c2-978581354b50",
				"${STEMCELL_ID}",
				{
					"ephemeral_disk": {
						"size": "40_960",
						"category": "cloud_efficiency"
					},
					"instance_name": "test-cc",
					"instance_type": "ecs.n4.small",
					"slbs": ["${SLB_ID}"],
					"system_disk": {
						"size": "61_440",
						"category": "cloud_efficiency"
					}
				},
				{
					"private": {
						"type": "manual",
						"ip": "${INTERNAL_IP}",
						"netmask": "${INTERNAL_NETMASK}",
						"cloud_properties": {
							"security_group_id": "${SECURITY_GROUP_ID}",
							"vswitch_id": "${VSWITCH_ID}"
						},
						"default": [
							"dns",
							"gateway"
						],
						"dns": [
							"8.8.8.8"
						],
						"gateway": "${INTERNAL_GW}"
					}
				},
				[],
				{}
			],
			"context": {
				"director_uuid": "911133bb-7d44-4811-bf8a-b215608bf084"
			}
		}`).P("STEMCELL_ID", stemcellId).
			P("SECURITY_GROUP_ID", securityGroupId).
			P("VSWITCH_ID", vswitchId).
			P("INTERNAL_IP", internalIp).
			P("INTERNAL_NETMASK", internalNetmask).
			P("INTERNAL_GW", internalGw).
			P("SLB_ID", slbId).
			ToBytes()

		r := caller.Run(in)
		Expect(r.GetError()).NotTo(HaveOccurred())
		cid := r.GetResultString()

		By("sleep 90 seconds to make sure vm is fully started")
		time.Sleep(time.Duration(90) * time.Second)

		By("delete vm")
		_, err := caller.Call("delete_vm", cid)
		Expect(err).NotTo(HaveOccurred())

		By("vm should not exists")
		exists, err := caller.CallGeneric("has_vm", cid)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(BeFalse())
	})

	It("can run the vm lifecycle with external ip", func() {
		if externalIp == "" {
			Skip("no CPI_EXTERNAL_IP settings")
		}
		By("create vm with external ip")
		in := mock.NewBuilder(`{
			"method": "create_vm",
			"arguments": [
				"be387a69-c5d5-4b94-86c2-978581354b50",
				"${STEMCELL_ID}",
				{
					"ephemeral_disk": {
						"size": "40_960",
						"category": "cloud_efficiency"
					},
					"instance_name": "test-cc",
					"instance_type": "ecs.n4.small",
					"system_disk": {
						"size": "61_440",
						"category": "cloud_efficiency"
					}
				},
				{
					"private": {
						"type": "manual",
						"ip": "${INTERNAL_IP}",
						"netmask": "${INTERNAL_NETMASK}",
						"cloud_properties": {
							"security_group_id": "${SECURITY_GROUP_ID}",
							"vswitch_id": "${VSWITCH_ID}"
						},
						"default": [
							"dns",
							"gateway"
						],
						"dns": [
							"8.8.8.8"
						],
						"gateway": "${INTERNAL_GW}"
					},
					"public": {
						"type": "vip",
						"ip": "${EXTERNAL_IP}",
						"type": "vip"
					}
				},
				[],
				{}
			],
			"context": {
				"director_uuid": "911133bb-7d44-4811-bf8a-b215608bf084"
			}
		}`).
			P("STEMCELL_ID", stemcellId).
			P("SECURITY_GROUP_ID", securityGroupId).
			P("VSWITCH_ID", vswitchId).
			P("INTERNAL_IP", internalIp).
			P("INTERNAL_NETMASK", internalNetmask).
			P("INTERNAL_GW", internalGw).
			P("EXTERNAL_IP", externalIp).
			ToBytes()

		r := caller.Run(in)
		Expect(r.GetError()).NotTo(HaveOccurred())
		cid := r.GetResultString()

		By("set vm metadata")
		in = mock.NewBuilder(`{
			"method": "set_vm_metadata",
			"arguments": [
				"${INSTANCE_ID}", {
					"director": "my-bosh",
					"deployment": "ali-cf-215",
					"id": "4daa76e3-db6d-4550-b9b4-c504e7865f62",
					"job": "dea_ng",
					"index": "0",
					"name": "dea_ng/4daa76e3-db6d-4550-b9b4-c504e7865f62",
					"created_at": "2017-10-25T09:00:07Z"
				}
			],
			"context": {
					"director_uuid": "580da067-b2ff-4eb6-b271-23cc76409121",
					"request_id": "cpi-440280"
			}
		}`).P("${INSTANCE_ID}", cid).ToBytes()

		r = caller.Run(in)
		Expect(r.GetError()).NotTo(HaveOccurred())

		By("sleep 90 seconds to make sure vm is fully started")
		time.Sleep(time.Duration(90) * time.Second)

		By("delete vm")
		_, err := caller.Call("delete_vm", cid)
		Expect(err).NotTo(HaveOccurred())

		By("vm should not exists")
		exists, err := caller.CallGeneric("has_vm", cid)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(BeFalse())
	})

	It("can run the spot instance lifecycle", func() {
		By("create spot instance")
		in := mock.NewBuilder(`{
			"method": "create_vm",
			"arguments": [
				"be387a69-c5d5-4b94-86c2-978581354b50",
				"${STEMCELL_ID}", {
					"ephemeral_disk": {
						"size": "40_960",
						"category": "cloud_efficiency"
					},
					"instance_name": "test-cc",
					"instance_type": "ecs.n4.small",
					"spot_strategy": "${SPOT_STRATEGY}",
					"spot_price_limit": ${SPOT_PRICE_LIMIT},
					"system_disk": {
						"size": "61_440",
						"category": "cloud_efficiency"
					}
				}, {
					"private": {
						"type": "manual",
						"ip": "${INTERNAL_IP}",
						"netmask": "${INTERNAL_NETMASK}",
						"cloud_properties": {
							"security_group_id": "${SECURITY_GROUP_ID}",
							"vswitch_id": "${VSWITCH_ID}"
						},
						"default": [
							"dns",
							"gateway"
						],
						"dns": [
							"8.8.8.8"
						],
						"gateway": "${INTERNAL_GW}"
					}
				},
				[],
				{}
			],
			"context": {
				"director_uuid": "911133bb-7d44-4811-bf8a-b215608bf084"
			}
		}`).
			P("STEMCELL_ID", stemcellId).
			P("SECURITY_GROUP_ID", securityGroupId).
			P("VSWITCH_ID", vswitchId).
			P("INTERNAL_IP", internalIp).
			P("INTERNAL_NETMASK", internalNetmask).
			P("INTERNAL_GW", internalGw).
			P("SPOT_STRATEGY", spotStrategy).
			P("SPOT_PRICE_LIMIT", spotPriceLimit).
			ToBytes()

		r := caller.Run(in)
		Expect(r.GetError()).NotTo(HaveOccurred())
		cid := r.GetResultString()

		By("sleep for awhile")
		time.Sleep(time.Duration(90) * time.Second)

		By("delete vm")
		_, err := caller.Call("delete_vm", cid)
		Expect(err).NotTo(HaveOccurred())

		By("vm should not exists")
		exists, err := caller.CallGeneric("has_vm", cid)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(BeFalse())
	})

	//It("can create vm with manual ip, and delete it", func() {})
	//It("can create vm with dynamic ip, and delete it", func() {})
	//It("can create vm with external ip, and delete it", func() {})
	//It("can create vm with key pair, and delete it", func() {})
	//It("can create vm, then start, stop and delete it", func() {})
	//It("can create vm, then start, reboot, stop and delete it", func() {})
	//It("can create vm, then start, reboot, stop and delete it", func() {})
})
