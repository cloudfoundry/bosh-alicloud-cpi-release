/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package action

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"bosh-alicloud-cpi/mock"
)

var _ = Describe("cpi:attach_disk", func() {
	It("can attach disk", func() {
		diskCid, _ := mockContext.NewDisk()
		instCid, _ := mockContext.NewInstance()

		in := mock.NewBuilder(`{
			"method": "attach_disk",
				"arguments": [
					"${INST_CID}",
					"${DISK_CID}"
				],
				"context": {
				"director_uuid": "911133bb-7d44-4811-bf8a-b215608bf084"
			}
		}`).P("DISK_CID", diskCid).P("INST_CID", instCid).ToBytes()
		r := caller.Run(in)
		Expect(r.GetError()).NotTo(HaveOccurred())
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