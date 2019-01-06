/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package action

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("cpi:snapshot_disk", func() {
	It("can create snapshot with exists disk", func() {
		diskCid, _ := mockContext.NewDisk("")
		ssid, err := caller.Call("snapshot_disk", diskCid, nil)
		Expect(err).NotTo(HaveOccurred())

		s, ok := mockContext.Snapshots[ssid]
		Expect(ok).Should(BeTrue())
		Expect(s).Should(Equal(diskCid))
	})
})
