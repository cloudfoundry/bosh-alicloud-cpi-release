/*
 * Copyright (C) 2017-2018 Alibaba Group Holding Limited
 */
package integration

import (
	"bosh-alicloud-cpi/action"
	"bosh-alicloud-cpi/alicloud"
	"bosh-alicloud-cpi/mock"
	"fmt"
	"os"
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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
                "region": "${CPI_REGION}",
		        "availability_zone": "${CPI_ZONE}",
                "access_key_id": "${CPI_ACCESS_KEY_ID}",
                "access_key_secret": "${CPI_ACCESS_KEY_SECRET}"
            },
            "registry": {
                "user": "${CPI_REGISTRY_USER}",
                "password": "${CPI_REGISTRY_PASSWORD}",
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

	err = ApplySystemEnv(&config)
	Expect(err).NotTo(HaveOccurred())

	logger := boshlog.NewWriterLogger(boshlog.LevelDebug, os.Stderr)

	services := action.Services{
		Stemcells: alicloud.NewStemcellManager(config, logger),
		Osses:     alicloud.NewOssManager(config, logger),
		Instances: alicloud.NewInstanceManager(config, logger),
		Disks:     alicloud.NewDiskManager(config, logger),
		Networks:  alicloud.NewNetworkManager(config, logger),
		Registry:  mock.NewRegistryMock(),
	}

	caller = action.NewCallerWithServices(config, logger, services)

	err = CleanInstances(config, services.Instances)
	Expect(err).NotTo(HaveOccurred())
})

func CleanInstances(config alicloud.Config, manager alicloud.InstanceManager) error {
	client, err := config.NewEcsClient("")
	if err != nil {
		return err
	}
	args := ecs.CreateDescribeInstancesRequest()
	args.PrivateIpAddresses = `["` + internalIp + `"]`
	args.VSwitchId = vswitchId

	insts, err := client.DescribeInstances(args)
	if err != nil {
		return fmt.Errorf("CleanInstances try DescribeInstances failed %s", err.Error())
	}

	if insts == nil || len(insts.Instances.Instance) <= 0 {
		return nil
	}
	for _, inst := range insts.Instances.Instance {
		_, err := caller.Call("delete_vm", inst.InstanceId)
		if err != nil {
			return fmt.Errorf("CleanInstances try delete_vm %s failed %s", inst.InstanceId, err.Error())
		}
	}
	return nil
}
