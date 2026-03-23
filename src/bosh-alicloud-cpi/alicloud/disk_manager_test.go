/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package alicloud

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("DiskManager", func() {
	It("should amend disk path", func() {
		By("cloud_efficiency /dev/xvdc -> /dev/vdc")
		Expect(AmendDiskPath("/dev/xvdc", DiskCategoryCloudEfficiency)).To(Equal("/dev/vdc"))
		By("cloud_efficiency /dev/vdc -> /dev/vdc")
		Expect(AmendDiskPath("/dev/vdc", DiskCategoryCloudEfficiency)).To(Equal("/dev/vdc"))
		By("ssd /dev/xvdc -> /dev/xvdc")
		Expect(AmendDiskPath("/dev/xvdc", DiskCategoryCloudSSD)).To(Equal("/dev/vdc"))
		By("ssd /dev/vdc -> /dev/vdc")
		Expect(AmendDiskPath("/dev/vdc", DiskCategoryCloudEfficiency)).To(Equal("/dev/vdc"))
	})
})
