/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package integration

import (
	"os"
	"testing"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"bosh-alicloud-cpi/action"
	"bosh-alicloud-cpi/alicloud"
	"bosh-alicloud-cpi/mock"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Integration Suite")
}

var caller action.Caller

var configForIntegration = string(`{
    "cloud": {
        "plugin": "alicloud",
        "properties": {
            "alicloud": {
                "region_id": "${CPI_REGION}",
				"zone_id": "${CPI_ZONE}",
                "access_key_id": "${CPI_ACCESS_KEY_ID}",
                "access_key_secret": "${CPI_ACCESS_KEY_SECRET}"
            },
            "registry": {
                "user": "${CPI_REGISTRY_USER}",
                "password": "${CPI_REGISTRY_PASSWORD},
                "protocol": "http",
                "address": "${CPI_REGISTRY_ADDRESS}",
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
	config, err := alicloud.NewConfigFromBytes([]byte(configForIntegration))
	Expect(err).NotTo(HaveOccurred())

	logger := boshlog.NewWriterLogger(boshlog.LevelDebug, os.Stderr)

	services := action.Services {
		Stemcells: alicloud.NewStemcellManager(config),
		Instances: alicloud.NewInstanceManager(config),
		Disks: alicloud.NewDiskManager(config),
		Networks: alicloud.NewNetworkManager(config),
		Registry: mock.NewRegistryMock(),
	}

	caller = action.NewCallerWithServices(config, logger, services)
})