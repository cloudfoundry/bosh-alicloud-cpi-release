/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package action

import (
	"bosh-alicloud-cpi/mock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("cpi:info", func() {
	It("can run with input", func() {
		r := caller.Run([]byte(`{
			"method":"info",
			"arguments":[],
			"context":{"director_uuid":"18badfe3-ec32-46eb-bbb5-fbdec1cfc5fa","request_id":"cpi-570352"}
		}`))
		Expect(r.GetError()).NotTo(HaveOccurred())
	})

	It("can run with non-existent vm", func() {
		cid := mock.NewInstanceId()
		r, err := caller.CallGeneric("has_vm", cid)
		Expect(err).NotTo(HaveOccurred())

		b := r.(bool)
		Expect(b).To(BeFalse())
	})
})
