/*
 * Copyright (C) 2017-2018 Alibaba Group Holding Limited
 */
package action

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("cpi:resize_disk", func() {
	It("can run", func() {
		cid, _ := mockContext.NewDisk("")
		_, err := caller.Call("resize_disk", cid, 51200)
		Expect(err).NotTo(HaveOccurred())

		disk, ok := mockContext.Disks[cid]
		Expect(ok).Should(BeTrue())
		Expect(disk.Size).Should(Equal(50))
	})
})
