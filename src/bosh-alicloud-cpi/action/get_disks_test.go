/*
 * Copyright (C) 2017-2018 Alibaba Group Holding Limited
 */
package action

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("cpi:get_disks", func() {
	It("can get with existent disk", func() {
		//r :=  caller.Run([]byte(`
		//
		//
		//`), )

		instCid, _ := mockContext.NewInstance()
		cid1, _ := mockContext.NewDisk(instCid)
		cid2, _ := mockContext.NewDisk(instCid)

		r, err := caller.CallGeneric("get_disks", instCid)
		Expect(err).NotTo(HaveOccurred())

		disks := r.([]interface{})
		Expect(disks).Should(ConsistOf(cid1, cid2))
	})
})
