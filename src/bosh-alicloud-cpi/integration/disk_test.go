/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package integration

import (
	"bosh-alicloud-cpi/mock"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("integration:disk", func() {
	It("can run the disk lifecycle", func() {
		By("create vm")
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
					"instance_name": "bosh-test-cpi-integration",
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
							"security_group_ids": ["${SECURITY_GROUP_ID}"],
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
		}`).P("STEMCELL_ID", existingStemcell).
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
		diskCid, err := caller.Call("create_disk", 1024, "{}", instCid)
		Expect(err).NotTo(HaveOccurred())

		By("sleep for awhile")
		time.Sleep(time.Duration(90) * time.Second)

		By("resize disk")
		_, err = caller.Call("resize_disk", diskCid, 30720)
		Expect(err).NotTo(HaveOccurred())

		By("attach disk")
		_, err = caller.Call("attach_disk", instCid, diskCid)
		Expect(err).NotTo(HaveOccurred())

		By("set disk meta data")
		r = caller.Run(mock.NewBuilder(`{
        	"method": "set_disk_metadata",
        	"arguments": [
                "${DISK_ID}",
                {
					"director": "my-bosh",
					"deployment": "cf",
					"instance_id": "${INST_ID}",
					"job": "consul",
					"instance_index": "0",
					"instance_name": "consul/441e940e-2ffe-4208-993e-3e5f888e2b7e",
					"attached_at": "2017-11-03T06:21:27Z"
                }
	        ],
    	    "context": {
                "director_uuid": "d5555ed6-7688-4aae-9dff-4c4507042f3d",
                "request_id": "cpi-201248"
        	}
		}`).P("INST_ID", instCid).P("DISK_ID", diskCid).ToBytes())
		Expect(r.GetError()).NotTo(HaveOccurred())

		By("snapshot disk")
		ssid, err := caller.Call("snapshot_disk", diskCid, nil)
		Expect(err).NotTo(HaveOccurred())

		By("detach disk")
		_, err = caller.Call("detach_disk", instCid, diskCid)
		Expect(err).NotTo(HaveOccurred())

		By("delete disk")
		_, err = caller.Call("delete_disk", diskCid)
		Expect(err).NotTo(HaveOccurred())

		By("delete snapshot")
		_, err = caller.Call("delete_snapshot", ssid)
		Expect(err).NotTo(HaveOccurred())

		By("delete vm")
		_, err = caller.Call("delete_vm", instCid)
		Expect(err).NotTo(HaveOccurred())
	})
})
