/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package mock

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("mock:TestContext", func() {
	It("foo", func() {
		Expect(true).To(Equal(true))
	})
})
