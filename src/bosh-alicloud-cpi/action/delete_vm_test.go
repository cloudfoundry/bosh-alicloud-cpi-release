/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package action

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("cpi:delete_vm", func() {
	It("can delete_vm with right registry", func() {
		By("attach disk")
		r :=  caller.RunTest([]byte(`


		`), )
		Expect(r.GetError()).NotTo(HaveOccurred())
		//
		// TODO: use mock method to detect execution results
		// disks := caller.Disks.GetDisk()
		// Expect(disks.GetDiskStatus(id)).To(Equal())

		By("update registry right")

	})
})