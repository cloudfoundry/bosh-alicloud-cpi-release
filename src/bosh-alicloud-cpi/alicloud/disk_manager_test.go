/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package alicloud

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/denverdino/aliyungo/ecs"
)

var _ = Describe("DiskManager", func() {
	It("should amend disk path", func() {
		By("cloud_efficiency /dev/xvdc -> /dev/vdc")
		Expect(amendDiskPath("/dev/xvdc", ecs.DiskCategoryCloudEfficiency)).To(Equal("/dev/vdc"))
		By("cloud_efficiency /dev/vdc -> /dev/vdc")
		Expect(amendDiskPath("/dev/vdc", ecs.DiskCategoryCloudEfficiency)).To(Equal("/dev/vdc"))
		By("ssd /dev/xvdc -> /dev/xvdc")
		Expect(amendDiskPath("/dev/xvdc", ecs.DiskCategoryCloudSSD)).To(Equal("/dev/xvdc"))
		By("ssd /dev/vdc -> /dev/vdc")
		Expect(amendDiskPath("/dev/vdc", ecs.DiskCategoryCloudEfficiency)).To(Equal("/dev/vdc"))
	})
})