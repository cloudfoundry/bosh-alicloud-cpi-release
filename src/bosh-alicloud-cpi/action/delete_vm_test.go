/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package action

import (
	"bosh-alicloud-cpi/mock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("cpi:delete_vm", func() {
	It("can delete existent vm", func() {
		cid, _ := mockContext.NewInstance()

		_, err := caller.Call("delete_vm", cid)
		Expect(err).NotTo(HaveOccurred())
	})

	It("no error with non-exist vm", func() {
		cid := mock.NewInstanceId()
		_, err := caller.Call("delete_vm", cid)
		Expect(err).NotTo(HaveOccurred())
	})
})
