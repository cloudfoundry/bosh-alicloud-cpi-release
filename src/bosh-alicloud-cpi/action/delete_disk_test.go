/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package action

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/denverdino/aliyungo/ecs"
)


var _ = Describe("cpi:delete_disk", func() {
	It("can delete disk", func() {
		cid, _ := mockContext.NewDisk("")

		_, err := caller.Call("delete_disk", cid)
		Expect(err).NotTo(HaveOccurred())

		exist,err := caller.CallGeneric("has_disk", cid)
		Expect(err).NotTo(HaveOccurred())
		Expect(exist).To(BeFalse())
	})

	It("delete disk will failed with running instance", func() {
		instCid, _ := mockContext.NewInstance()
		diskCid, disk := mockContext.NewDisk(instCid)
		disk.Status = ecs.DiskStatusInUse

		_, err := caller.Call("delete_disk", diskCid)
		Expect(err).Should(HaveOccurred())
		Expect(err.Error()).Should(ContainSubstring("can't delete disk"))
	})
})