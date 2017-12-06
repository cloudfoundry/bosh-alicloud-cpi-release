/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package integration

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"bosh-alicloud-cpi/mock"
)

var _ = Describe("integration:stemcell", func() {
	It("can upload a stemcell(light)", func() {
		bytes := mock.NewBuilder(`{
			"method": "create_stemcell",
			"arguments": [
				"/var/vcap/data/tmp/director/stemcell20170926-9684-17ncxdz/image",
				{
					"architecture": "x86_64",
					"container_format": "",
					"disk": "50",
					"disk_format": "rawdisk",
					"hypervisor": "kvm",
					"image_id": {
						"cn-beijing": "${STEMCELL_ID}",
						"cn-qingdao": "m-m5e1bdeo0a0ujw2s0ovv",
						"cn-zhangjiakou": "m-8vbdd1mi6hwlddicqt2x"
					},
					"infrastructure": "alicloud",
					"name": "bosh-alicloud-kvm-hubuntu-trusty-go_agent",
					"os_distro": "ubuntu",
					"os_type": "linux",
					"root_device_name": "/dev/vda1",
					"source_url": "",
					"version": 1003
				}
			],
			"context": {
				"director_uuid": "073eac6e-7a35-4a49-8c42-68988ea16ca7"
			}
		}`).P("STEMCELL_ID", stemcellId).ToBytes()

		r := caller.Run(bytes)
		Expect(r.GetError()).NotTo(HaveOccurred())
		Expect(r.GetResultString()).Should(Equal(stemcellId))
	})

	It("can create a stemcell from oss", func() {
		bytes := mock.NewBuilder(`{
			"method": "create_stemcell",
			"arguments": [
				"/var/vcap/data/tmp/director/stemcell20170926-9684-17ncxdz/image",
				{
					"architecture": "x86_64",
					"container_format": "",
					"disk": 40960,
					"disk_format": "RAW",
					"hypervisor": "kvm",
					"infrastructure": "alicloud",
					"name": "bosh-alicloud-kvm-ubuntu-trusty-go_agent",
					"oss_bucket": "bosh-stemcell",
					"oss_object": "alicloud-kvm-ubuntu.raw",
					"os_type": "linux",
					"root_device_name": "/dev/vda1",
					"version": 1003
				}
			],
			"context": {
				"director_uuid": "073eac6e-7a35-4a49-8c42-68988ea16ca7"
			}
		}`).P("STEMCELL_ID", stemcellId).ToBytes()

		r := caller.Run(bytes)
		Expect(r.GetError()).NotTo(HaveOccurred())
	})

	// this case can run pass, but on CI cant not fetch tarball, so annotate it temporary
	//FIt("can create a stemcell from tarball", func() {
	//	bytes := mock.NewBuilder(`{
	//		"method": "create_stemcell",
	//		"arguments": [
	//			"/.bosh/bosh-stemcell-1008-alicloud-kvm-ubuntu-trusty-go_agent.tgz",
	//			{
	//				"architecture": "x86_64",
	//				"container_format": "",
	//				"disk": 40960,
	//				"disk_format": "RAW",
	//				"hypervisor": "kvm",
	//				"infrastructure": "alicloud",
	//				"name": "bosh-alicloud-kvm-ubuntu",
	//				"os_type": "linux",
	//				"os_distro": "Ubuntu",
	//				"root_device_name": "/dev/vda1",
	//				"version": "1008"
	//			}
	//		],
	//		"context": {
	//			"director_uuid": "073eac6e-7a35-4a49-8c42-68988ea16ca7"
	//		}
	//	}`).P("STEMCELL_ID", stemcellId).ToBytes()
	//
	//	r := caller.Run(bytes)
	//	Expect(r.GetError()).NotTo(HaveOccurred())
	//})

	It("can delete a stemcell(light)", func() {
		_, err := caller.Call("delete_stemcell", stemcellId)
		Expect(err).NotTo(HaveOccurred())
	})
})
