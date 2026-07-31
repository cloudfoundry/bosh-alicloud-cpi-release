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

	It("deletes backing snapshots when the stemcell is deleted", func() {
		ss1 := mock.NewSnapshotId()
		ss2 := mock.NewSnapshotId()
		id, _ := mockContext.NewStemcellWithSnapshots([]string{ss1, ss2})

		_, err := caller.Call("delete_stemcell", id)
		Expect(err).NotTo(HaveOccurred())

		_, ok := mockContext.Stemcells[id]
		Expect(ok).To(BeFalse(), "stemcell should be removed")

		_, ok = mockContext.Snapshots[ss1]
		Expect(ok).To(BeFalse(), "snapshot ss1 should be removed")

		_, ok = mockContext.Snapshots[ss2]
		Expect(ok).To(BeFalse(), "snapshot ss2 should be removed")
	})

	It("succeeds silently when the stemcell does not exist", func() {
		_, err := caller.Call("delete_stemcell", mock.NewStemcellId())
		Expect(err).To(HaveOccurred())
	})
})
