/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package action

import (
	"encoding/json"

	"bosh-alicloud-cpi/alicloud"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func parseCloudProps(in string) (Disks, error) {
	var props InstanceProps

	err := json.Unmarshal([]byte(in), &props)
	if err != nil {
		return Disks{}, err
	}
	return NewDisksWithProps(props.SystemDisk, props.EphemeralDisk)
}

var _ = Describe("Disks", func() {
	It("can work with system_disk & ephemeral_disk", func() {
		disks, err := parseCloudProps(`{
			"ephemeral_disk": {
				"size": "40_960",
				"category": "ephemeral_ssd"
			},
			"image_id": "m-2ze200tcuotb5uk2kol4",
			"instance_name": "bosh-test-cpi-integration",
			"instance_type": "ecs.n4.small",
			"system_disk": {
				"size": 51200,
				"type": "cloud_efficiency"
			}
		}`)
		Expect(err).NotTo(HaveOccurred())
		Expect(disks.SystemDisk.sizeGB).To(Equal(50))
		Expect(disks.SystemDisk.GetCategory()).To(Equal(alicloud.DiskCategoryCloudEfficiency))
		Expect(disks.EphemeralDisk.sizeGB).To(Equal(40))
		Expect(disks.EphemeralDisk.GetCategory()).To(Equal(alicloud.DiskCategoryEphemeralSSD))
	})
	It("can work with alternative size format and emit types", func() {
		disks, err := parseCloudProps(`{
			"ephemeral_disk": {
				"size": "81920"
			},
			"image_id": "m-2ze200tcuotb5uk2kol4",
			"instance_name": "bosh-test-cpi-integration",
			"instance_type": "ecs.n4.small",
			"system_disk": {
				"size": "102_400"
			}
		}`)

		Expect(err).NotTo(HaveOccurred())
		Expect(disks.SystemDisk.sizeGB).To(Equal(100))
		Expect(disks.SystemDisk.GetCategory()).To(Equal(DefaultDiskCategory))
		Expect(disks.EphemeralDisk.sizeGB).To(Equal(80))
		Expect(disks.EphemeralDisk.GetCategory()).To(Equal(DefaultDiskCategory))
	})
	It("can work without system_disk", func() {
		disks, err := parseCloudProps(`{
			"ephemeral_disk": {
				"size": 40960,
				"category": "cloud"
			}
		}`)
		Expect(err).NotTo(HaveOccurred())
		Expect(disks.SystemDisk.sizeGB).To(Equal(DefaultSystemDiskSizeGB))
		Expect(disks.SystemDisk.ecsCategory).To(Equal(alicloud.DiskCategoryCloudEfficiency))
		Expect(disks.EphemeralDisk.sizeGB).To(Equal(40))
		Expect(disks.EphemeralDisk.GetCategory()).To(Equal(alicloud.DiskCategoryCloud))
	})
	It("can work without ephemeral_disk", func() {
		disks, err := parseCloudProps(`{
			"system_disk": {
				"size": 81920,
				"category": "cloud_ssd"
			}
		}
		`)
		Expect(err).NotTo(HaveOccurred())
		Expect(disks.SystemDisk.sizeGB).To(Equal(80))
		Expect(disks.SystemDisk.ecsCategory).To(Equal(alicloud.DiskCategoryCloudSSD))
	})
	It("will check system disk category", func() {
		_, err := parseCloudProps(`{
			"system_disk": {
				"size": 81920,
				"category": "cloud"
			}
		}`)
		Expect(err).Should(HaveOccurred())
		//Expect(disks.SystemDisk.sizeGB).To(Equal(40))
		//Expect(disks.SystemDisk.sizeGB).To(Equal(60))
	})
})
