/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package action

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("cpi:attach_disk", func() {
	It("can attach disk", func() {
		diskCid, _ := mockContext.NewDisk("")
		instCid, _ := mockContext.NewInstance()

		_, err := caller.Call("attach_disk", instCid, diskCid)
		Expect(err).NotTo(HaveOccurred())
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