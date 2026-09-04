/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package action

import (
	"bosh-alicloud-cpi/alicloud"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("cpi:detach_disk", func() {
	It("can detach disk", func() {
		instCid, _ := mockContext.NewInstance()
		diskCid, disk := mockContext.NewDisk(instCid)
		disk.Status = string(alicloud.DiskStatusInUse)

		_, err := caller.Call("detach_disk", instCid, diskCid)
		Expect(err).NotTo(HaveOccurred())
		Expect(disk.InstanceId).Should(Equal(""))
	})

	It("fails immediately when disk is Modifying", func() {
		// "Modifying" is undocumented in the AliCloud DescribeDisks enum but
		// observed in production (Task 26321593). This test pins the string
		// literal via DiskStatusModifying so a casing/spelling change is caught.
		instCid, _ := mockContext.NewInstance()
		diskCid, disk := mockContext.NewDisk(instCid)
		disk.Status = string(alicloud.DiskStatusModifying)

		_, err := caller.Call("detach_disk", instCid, diskCid)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(ContainSubstring("Modifying"))
	})

	//It("can attach disk with right registry", func() {
	//	By("attach disk")
	//	diskCid, _ := mockContext.NewDisk()
	//	instCid, _ := mockContext.NewInstance()
	//
	//	in := mock.NewBuilder(`
	//
	//	`).
	//	P("DISK_CID", diskCid).
	//	P("INST_CID", instCid).
	//	ToBytes()
	//
	//	r :=  caller.Run(in)
	//
	//	Expect(r.Error).NotTo(HaveOccurred())
	//
	//	//
	//	// TODO: use mock method to detect execution results
	//	// disks := caller.Disks.GetDisk()
	//	// Expect(disks.GetDiskStatus(id)).To(Equal())
	//
	//	By("update registry right")
	//
	//})
})
