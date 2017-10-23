/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package action

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"bosh-alicloud-cpi/mock"
)

var _ = Describe("cpi:has_vm", func() {
	It("can run with existent vm", func() {
		cid, _ := mockContext.NewInstance()
		r, err := caller.CallGeneric("has_vm", cid)
		Expect(err).NotTo(HaveOccurred())

		b := r.(bool)
		Expect(b).To(BeTrue())
	})

	It("can run with non-existent vm", func() {
		cid := mock.NewInstanceId()
		r, err := caller.CallGeneric("has_vm", cid)
		Expect(err).NotTo(HaveOccurred())

		b := r.(bool)
		Expect(b).To(BeFalse())
	})
})