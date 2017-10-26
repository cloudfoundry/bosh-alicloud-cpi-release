/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package integration

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"bosh-alicloud-cpi/mock"
	"time"
)

var _ = Describe("integration:vm", func() {
	It("can run the vm lifecycle", func() {
		By("create vm")
		in := mock.NewBuilder(`{
			"method": "create_vm",
			"arguments": [
				"be387a69-c5d5-4b94-86c2-978581354b50",
				"m-2zehhdtfg22hq46reabf",
				{
					"ephemeral_disk": {
						"size": "40_960",
						"category": "cloud_efficiency"
					},
					"image_id": "${STEMCELL_ID}",
					"instance_name": "test-cc",
					"instance_type": "ecs.n4.small",
					"system_disk": {
						"size": "61_440",
						"category": "cloud_efficiency"
					}
				},
				{
					"private": {
						"type": "manual",
						"ip": "${INTERNAL_IP}",
						"netmask": "${INTERNAL_NETMASK}",
						"cloud_properties": {
							"security_group_id": "${SECURITY_GROUP_ID}",
							"vswitch_id": "${VSWITCH_ID}"
						},
						"default": [
							"dns",
							"gateway"
						],
						"dns": [
							"8.8.8.8"
						],
						"gateway": "${INTERNAL_GW}"
					}
				},
				[],
				{}
			],
			"context": {
				"director_uuid": "911133bb-7d44-4811-bf8a-b215608bf084"
			}
		}`).
		P("STEMCELL_ID", stemcellId).
		P("SECURITY_GROUP_ID", securityGroupId).
		P("VSWITCH_ID", vswitchId).
		P("INTERNAL_IP", internalIp).
		P("INTERNAL_NETMASK", internalNetmask).
		P("INTERNAL_GW", internalGw).
		ToBytes()

		r := caller.Run(in)
		Expect(r.GetError()).NotTo(HaveOccurred())
		cid := r.GetResultString()

		By("sleep 90 seconds to make sure vm is fully started")
		time.Sleep(time.Duration(90) * time.Second)

		By("delete vm")
		_, err := caller.Call("delete_vm", cid)
		Expect(err).NotTo(HaveOccurred())

		By("vm should not exists")
		exists, err := caller.CallGeneric("has_vm", cid)
		Expect(err).NotTo(HaveOccurred())
		Expect(exists).To(BeFalse())
	})

	It("can run the vm lifecycle with persistent disk", func() {
		//By("create persistent disk")
		//caller.Call("create_disk", 30720, "{}", )
		//
		//By("create vm with persistent disk")
		//
		//By("delete vm")
		//
		// By("delete disk")
	})


	It("can run the vm lifecycle with persistent disk", func() {
		//By("create vm")
		//By("create vm")
		//in := mock.NewBuilder(`{
		//	"method": "create_vm",
		//	"arguments": [
		//		"be387a69-c5d5-4b94-86c2-978581354b50",
		//		"m-2zehhdtfg22hq46reabf",
		//		{
		//			"ephemeral_disk": {
		//				"size": "40_960",
		//				"type": "cloud_efficiency"
		//			},
		//			"image_id": "${STEMCELL_ID}",
		//			"instance_name": "test-cc",
		//			"instance_type": "ecs.n4.small",
		//			"system_disk": {
		//				"size": "61_440",
		//				"type": "cloud_efficiency"
		//			}
		//		},
		//		{
		//			"private": {
		//				"type": "manual"
		//				"ip": "${NETWORK_IP}",
		//				"netmask": "255.240.0.0",
		//				"cloud_properties": {
		//					"security_group_id": "${SECURITY_GROUP_ID}",
		//					"vswitch_id": "${VSWITCH_ID}"
		//				},
		//				"default": [
		//					"dns",
		//					"gateway"
		//				],
		//				"dns": [
		//					"8.8.8.8"
		//				],
		//				"gateway": "${NETWORK_GATEWAY}""
		//			}
		//		},
		//		[],
		//		{}
		//	],
		//	"context": {
		//		"director_uuid": "911133bb-7d44-4811-bf8a-b215608bf084"
		//	}
		//}`).
		//P("STEMCELL_ID", stemcellId).
		//P("NETWORK_IP", networkAddress).
		//P("SECURITY_GROUP_ID", securityGroupId).
		//P("VSWITCH_ID", vswitchId).
		//P("NETWORK_GETWAY", networkGateway).
		//ToBytes()
		//
		//r := caller.Run(in)
		//Expect(r.GetError()).NotTo(HaveOccurred())
		//instCid := r.GetResultString()
		//
		//By("create persistent disk")
		//diskCid, err := caller.Call("create_disk", 30720, `{ "type": "cloud_ssd"}`, instCid)
		//Expect(err).NotTo(HaveOccurred())
		//
		//By("attach disk")
		//_, err = caller.Call("attach_disk", instCid, diskCid)
		//Expect(err).NotTo(HaveOccurred())
		//
		//By("verify disks")
		//r2, err := caller.CallGeneric("get_disks", instCid)
		//Expect(err).NotTo(HaveOccurred())
		//disks := r2.([]interface{})
		//Expect(disks).Should(ConsistOf(diskCid))
		//
		//By("detach disk")
		//
		//By("delete vm")
		//_, err := caller.
		//
		//By("delete disk")
	})

	//It("can create vm with manual ip, and delete it", func() {})
	//It("can create vm with dynamic ip, and delete it", func() {})
	//It("can create vm with external ip, and delete it", func() {})
	//It("can create vm with key pair, and delete it", func() {})
	//It("can create vm, then start, stop and delete it", func() {})
	//It("can create vm, then start, reboot, stop and delete it", func() {})
	//It("can create vm, then start, reboot, stop and delete it", func() {})

	//It("creates a VM with an invalid configuration and receives an error message with logs", func() {
	//	testing.Short()
	//	request := fmt.Sprintf(`{
	//		"method": "create_vm",
	//		"arguments": [
	//		"7bc16fab-52c3-4bb9-a5c3-560445986860",
	//			"%v",
	//			{
	//				"ephemeral_disk": {
	//					"size": 50,
	//					"type": "cloud_efficiency"
	//				},
	//				"image_id": "m-temp1234",
	//				"instance_type": "ecs.n4.large",
	//				"system_disk": {
	//					"size": 50,
	//					"type": "cloud_efficiency"
	//				}
	//			},
	//			{
	//				"private": {
	//					"ip": "172.16.0.101",
	//					"netmask": "255.240.0.0",
	//					"cloud_properties": {
	//						"security_group_id": "%v",
	//						"vswitch_id": "%v"
	//					},
	//					"default": [
	//						"dns",
	//						"gateway"
	//					],
	//					"gateway": "172.16.0.1"
	//				}
	//			}
	//		]
	//		}`, boshStemcellImageId, securityGroupId, vswitchId)
	//	resp, err := execCPI(request)
	//	Expect(err).ToNot(HaveOccurred())
	//	Expect(resp.Error.Message).ToNot(BeEmpty())
	//})
	//
	//It("executes the VM lifecycle", func() {
	//	var vmCID string
	//	By("creating a VM")
	//	request := fmt.Sprintf(`{
	//		"method": "create_vm",
	//		"arguments": [
	//			"182a951a-2f8e-4d22-6489-d78b4a8b6f8a",
	//			"m-2zeggz4i4n2z510ajcvw",
	//			{
	//				"availability_zone": "%v",
	//				"ephemeral_disk": {
	//					"size": "100",
	//					"type": "cloud_efficiency"
	//				},
	//				"halt_mark": "true",
	//				"instance_charge_type": "PostPaid",
	//				"instance_type": "ecs.n4.large"
	//			},
	//			{
	//				"public": {
	//					"cloud_properties": {
	//						"internet_charge_type": "PayByTraffic"
	//					},
	//					"ip": "%v",
	//					"type": "vip"
	//				},
	//				"default": {
	//					"cloud_properties": {
	//						"internet_charge_type": "PayByTraffic",
	//						"security_group_id": "%v",
	//						"vswitch_id": "%v"
	//					},
	//					"default": [
	//						"dns",
	//						"gateway"
	//					],
	//					"dns": [
	//						"8.8.8.8"
	//					],
	//					"gateway": "172.16.0.1",
	//					"ip": "172.16.0.3",
	//					"netmask": "255.255.255.0",
	//					"type": "manual"
	//				}
	//			},
	//			[
	//			],
	//			{
	//				"bosh": {
	//					"mbus": {
	//						"cert": {
	//							"ca": "-----BEGIN CERTIFICATE-----\nMIIDEzCCAfugAwIBAgIQe1NhaUZY50HsFPOw5zhmzzANBgkqhkiG9w0BAQsFADAz\nMQwwCgYDVQQGEwNVU0ExFjAUBgNVBAoTDUNsb3VkIEZvdW5kcnkxCzAJBgNVBAMT\nAmNhMB4XDTE3MDkyMTA4MjA1OVoXDTE4MDkyMTA4MjA1OVowMzEMMAoGA1UEBhMD\nVVNBMRYwFAYDVQQKEw1DbG91ZCBGb3VuZHJ5MQswCQYDVQQDEwJjYTCCASIwDQYJ\nKoZIhvcNAQEBBQADggEPADCCAQoCggEBAL2gfrohAQm2E5LYzqC0QEER7HktvhAl\n/I+m0MbDZwHH/dXUjPV/5+Xi2w7X1llqnEyDP2cvWbJ4EkSqsaG7UMwb+7sbkLGL\ncG/BJK/mFVGLLPmpln3ZnQ9zzIQ46sS8Dxy6ViV9oK53XCe2uHphHqBNNJ9NHwrp\nx+cSrBjRAQmH2r4KikHqIngEVX2qN++8ZS4nrw/7WRI90Scd2YxJkUW/HTaklKW0\n4PBtVUBRxAe1L/MRUF9T5lgzJbiVDX+0XHQ58HryC6uIyQzAzw9oNyn83ymTGg+a\n1Ni9xGKIcEYZIzQdZvJXl2huLExyVnNNKeGmLnpGRS0no7GMigp1cWcCAwEAAaMj\nMCEwDgYDVR0PAQH/BAQDAgEGMA8GA1UdEwEB/wQFMAMBAf8wDQYJKoZIhvcNAQEL\nBQADggEBAECF/zgrONd+EH+IMCpmWDqwizXx5IQD+iSVTTAWTIqOfmm6WvHCLCLP\nplmh+bB/PqlCxjRl7X7MpeYNVQl1arenvzOGbDtH3h7lxW8wSYiJXslY3+HS0KWF\nF2G6XkT7Fz/YPsE/CiccUwAiADJWwWnqotr1jvNRyccCTUEVU4zkk2x0V4NYYRxh\nKQcV2gphicz3gbyhQyuRRmrfvBtrWkdwxvLDbfxeLFpw7PS78SEesD1AevUiuV/A\nJv6wI7o5843QRkFcEX3wfn4co8yMRW5GV2pXfFvCQirvuEmqDLYMKcc56xAOXfjD\nYoRwQYzxgLuWbDv002upSjw4zPY7PVI=\n-----END CERTIFICATE-----\n",
	//							"certificate": "-----BEGIN CERTIFICATE-----\nMIIDPzCCAiegAwIBAgIRANufxsY2UPV9Py5YbkYKAxEwDQYJKoZIhvcNAQELBQAw\nMzEMMAoGA1UEBhMDVVNBMRYwFAYDVQQKEw1DbG91ZCBGb3VuZHJ5MQswCQYDVQQD\nEwJjYTAeFw0xNzA5MjEwODIwNTlaFw0xODA5MjEwODIwNTlaMDsxDDAKBgNVBAYT\nA1VTQTEWMBQGA1UEChMNQ2xvdWQgRm91bmRyeTETMBEGA1UEAxMKMTcyLjE2LjAu\nMjCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBALLgCxMxoqg/64gj1jC6\ng9s148lMKi6zcI/LdUpTd2R44Do8S8E9z5UugmpG2itqijx9Od6oA+KwgbcqIosd\n5kiwQuGhdp4otTwQPKUbtCghKMG0qN41nB3zFXaSsOQUvVZecppj08ILWYnyb3QW\nWTZc1AyCwOCUjU4+cl6Zu5i0hqMwRSWF/qrYttN1OxNHfltLlCqYiyiizFqvvE7a\njp6kL9o/hG3EyAznQRaE+uUPzUEHeRD+jLiekwT4GknNt2nxuF3+3z21zyivgAJp\nvkNMwoCOelaMFShPyRGQceHBED6zEqDP5E4dqs+qhYscYGMWVcNWD2iSdAt+bOy+\nFTUCAwEAAaNGMEQwDgYDVR0PAQH/BAQDAgWgMBMGA1UdJQQMMAoGCCsGAQUFBwMB\nMAwGA1UdEwEB/wQCMAAwDwYDVR0RBAgwBocErBAAAjANBgkqhkiG9w0BAQsFAAOC\nAQEAOQxjxoB4JnJ9syLSMZTPdjDImXgnVmtB+hRax/ZsWsrlWWgsLV09SxtvFjg+\nHr1IeMSEg0zpIBJ9btseQTzLBaR3t4Qbeg18q2GVKD3kU6KeW5ucbpa0IzIBx51S\n2DzpjP0Eb0VyweAFMZi18OWINhdFTrOdf6pTa/H9G/E0Uu2rdqh5gKBFRrfBd+Hl\n3V0/wxih2knR+qLcMVnr0kpSZtuKmVzqWOzWa0OWd1iVIrFSkPousV29AQtdG4rO\nTCQ4VmmSfhBNJr8BlVi23r3WxkKMJWtUSAQg5gyF+KnwWzaGAeNLf7suQUFJingK\n134NfboZBF4Z3A3MRXez5xBNYQ==\n-----END CERTIFICATE-----\n",
	//							"private_key": "-----BEGIN RSA PRIVATE KEY-----\nMIIEogIBAAKCAQEAsuALEzGiqD/riCPWMLqD2zXjyUwqLrNwj8t1SlN3ZHjgOjxL\nwT3PlS6CakbaK2qKPH053qgD4rCBtyoiix3mSLBC4aF2nii1PBA8pRu0KCEowbSo\n3jWcHfMVdpKw5BS9Vl5ymmPTwgtZifJvdBZZNlzUDILA4JSNTj5yXpm7mLSGozBF\nJYX+qti203U7E0d+W0uUKpiLKKLMWq+8TtqOnqQv2j+EbcTIDOdBFoT65Q/NQQd5\nEP6MuJ6TBPgaSc23afG4Xf7fPbXPKK+AAmm+Q0zCgI56VowVKE/JEZBx4cEQPrMS\noM/kTh2qz6qFixxgYxZVw1YPaJJ0C35s7L4VNQIDAQABAoIBAFeKwqDQJ/UD43er\nYkZS4flEtIht2C8m7q3RO0P2+XWYmtSlccXPRGqUaossxdV9vM3B07Kes9gb3kAQ\nRPuk1HE6omDerrjU323X3HZJyq/hGptCmWq2/gLCVvzC6gOWCtvcOWZJ+Pb8qwOS\nPO2pilvKrpS44UCIM2fZtAuMXX1r+hf91FwxHWrHlo3Sb0I4sbdQIcomskZybGUE\nQTfzGuhJIFa9bLeHV81mWqOhMP636eMHZ7F4EYpsK9S8D4l+FgNX21wB/6zIXiRm\nxqp3HS1ryF1GXusPxuSomJbmk0cG6OpV5PFX4FJCvMWXmuVDwJOces4vQKgElTLa\ncYjSPbkCgYEA613dC2P7z321oetrxUUOgGH9uQ3hzsSPj/IlMd4/uqBcr90NW8uM\nureldgrpZ0oxKRKP8vUwOTPK9RjotuBoQPrI2BpeCGkjBfhq99B/WrZfz1doJ6jW\na71xsTyTlfjbpNW/Xh/clqUvlS632O4eGwPTq69LLf7djkW4gdItc48CgYEAwo5j\nhWen5QtqPAVKFgmrZmd/nri3mDGHyq0Eg0lpr9d9HW8cQv9AV56K0u5liVt9Lo7X\neOJbUfNXY9+DlNpaqNmmhA/+KXL1ckNumijVl4GNXVa367SJSOaygQA02ymkFUd4\njj8MWiRVCcP3y2YrvDDAYuTCZsIl8Dudj4NxuPsCgYACTypzCSkYURBuJUQqbFIH\nGm8F2MgFYlJSRDrvMVIIv7gJFa8i3m1kC5c5AERn+gdfcsosxRETDpoIK5Vk7fC3\n6n37+M5BYN6yGUzbX5VQS4fHHgFsmjB4YCR0a7a6+vUUufAluURNyhMccJfnLfbn\npvL1tUOUkPKVicOUqn49qwKBgE+7XtnLMylQ1kamvEfvyoh7HfgEJ2l90vKimVjc\ney2PGD05zdE/HjVKSgZLoNz72397Fp751QbuvP+3GAumuMS9/dndXAHMlP4w2GDh\nHzep5i88XL+CC0kPElR/qymuFQqLccKJ4BwJC7im0SRQSNgk+pMMwQavxjB/ngC0\nk6SFAoGAXvVfNasp3FkidMFBBlHA/2OVJhHfBFrupeQYwggW7WUfH69H5dop5Uu6\nAnOJSQi1FL4soc5Db4hECu8x04+CZdfaBippOHG8IxV/5Sfavhqw6q30e8L0uKu4\n8XqFMVz83X5jfBBDUZT1q8D40Q8GCVliajXDdG2UDcbs383t23M=\n-----END RSA PRIVATE KEY-----\n"
	//						}
	//					},
	//					"password": "*"
	//				}
	//			}
	//		],
	//		"context": {
	//			"director_uuid": "478b5c95-c143-4223-737f-7c1c834eebc0"
	//		}
	//	}`, "cn-beijing-c", "47.94.216.146", "sg-2zec8ubi1q5aeo5mqcbb", "vsw-2zevwt3w7h5u761o405rd")
	//	vmCID = assertSucceedsWithResult(request).(string)
	//	log.Printf("VM cid:", vmCID)
		//By("locating the VM")
		//request = fmt.Sprintf(`{
		//	  "method": "has_vm",
		//	  "arguments": ["%v"]
		//	}`, vmCID)
		//exists := assertSucceedsWithResult(request).(bool)
		//Expect(exists).To(Equal(true))
		//
		////expectLabels := map[string]string{
		////	"label-1-key": "label-1-value",
		////	"label-2-key": "label-2-value",
		////}
		////assertValidVMB(vmCID, func(instance *computebeta.Instance) {
		////	// Labels should be an exact match
		////	Expect(instance.Labels).To(BeEquivalentTo(expectLabels))
		////})
		//
		//m := map[string]string{
		//	"director":           "val-that-is-definitely-for-sure-absolutely-longer-than-the-allowable-enforced-63-char-limit-and-should-be-truncated",
		//	"name":               "val_with_underscores_ending_in_dash-",
		//	"deployment":         "deployment-name",
		//	"job":                "job-name",
		//	"index":              "0",
		//	"integration-delete": "",
		//}
		////expectLabels = map[string]string{
		////	"director":    "val-that-is-definitely-for-sure-absolutely-longer-than-the-al",
		////	"name":        "val-with-underscores-ending-in-dash",
		////	"deployment":  "deployment-name",
		////	"job":         "job-name",
		////	"index":       "n0",
		////	"label-1-key": "label-1-value",
		////	"label-2-key": "label-2-value",
		////}
		//mj, _ := json.Marshal(m)
		//request = fmt.Sprintf(`{
		//	  "method": "set_vm_metadata",
		//	  "arguments": [
		//		"%v",
		//		%v
		//	  ]
		//	}`, vmCID, string(mj))
		//assertSucceeds(request)
		////assertValidVMB(vmCID, func(instance *computebeta.Instance) {
		////	// Labels should be an exact match
		////	Expect(instance.Labels).To(BeEquivalentTo(expectLabels))
		////})
		//
		//By("rebooting the VM")
		//request = fmt.Sprintf(`{
		//	  "method": "reboot_vm",
		//	  "arguments": ["%v"]
		//	}`, vmCID)
		//assertSucceeds(request)
		//
		//By("deleting the VM")
		//request = fmt.Sprintf(`{
		//	  "method": "delete_vm",
		//	  "arguments": ["%v"]
		//	}`, vmCID)
		//assertSucceeds(request)

	//})
})
