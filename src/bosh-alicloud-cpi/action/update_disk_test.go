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
	Context("in-place mutations (ModifyDiskSpec path)", func() {
		It("changes disk category when the target category differs", func() {
			cid, disk := mockContext.NewDisk("")
			Expect(disk.Category).To(Equal(string(alicloud.DiskCategoryCloudEfficiency)))

			r, err := caller.CallGenericAPIVersion("update_disk", 2, cid, disk.Size*1024, map[string]interface{}{
				"category": string(alicloud.DiskCategoryCloudESSD),
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(r).To(Equal(cid))

			updated, ok := mockContext.Disks[cid]
			Expect(ok).To(BeTrue())
			Expect(updated.Category).To(Equal(string(alicloud.DiskCategoryCloudESSD)))
		})

		It("is a no-op when category already matches and size is unchanged", func() {
			cid, disk := mockContext.NewDisk("")
			disk.Category = string(alicloud.DiskCategoryCloudESSD)

			r, err := caller.CallGenericAPIVersion("update_disk", 2, cid, disk.Size*1024, map[string]interface{}{
				"category": string(alicloud.DiskCategoryCloudESSD),
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(r).To(Equal(cid))

			updated := mockContext.Disks[cid]
			Expect(updated.Category).To(Equal(string(alicloud.DiskCategoryCloudESSD)))
		})

		It("resizes when category matches but the requested size is larger", func() {
			cid, disk := mockContext.NewDisk("")
			disk.Category = string(alicloud.DiskCategoryCloudESSD)
			initialSize := disk.Size

			_, err := caller.CallGenericAPIVersion("update_disk", 2, cid, (initialSize+10)*1024, map[string]interface{}{
				"category": string(alicloud.DiskCategoryCloudESSD),
			})
			Expect(err).NotTo(HaveOccurred())

			updated := mockContext.Disks[cid]
			Expect(updated.Size).To(BeNumerically(">", initialSize))
		})

		It("applies category change AND resize in a single call", func() {
			cid, disk := mockContext.NewDisk("")
			Expect(disk.Category).To(Equal(string(alicloud.DiskCategoryCloudEfficiency)))
			initialSize := disk.Size

			_, err := caller.CallGenericAPIVersion("update_disk", 2, cid, (initialSize+10)*1024, map[string]interface{}{
				"category": string(alicloud.DiskCategoryCloudESSD),
			})
			Expect(err).NotTo(HaveOccurred())

			updated := mockContext.Disks[cid]
			Expect(updated.Category).To(Equal(string(alicloud.DiskCategoryCloudESSD)))
			Expect(updated.Size).To(BeNumerically(">", initialSize))
		})
	})

	Context("performance_level updates (ESSD)", func() {
		It("applies a performance_level change on an ESSD disk", func() {
			cid, disk := mockContext.NewDisk("")
			disk.Category = string(alicloud.DiskCategoryCloudESSD)
			disk.PerformanceLevel = "PL1"

			_, err := caller.CallGenericAPIVersion("update_disk", 2, cid, disk.Size*1024, map[string]interface{}{
				"category":          string(alicloud.DiskCategoryCloudESSD),
				"performance_level": "PL2",
			})
			Expect(err).NotTo(HaveOccurred())

			updated := mockContext.Disks[cid]
			Expect(updated.PerformanceLevel).To(Equal("PL2"))
			Expect(updated.Category).To(Equal(string(alicloud.DiskCategoryCloudESSD)))
		})

		It("applies category change AND performance_level in a single call", func() {
			cid, disk := mockContext.NewDisk("")
			Expect(disk.Category).To(Equal(string(alicloud.DiskCategoryCloudEfficiency)))

			_, err := caller.CallGenericAPIVersion("update_disk", 2, cid, disk.Size*1024, map[string]interface{}{
				"category":          string(alicloud.DiskCategoryCloudESSD),
				"performance_level": "PL2",
			})
			Expect(err).NotTo(HaveOccurred())

			updated := mockContext.Disks[cid]
			Expect(updated.Category).To(Equal(string(alicloud.DiskCategoryCloudESSD)))
			Expect(updated.PerformanceLevel).To(Equal("PL2"))
		})

		It("leaves performance_level unchanged when the target matches the current level", func() {
			cid, disk := mockContext.NewDisk("")
			disk.Category = string(alicloud.DiskCategoryCloudESSD)
			disk.PerformanceLevel = "PL1"

			_, err := caller.CallGenericAPIVersion("update_disk", 2, cid, disk.Size*1024, map[string]interface{}{
				"category":          string(alicloud.DiskCategoryCloudESSD),
				"performance_level": "PL1",
			})
			Expect(err).NotTo(HaveOccurred())

			updated := mockContext.Disks[cid]
			Expect(updated.PerformanceLevel).To(Equal("PL1"))
		})
	})

	Context("unsupported transitions", func() {
		It("returns Bosh::Clouds::NotSupported when AliCloud refuses the in-place category change", func() {
			// cloud_essd -> cloud_efficiency is a downgrade; the mock's
			// ModifyDiskCategory simulates AliCloud returning
			// InvalidDiskCategory.NotSupported, which the CPI translates to
			// Bosh::Clouds::NotSupported so the director can fall back.
			cid, disk := mockContext.NewDisk("")
			disk.Category = string(alicloud.DiskCategoryCloudESSD)

			_, err := caller.CallGenericAPIVersion("update_disk", 2, cid, disk.Size*1024, map[string]interface{}{
				"category": string(alicloud.DiskCategoryCloudEfficiency),
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(`"type":"Bosh::Clouds::NotSupported"`))

			// Disk must be untouched — director drives the recovery.
			Expect(mockContext.Disks[cid].Category).To(Equal(string(alicloud.DiskCategoryCloudESSD)))
		})
	})

	Context("failure modes", func() {
		It("returns Bosh::Clouds::NotSupported for a shrink so the director falls back to copy migration", func() {
			// AliCloud can't shrink in-place. The CPI returns NotSupported (not a
			// generic error) so the director's copy/migrate path handles the shrink,
			// matching the behavior from before update_disk existed.
			cid, disk := mockContext.NewDisk("")
			disk.Category = string(alicloud.DiskCategoryCloudESSD)
			initialSize := disk.Size

			_, err := caller.CallGenericAPIVersion("update_disk", 2, cid, (initialSize-5)*1024, map[string]interface{}{
				"category": string(alicloud.DiskCategoryCloudESSD),
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring(`"type":"Bosh::Clouds::NotSupported"`))
			Expect(err.Error()).To(ContainSubstring("cannot shrink"))

			updated := mockContext.Disks[cid]
			Expect(updated.Size).To(Equal(initialSize))
		})

		It("errors when the disk never returns to Available after the category change", func() {
			cid, disk := mockContext.NewDisk("")
			disk.Status = string(alicloud.DiskStatusCreating)

			_, err := caller.CallGenericAPIVersion("update_disk", 2, cid, disk.Size*1024, map[string]interface{}{
				"category": string(alicloud.DiskCategoryCloudESSD),
			})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("WaitForDiskStatus"))
		})

		It("returns an error when the disk does not exist", func() {
			_, err := caller.CallGenericAPIVersion("update_disk", 2, "non-existent-disk", 20480, map[string]interface{}{
				"category": string(alicloud.DiskCategoryCloudESSD),
			})
			Expect(err).To(HaveOccurred())
		})
	})
})
