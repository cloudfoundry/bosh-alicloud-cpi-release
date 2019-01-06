/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package action

import (
	"bosh-alicloud-cpi/mock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("cpi:delete_snapshot", func() {
	It("can delete exists snapshot", func() {
		ssid := mockContext.NewSnapshot("")
		_, err := caller.Call("delete_snapshot", ssid)
		Expect(err).NotTo(HaveOccurred())

		_, ok := mockContext.Snapshots[ssid]
		Expect(ok).Should(BeFalse())
	})

	It("raise error with non exists snapshot", func() {
		ssid := mock.NewSnapshotId()
		_, err := caller.Call("delete_snapshot", ssid)
		Expect(err).To(HaveOccurred())
	})
})
