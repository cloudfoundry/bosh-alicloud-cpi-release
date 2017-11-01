/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package action

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("cpi:set_vm_metadata", func() {
	It("can run reboot vm", func() {
		r :=  caller.Run([]byte(`{
			"method": "set_vm_metadata",
			"arguments": [
				"i-2ze1zz9zmbblcsk0e4i6", {
					"director": "my-bosh",
					"deployment": "ali-cf-215",
					"id": "4daa76e3-db6d-4550-b9b4-c504e7865f62",
					"job": "dea_ng",
					"index": "0",
					"name": "dea_ng/4daa76e3-db6d-4550-b9b4-c504e7865f62",
					"created_at": "2017-10-25T09:00:07Z"
				}
			],
			"context": {
					"director_uuid": "580da067-b2ff-4eb6-b271-23cc76409121",
					"request_id": "cpi-440280"
			}
		}`), )
		Expect(r.GetError()).NotTo(HaveOccurred())
	})
})


