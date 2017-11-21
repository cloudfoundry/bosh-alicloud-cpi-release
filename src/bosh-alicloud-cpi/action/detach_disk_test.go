/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package action

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/denverdino/aliyungo/ecs"
)

var _ = Describe("cpi:detach_disk", func() {
	It("can detach disk", func() {
		instCid, _ := mockContext.NewInstance()
		diskCid, disk := mockContext.NewDisk(instCid)
		disk.Status = ecs.DiskStatusInUse

		_, err := caller.Call("detach_disk", instCid, diskCid)
		Expect(err).NotTo(HaveOccurred())
		Expect(disk.InstanceId).Should(Equal(""))
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