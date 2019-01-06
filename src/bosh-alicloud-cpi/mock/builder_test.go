/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package mock

import (
	"encoding/json"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type TestJson struct {
	Key    string `json:"Key"`
	Secret string `json:Secret`
}

var _ = Describe("mock:builder", func() {
	It("should make json with parameter for input", func() {
		By("generate test parameters")
		key := RandomString(8)
		secret := RandomString(32)

		Expect(len(key)).To(Equal(8))
		Expect(len(secret)).To(Equal(32))

		By("build json and parse")
		bytes := NewBuilder(`{
			"Key": "${KEY}",
			"Secret": "${SECRET}"
		}`).
			P("KEY", key).
			P("${SECRET}", secret). // compatible for key surround with ${} or not
			ToBytes()

		var v TestJson
		err := json.Unmarshal(bytes, &v)
		Expect(err).NotTo(HaveOccurred())

		By("with correct field values")
		Expect(key).To(Equal(v.Key))
		Expect(secret).To(Equal(v.Secret))
	})

	It("should apply env", func() {

	})
})
