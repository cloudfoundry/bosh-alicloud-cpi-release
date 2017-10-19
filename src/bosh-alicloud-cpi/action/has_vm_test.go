/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package action

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("cpi:has_vm", func() {
	It("can run when vm exists", func() {
		r :=  caller.RunTest([]byte(`


		`), )
		Expect(r.Error).NotTo(HaveOccurred())
	})

	It("can run when vm not exists", func() {
		r :=  caller.RunTest([]byte(`


		`), )
		Expect(r.Error).NotTo(HaveOccurred())
	})
})