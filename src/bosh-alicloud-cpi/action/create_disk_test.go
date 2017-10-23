/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package action

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)


var _ = Describe("cpi:create_disk", func() {
	It("can create disk with ", func() {
		By("create disk right")

		r := caller.Run([]byte(`{
			"method": "create_disk",
				"arguments": [
					30720,
					{},
				],
				"context": {
				"director_uuid": "911133bb-7d44-4811-bf8a-b215608bf084"
			}
		}`))
		Expect(r.GetError()).NotTo(HaveOccurred())
	})
	//It("create disk with instance", func() {
	//	By("create disk right")
	//
	//	r := caller.Run([]byte(`{
	//		"method": "create_disk",
	//			"arguments": [
	//				30720,
	//				{},
	//				"i-2zefl7hfr7yb97ni5skw"
	//			],
	//			"context": {
	//			"director_uuid": "911133bb-7d44-4811-bf8a-b215608bf084"
	//		}
	//	}`))
	//	Expect(r.GetError()).NotTo(HaveOccurred())
	//})
})