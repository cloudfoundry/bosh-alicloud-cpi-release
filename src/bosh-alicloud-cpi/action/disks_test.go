/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package action

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"encoding/json"
	"github.com/denverdino/aliyungo/ecs"
)

func parseCloudProps(in string) (Disks, error) {
	var props InstanceProps

	err := json.Unmarshal([]byte(in), &props)
	if err != nil {
		return Disks{}, err
	}
	return NewDisks(props.SystemDisk, props.EphemeralDisk)
}

var _ = Describe("Disks", func() {
	It("can work with system_disk & ephemeral_disk", func() {
		disks, err := parseCloudProps(`{
			"ephemeral_disk": {
				"size": "40_960",
				"type": "ephemeral_ssd"
			},
			"image_id": "m-2ze200tcuotb5uk2kol4",
			"instance_name": "test-cc",
			"instance_type": "ecs.n4.small",
			"system_disk": {
				"size": 51200,
				"type": "cloud_efficiency"
			}
		}`)
		Expect(err).NotTo(HaveOccurred())
		Expect(disks.SystemDisk.sizeGB).To(Equal(50))
		Expect(disks.SystemDisk.GetCategory()).To(Equal(ecs.DiskCategoryCloudEfficiency))
		Expect(disks.EphemeralDisk.sizeGB).To(Equal(40))
		Expect(disks.EphemeralDisk.GetCategory()).To(Equal(ecs.DiskCategoryEphemeralSSD))
	})
	It("can work with alternative size format and emit types", func() {
		disks, err := parseCloudProps(`{
			"ephemeral_disk": {
				"size": "81920"
			},
			"image_id": "m-2ze200tcuotb5uk2kol4",
			"instance_name": "test-cc",
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
				"type": "cloud"
			}
		}`)
		Expect(err).NotTo(HaveOccurred())
		Expect(disks.SystemDisk.sizeGB).To(Equal(DefaultSystemDiskSizeGB))
		Expect(disks.SystemDisk.ecsCategory).To(Equal(ecs.DiskCategoryCloudEfficiency))
		Expect(disks.EphemeralDisk.sizeGB).To(Equal(40))
		Expect(disks.EphemeralDisk.GetCategory()).To(Equal(ecs.DiskCategoryCloud))
	})
	It("can work without ephemeral_disk", func() {
		disks, err := parseCloudProps(`{
			"system_disk": {
				"size": 81920,
				"type": "cloud_ssd"
			}
		}
		`)
		Expect(err).NotTo(HaveOccurred())
		Expect(disks.SystemDisk.sizeGB).To(Equal(80))
		Expect(disks.SystemDisk.ecsCategory).To(Equal(ecs.DiskCategoryCloudSSD))
	})
	It("will check system disk category", func() {
		_, err := parseCloudProps(`{
			"system_disk": {
				"size": 81920,
				"type": "cloud"
			}
		}`)
		Expect(err).Should(HaveOccurred())
		//Expect(disks.SystemDisk.sizeGB).To(Equal(40))
		//Expect(disks.SystemDisk.sizeGB).To(Equal(60))
	})
})
