/*
 * Copyright (C) 2017-2018 Alibaba Group Holding Limited
 */
package registry

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Sample", func() {
	BeforeEach(func() {

	})

	It("Short", func() {
		s := "Hello"
		Expect(s).To(Equal("Hello"))
	})

	//It("Integration", func() {
	//	s := testing.Short()
	//	Expect(s).To(Equal(true))
	//})
})
