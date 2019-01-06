/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package action

import (
	"bosh-alicloud-cpi/mock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
