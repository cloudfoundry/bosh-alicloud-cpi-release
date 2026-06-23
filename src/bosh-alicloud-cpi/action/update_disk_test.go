/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package action

import (
	"bosh-alicloud-cpi/alicloud"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("cpi:update_disk", func() {
	It("changes disk category when target category differs", func() {
		cid, disk := mockContext.NewDisk("")
		Expect(disk.Category).To(Equal(string(alicloud.DiskCategoryCloudEfficiency)))

		r, err := caller.CallGeneric("update_disk", cid, 20480, map[string]interface{}{
			"category": string(alicloud.DiskCategoryCloudESSD),
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(r).To(Equal(cid))

		updated, ok := mockContext.Disks[cid]
		Expect(ok).To(BeTrue())
		Expect(updated.Category).To(Equal(string(alicloud.DiskCategoryCloudESSD)))
	})

	It("is a no-op when category already matches", func() {
		cid, disk := mockContext.NewDisk("")
		disk.Category = string(alicloud.DiskCategoryCloudESSD)

		r, err := caller.CallGeneric("update_disk", cid, 20480, map[string]interface{}{
			"category": string(alicloud.DiskCategoryCloudESSD),
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(r).To(Equal(cid))

		updated := mockContext.Disks[cid]
		Expect(updated.Category).To(Equal(string(alicloud.DiskCategoryCloudESSD)))
	})

	It("resizes disk when category matches but new size is larger", func() {
		cid, disk := mockContext.NewDisk("")
		disk.Category = string(alicloud.DiskCategoryCloudESSD)
		initialSize := disk.Size

		_, err := caller.CallGeneric("update_disk", cid, (initialSize+10)*1024, map[string]interface{}{
			"category": string(alicloud.DiskCategoryCloudESSD),
		})
		Expect(err).NotTo(HaveOccurred())

		updated := mockContext.Disks[cid]
		Expect(updated.Size).To(BeNumerically(">", initialSize))
	})

	It("returns error when disk does not exist", func() {
		_, err := caller.CallGeneric("update_disk", "non-existent-disk", 20480, map[string]interface{}{
			"category": string(alicloud.DiskCategoryCloudESSD),
		})
		Expect(err).To(HaveOccurred())
	})
})
