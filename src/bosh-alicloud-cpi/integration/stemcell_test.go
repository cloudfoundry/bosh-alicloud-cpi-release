/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package integration

import (
	"bosh-alicloud-cpi/mock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("integration:stemcell", func() {
	var stemcellCID string
	It("executes the stemcell lifecycle with an oss bucket", func() {
		bytes := mock.NewBuilder(`{
			"method": "create_stemcell",
			"arguments": [
				"",
				{
					"architecture": "x86_64",
					"container_format": "bare",
					"disk": 3072,
					"disk_format": "raw",
					"hypervisor": "kvm",
					"infrastructure": "alicloud",
					"name": "bosh-alicloud-kvm-ubuntu-xenial-go_agent",
					"os_type": "linux",
					"os_distro": "ubuntu",
					"root_device_name": "/dev/vda1",
					"version": "${STEMCELL_VERSION}",
					"oss_bucket": "${OSS_BUCKET}",
					"oss_object": "${OSS_OBJECT}"
				}
			]
		}`).P("STEMCELL_VERSION", stemcellVersion).
			P("OSS_BUCKET", imageOssBucket).
			P("OSS_OBJECT", imageOssObject).
			ToBytes()

		r := caller.Run(bytes)
		Expect(r.GetError()).NotTo(HaveOccurred())
		stemcellCID = r.GetResultString()
	})

	It("can delete a stemcell(light)", func() {
		_, err := caller.Call("delete_stemcell", stemcellCID)
		Expect(err).NotTo(HaveOccurred())
	})
})
