/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package alicloud

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"github.com/denverdino/aliyungo/ecs"
	. "github.com/onsi/ginkgo"
	"os"
)

var TestConfig = []byte(`
{
    "cloud": {
        "plugin": "alicloud",
        "properties": {
            "alicloud": {
                "region_id": "cn-beijing",
                "zone_id": "cn-beijing-a",
                "access_key_id": "${ACCESS_KEY_ID}",
                "access_key_secret": "${ACCESS_KEY_CONFIG}"
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
`)

var _ = Describe("Testing with Ginkgo", func() {
	It("get vm status", func() {

		config, _ := NewConfigFromBytes(TestConfig)
		logger := boshlog.NewWriterLogger(boshlog.LevelDebug, os.Stderr)
		runner := NewRunner(logger, config)

		status, err := runner.GetInstanceStatus(testInstId)

		if err != nil {
			GinkgoT().Error("VmStatus", "TestGetVmStatus failed, ", err)
		}

		logger.Info("VmStatus:", "Status %s", status)
		GinkgoT().Log(status)
	})
	It("stop start vm", func() {

		config, _ := NewConfigFromBytes(TestConfig)
		logger := boshlog.NewWriterLogger(boshlog.LevelDebug, os.Stderr)
		runner := NewRunner(logger, config)
		GinkgoT().Log("Stopping VM ...")
		err := runner.StopInstance(testInstId)
		if err != nil {
			GinkgoT().Error("VmStatus", "Stoping failed, ", err)
		}
		err = runner.WaitForInstanceStatus(testInstId, ecs.Stopped)
		if err != nil {
			GinkgoT().Error("VmStatus", "WaitForVm failed, ", err)
		}
		GinkgoT().Log("Stopping VM Done")
		GinkgoT().Log("Starting VM ...")
		err = runner.StartInstance(testInstId)
		if err != nil {
			GinkgoT().Error("VmStatus", "Starting failed, ", err)
		}
		err = runner.WaitForInstanceStatus(testInstId, ecs.Running)
		if err != nil {
			GinkgoT().Error("VmStatus", "WaitForVm failed, ", err)
		}
		GinkgoT().Log("Stopping VM Done")
	})
})
var testInstId = "i-2zefl7hfr7yb97ni5skw"
