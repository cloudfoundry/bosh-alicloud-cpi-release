/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package alicloud

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("LoadConfig from Json", func() {
	var (
		config Config
		err error
	)

	BeforeEach(func() {
		config, err = NewConfigFromBytes([]byte(`{
			"cloud": {
				"plugin": "alicloud",
				"properties": {
					"alicloud": {
						"region_id": "cn-beijing",
						"zone_id": "cn-beijing-a",
						"endpoint": "cn-zhangjiakou.aliyuncs.com",
						"access_key_id": "${ACCESS_KEY_ID}",
						"access_key_secret": "${ACCESS_KEY_SECRET}"
					},
					"registry": {
						"user": "registry",
						"password": "registry-password",
						"protocol": "http",
						"address": "10.0.0.2",
						"port": "6901"
					},
					"agent": {
						"ntp": ["ntp", "ntp"],
						"mbus": "http://mbus:mbus@0.0.0.0:6868",
						"blobstore": {
							"provider": "dav",
							"options": {
								"endpoint": "http://10.0.0.2:25250",
								"user": "agent",
								"password": "agent-password"
							}
						}
					}
				}
			}
		}`))
	})

	It("unmarshal without error", func() {
		Expect(err).NotTo(HaveOccurred())
	})

	It("validate without error", func() {
		err := config.Validate()
		Expect(err).NotTo(HaveOccurred())
	})

	It("generate registry settings correctly", func() {
		s1 := RegistryConfig{
			Protocol: "http",
			User:"registry",
			Password: "registry-password",
			Host: "10.0.0.2",
			Port: "6901",
		}
		s2 := `{"registry":{"endpoint":"http://registry:registry-password@10.0.0.2:6901"}}`
		Expect(s2).To(Equal(s1.ToInstanceUserData()))
	})

	//It("", )
	//
	//It("loading from JSON", func() {
	//	Context("when the JSON parses succesfully", func() {
	//		It("should populate the fields correctly", func() {
	//			Expect(config.OpenApi.RegionId).To(Equal("cn-beijing"))
	//		})
	//
	//		It("should not error", func() {
	//			Expect(err).NotTo(HaveOccurred())
	//		})
	//	})
	//})
})
