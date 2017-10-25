/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package action

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"bosh-alicloud-cpi/mock"
)

var _ = Describe("cpi:has_disk", func() {
	It("can run with existent disk", func() {
		cid, _ := mockContext.NewDisk("")
		r, err := caller.CallGeneric("has_disk", cid)
		Expect(err ).NotTo(HaveOccurred())

		b := r.(bool)
		Expect(b).To(BeTrue())
	})
	It("can run with non-existed disk", func() {
		cid := mock.NewDiskId()
		r, err := caller.CallGeneric("has_disk", cid)
		Expect(err ).NotTo(HaveOccurred())

		b := r.(bool)
		Expect(b).To(BeFalse())
	})
})