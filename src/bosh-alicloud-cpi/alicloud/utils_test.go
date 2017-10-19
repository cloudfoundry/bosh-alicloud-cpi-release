/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package alicloud

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Utils", func() {
	It("GenerateHexString", func() {
		s := GenerateHexString(20)
		Expect(len(s)).To(20)
	})
})