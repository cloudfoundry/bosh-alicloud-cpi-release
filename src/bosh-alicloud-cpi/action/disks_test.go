/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package action

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
	"encoding/json"
)

func parseCloudProps(in string) (Disks, error) {
	var props InstanceProps

	err := json.Unmarshal([]byte(in), &props)
	if err != nil {
		return Disks{}, err
	}
	return NewDisks(props.SystemDisk, []DiskInfo{props.EphemeralDisk})
}

var _ = Describe("Disks", func() {
	It("can work with cloud properties", func() {
		disks, err := parseCloudProps(`{
			"ephemeral_disk": {
				"size": 40_000,
				"type": "cloud_efficiency"
			},
			"image_id": "m-2ze200tcuotb5uk2kol4",
			"instance_name": "test-cc",
			"instance_type": "ecs.n4.small",
			"system_disk": {
				"size": 60_000,
				"type": "cloud_efficiency"
			}
		}`)
		Expect(err).NotTo(HaveOccurred())
		Expect(disks.SystemDisk.sizeGB).To(Equal(40))
		Expect(disks.SystemDisk.sizeGB).To(Equal(60))
	})
	It("can work with alternative size format", func() {
		disks, err := parseCloudProps(`{
			"ephemeral_disk": {
				"size": 40000,
				"type": "cloud_efficiency"
			},
			"image_id": "m-2ze200tcuotb5uk2kol4",
			"instance_name": "test-cc",
			"instance_type": "ecs.n4.small",
			"system_disk": {
				"size": 10_240,
				"type": "cloud_efficiency"
			}
		}`)

		Expect(err).NotTo(HaveOccurred())
		Expect(disks.SystemDisk.sizeGB).To(Equal(40))
		Expect(disks.SystemDisk.sizeGB).To(Equal(60))
	})
	It("can work without ephemeral_disk", func() {
		disks, err := parseCloudProps(`{
			"ephemeral_disk": {
				"size": 40000,
				"type": "cloud_efficiency"
			},
			"image_id": "m-2ze200tcuotb5uk2kol4",
			"instance_name": "test-cc",
			"instance_type": "ecs.n4.small",
		}
		`)
		Expect(err).NotTo(HaveOccurred())
		Expect(disks.SystemDisk.sizeGB).To(Equal(40))
		Expect(disks.SystemDisk.sizeGB).To(Equal(60))
	})
	It("can amend path for alicloud disk", func() {

		//Expect().NotTo(HaveOccurred())
		//Expect(disks.SystemDisk.sizeGB).To(Equal(40))
		//Expect(disks.SystemDisk.sizeGB).To(Equal(60))
	})
})


func TestDeviceTrim(t *testing.T) {
	device := "/dev/xvdc"
	if device[5] == 'x' {
		device = "/dev/" + string(device[6:])
	}

	t.Log(device)
}

