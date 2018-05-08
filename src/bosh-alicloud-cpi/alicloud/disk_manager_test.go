/*
 * Copyright (C) 2017-2018 Alibaba Group Holding Limited
 */
package alicloud

import (
	"bosh-alicloud-cpi/alicloud"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DiskManager", func() {
	It("should amend disk path", func() {
		By("cloud_efficiency /dev/xvdc -> /dev/vdc")
		Expect(AmendDiskPath("/dev/xvdc", alicloud.DiskCategoryCloudEfficiency)).To(Equal("/dev/vdc"))
		By("cloud_efficiency /dev/vdc -> /dev/vdc")
		Expect(AmendDiskPath("/dev/vdc", alicloud.DiskCategoryCloudEfficiency)).To(Equal("/dev/vdc"))
		By("ssd /dev/xvdc -> /dev/xvdc")
		Expect(AmendDiskPath("/dev/xvdc", alicloud.DiskCategoryCloudSSD)).To(Equal("/dev/xvdc"))
		By("ssd /dev/vdc -> /dev/vdc")
		Expect(AmendDiskPath("/dev/vdc", alicloud.DiskCategoryCloudEfficiency)).To(Equal("/dev/vdc"))
	})
})
