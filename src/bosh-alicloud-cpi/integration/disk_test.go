/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package integration

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"bosh-alicloud-cpi/mock"
)

var _ = Describe("integration:disk", func() {
	It("can run the disk lifecycle", func() {
		By("create vm")
		in := mock.NewBuilder(`{
			"method": "create_vm",
			"arguments": [
				"be387a69-c5d5-4b94-86c2-978581354b50",
				"m-2zehhdtfg22hq46reabf",
				{
					"ephemeral_disk": {
						"size": "40_960",
						"category": "cloud_efficiency"
					},
					"image_id": "${STEMCELL_ID}",
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
			ToBytes()

		r := caller.Run(in)
		Expect(r.GetError()).NotTo(HaveOccurred())
		instCid := r.GetResultString()

		By("create disk")
		diskCid, err := caller.Call("create_disk", 30720, "{}", instCid)
		Expect(err).NotTo(HaveOccurred())

		By("attach disk")
		_, err = caller.Call("attach_disk", instCid, diskCid)
		Expect(err).NotTo(HaveOccurred())

		By("detach disk")
		_, err = caller.Call("detach_disk", instCid, diskCid)
		Expect(err).NotTo(HaveOccurred())

		By("delete disk")
		_, err = caller.Call("delete_disk", diskCid)
		Expect(err).NotTo(HaveOccurred())

		By("delete vm")
		_, err = caller.Call("delete_vm", instCid)
		Expect(err).NotTo(HaveOccurred())
	})
})
