/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package action

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"bosh-alicloud-cpi/mock"
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