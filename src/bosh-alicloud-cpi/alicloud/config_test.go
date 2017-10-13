package alicloud

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Testing with Ginkgo", func() {
	var (
		config Config
		err error
	)

	BeforeEach(func() {
		config, err = NewConfigFromBytes([]byte(`
			{
				"cloud": {
					"plugin": "alicloud",
					"properties": {
						"alicloud": {
							"region_id": "cn-beijing",
							"zone_id": "cn-beijing-a",
							"access_key_id": "OAtTkdQDHJyYGCmU",
							"access_key_secret": "PzosfNHKsG6sFQKNnxvMjqQNx8x2of"
						},
						"registry": {
							"user": "registry",
							"password": "2a57f7c0-7726-4e76-43aa-00b10b073229",
							"protocol": "http",
							"address": "10.0.0.2",
							"port": "6901"
						},
						"agent": {
							"ntp": "ntp",
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
			}
		`))

	})

	It("loading from JSON", func() {
		Context("when the JSON parses succesfully", func() {
			It("should populate the fields correctly", func() {
				Expect(config.OpenApi.RegionId).To(Equal("cn-beijing"))
			})

			It("should not error", func() {
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})
