/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package action

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"os"
	"testing"
	"bosh-alicloud-cpi/mock"
	"bosh-alicloud-cpi/alicloud"
)

func TestActions(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Actions Suite")
}

var caller Caller
var mockContext mock.TestContext

var configForTest = []byte(`
{
    "cloud": {
        "plugin": "alicloud",
        "properties": {
            "alicloud": {
                "region_id": "cn-beijing",
                "access_key_id": "---",
                "access_key_secret": "---"
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
}
`)

var _ = BeforeSuite(func() {
	config, err := alicloud.NewConfigFromBytes(configForTest)
	Expect(err).NotTo(HaveOccurred())

	err = config.Validate()
	Expect(err).NotTo(HaveOccurred())

	logger := boshlog.NewWriterLogger(boshlog.LevelDebug, os.Stderr)
	caller = NewTestCaller(config, logger, mockContext)
})
