/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package action

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("cpi:delete_vm", func() {
	It("can delete existent vm", func() {
		cid, _ := mockContext.NewInstance()

		_, err :=  caller.Call("delete_vm", cid)
		Expect(err).NotTo(HaveOccurred())
	})
})