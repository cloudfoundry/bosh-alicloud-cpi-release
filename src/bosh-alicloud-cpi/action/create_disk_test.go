/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package action

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("cpi:create_disk", func() {
	It("can create disk", func() {
		By("create disk right")
		instCid, _ := mockContext.NewInstance()

		cid, err := caller.Call("create_disk", 30720, `{ "type": "cloud" }`, instCid)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(cid)).To(Equal(22))

		//r = caller.Call("has_disk", cid)
		//Expect(r.GetError()).NotTo(HaveOccurred())
		//Expect(r.Result)
	})

	It("can create disk with small size", func() {
		instCid, _ := mockContext.NewInstance()

		cid, err := caller.Call("create_disk", 1024, `{ "type": "cloud" }`, instCid)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(cid)).To(Equal(22))

		//r = caller.Call("has_disk", cid)
		//Expect(r.GetError()).NotTo(HaveOccurred())
		//Expect(r.Result)
	})

	//It("create disk with instance", func() {
	//	By("create disk right")
	//
	//	r := caller.Run([]byte(`{
	//		"method": "create_disk",
	//			"arguments": [
	//				30720,
	//				{},
	//				""
	//			],
	//			"context": {
	//			"director_uuid": "911133bb-7d44-4811-bf8a-b215608bf084"
	//		}
	//	}`))
	//	Expect(r.GetError()).NotTo(HaveOccurred())
	//})
})
