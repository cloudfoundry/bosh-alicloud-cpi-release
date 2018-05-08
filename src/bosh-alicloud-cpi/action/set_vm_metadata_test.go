/*
 * Copyright (C) 2017-2018 Alibaba Group Holding Limited
 */
package action

import (
	"bosh-alicloud-cpi/mock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("cpi:set_vm_metadata", func() {
	It("can run reboot vm", func() {
		id, _ := mockContext.NewInstance()
		r := caller.Run(mock.NewBuilder(`{
			"method": "set_vm_metadata",
			"arguments": [
				"${INST_ID}", {
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
		}`).P("INST_ID", id).ToBytes())
		Expect(r.GetError()).NotTo(HaveOccurred())

		inst, ok := mockContext.Instances[id]
		Expect(ok).To(BeTrue())

		Expect(inst.InstanceName).Should(Equal("dea_ng.4daa76e3-db6d-4550-b9b4-c504e7865f62"))
	})

	It("instance name can be normalized", func() {
		s := normalizeName("dea_ng/4daa76e3-db6d-4550-b9b4-c504e7865f62", "i_")
		Expect(s).Should(Equal("dea_ng.4daa76e3-db6d-4550-b9b4-c504e7865f62"))

		s = normalizeName("a", "i_")
		Expect(s).Should(Equal("i_a"))

		s = normalizeName("11123", "i_")
		Expect(s).Should(Equal("i_11123"))
	})
})
