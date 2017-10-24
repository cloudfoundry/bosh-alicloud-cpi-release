/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package action

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)


var _ = Describe("cpi:delete_disk", func() {
	It("can delete disk", func() {
		cid, _ := mockContext.NewDisk("")

		_, err := caller.Call("delete_disk", cid)
		Expect(err).NotTo(HaveOccurred())
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
	//				"i-2zefl7hfr7yb97ni5skw"
	//			],
	//			"context": {
	//			"director_uuid": "911133bb-7d44-4811-bf8a-b215608bf084"
	//		}
	//	}`))
	//	Expect(r.GetError()).NotTo(HaveOccurred())
	//})
})