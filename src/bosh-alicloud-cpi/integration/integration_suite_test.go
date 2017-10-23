/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package integration

import (
	"testing"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"bosh-alicloud-cpi/action"
	"bosh-alicloud-cpi/alicloud"
	"os"
	"bosh-alicloud-cpi/mock"
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

var caller action.Caller

var configForIntergation = string(`{
    "cloud": {
        "plugin": "alicloud",
        "properties": {
            "alicloud": {
                "region_id": "${REGION}",
                "access_key_id": "${ACCESS_KEY_ID}",
                "access_key_secret": "${ACCESS_KEY_SECRET}"
            },
            "registry": {
                "user": "registry",
                "password": "registry",
                "protocol": "http",
                "address": "127.0.0.1",
                "port": 25777
            },
            "agent": {
                "ntp": ["0.pool.ntp.org", "1.pool.ntp.org"],
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
}`)

var _ = BeforeSuite(func() {
	b, err := mock.NewBuilder(configForIntergation).ApplyEnvs()
	Expect(err).NotTo(HaveOccurred())

	config, err := alicloud.NewConfigFromBytes(b.ToBytes())
	Expect(err).NotTo(HaveOccurred())

	err = config.Validate()
	Expect(err).NotTo(HaveOccurred())

	logger := boshlog.NewWriterLogger(boshlog.LevelDebug, os.Stderr)
	caller = NewTestCaller(config, logger, mockContext)
})