/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package action

import (
	"bosh-alicloud-cpi/mock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("cpi:delete_stemcell", func() {
	It("deletes an existing stemcell", func() {
		id, _ := mockContext.NewStemcell()

		_, err := caller.Call("delete_stemcell", id)
		Expect(err).NotTo(HaveOccurred())

		_, ok := mockContext.Stemcells[id]
		Expect(ok).To(BeFalse())
	})

	It("returns an error when the stemcell does not exist", func() {
		// NOTE: the mock's DeleteStemcell errors on a missing image. The real
		// StemcellManagerImpl.DeleteStemcell is idempotent (returns nil) — that
		// path is covered by the alicloud-package unit test.
		_, err := caller.Call("delete_stemcell", mock.NewStemcellId())
		Expect(err).To(HaveOccurred())
	})
})
