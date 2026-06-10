/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package action

import (
	"bosh-alicloud-cpi/alicloud"
	"bosh-alicloud-cpi/mock"
	"fmt"
	"os"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// nvmeStemcellManagerSpy wraps StemcellManagerMock and records whether
// EnableNvmeSupport was called and with which image ID.
type nvmeStemcellManagerSpy struct {
	alicloud.StemcellManager
	nvmeCalledWith []string
	nvmeErr        error
}

func (s *nvmeStemcellManagerSpy) EnableNvmeSupport(imageId string) error {
	s.nvmeCalledWith = append(s.nvmeCalledWith, imageId)
	return s.nvmeErr
}

func (s *nvmeStemcellManagerSpy) CopyImage(args *ecs.CopyImageRequest) (string, error) {
	return s.StemcellManager.CopyImage(args)
}

func newEncryptedCallerWithSpy(spy *nvmeStemcellManagerSpy, config alicloud.Config) Caller {
	logger := boshlog.NewWriterLogger(boshlog.LevelWarn, os.Stderr)
	mc := mock.NewTestContext(config)
	services := Services{
		Stemcells: spy,
		Osses:     mock.NewOssManagerMock(mc),
		Instances: mock.NewInstanceManagerMock(mc),
		Disks:     mock.NewDiskManagerMock(mc),
		Networks:  mock.NewNetworkManagerMock(mc),
		Registry:  mock.NewRegistryMock(),
	}
	return NewCallerWithServices(config, logger, services)
}

var _ = Describe("create_stemcell", func() {
	It("can create stemcell", func() {
		r := caller.Run([]byte(`{
			"method": "create_stemcell",
			"arguments": [
				"/var/vcap/data/tmp/director/stemcell20170926-9684-17ncxdz/image",
				{
					"architecture": "x86_64",
					"container_format": "",
					"disk": "50",
					"disk_format": "rawdisk",
					"hypervisor": "kvm",
					"image_id": {
						"cn-beijing": "m-2zehhdtfg22hq46reabf",
						"cn-qingdao": "m-m5e1bdeo0a0ujw2s0ovv",
						"cn-zhangjiakou": "m-8vbdd1mi6hwlddicqt2x"
					},
					"infrastructure": "alicloud",
					"name": "bosh-alicloud-kvm-hubuntu-trusty-go_agent",
					"os_distro": "ubuntu",
					"os_type": "linux",
					"root_device_name": "/dev/vda1",
					"source_url": "",
					"version": 1003
				}
			],
			"context": {
				"director_uuid": "073eac6e-7a35-4a49-8c42-68988ea16ca7"
			}
		}`))
		Expect(r.GetError()).NotTo(HaveOccurred())
	})

	It("can import image from oss, and then create stemcell", func() {
		r := caller.Run([]byte(`{
			"method": "create_stemcell",
			"arguments": [
				"/var/vcap/data/tmp/director/stemcell20170926-9684-17ncxdz/image",
				{
					"architecture": "x86_64",
					"container_format": "",
					"disk": "50",
					"disk_format": "rawdisk",
					"hypervisor": "kvm",
					"format": "RAW",
					"oss_bucket": "bosh-alicloud",
					"oss_object": "alicloud-kvm-ubuntu-1008.raw",
					"infrastructure": "alicloud",
					"name": "bosh-alicloud-kvm-hubuntu-trusty-go_agent",
					"os_distro": "ubuntu",
					"os_type": "linux",
					"root_device_name": "/dev/vda1",
					"version": 1003
				}
			],
			"context": {
				"director_uuid": "073eac6e-7a35-4a49-8c42-68988ea16ca7"
			}
		}`))
		Expect(r.GetError()).NotTo(HaveOccurred())
	})

	It("can create stemcell from tarball", func() {
		r := caller.Run([]byte(`{
			"method": "create_stemcell",
			"arguments": [
				"/var/vcap/data/tmp/director/stemcell20170926-9684-17ncxdz/image",
				{
					"architecture": "x86_64",
					"container_format": "",
					"disk": "50",
					"disk_format": "rawdisk",
					"hypervisor": "kvm",
					"format": "RAW",
					"infrastructure": "alicloud",
					"name": "bosh-alicloud-kvm-hubuntu-trusty-go_agent",
					"os_distro": "ubuntu",
					"os_type": "linux",
					"root_device_name": "/dev/vda1",
					"version": 1003
				}
			],
			"context": {
				"director_uuid": "073eac6e-7a35-4a49-8c42-68988ea16ca7"
			}
		}`))
		Expect(r.GetError()).NotTo(HaveOccurred())
	})
})

var _ = Describe("create_stemcell EnableNvmeSupport", func() {
	// stemcell arguments shared across cases — uses the image_id (region map) path so
	// no real file I/O or OSS upload is needed.
	const stemcellArgs = `[
		"/var/vcap/data/tmp/director/stemcell/image",
		{
			"architecture": "x86_64",
			"disk": "50",
			"disk_format": "rawdisk",
			"hypervisor": "kvm",
			"image_id": {"cn-beijing": "m-2zehhdtfg22hq46reabf"},
			"infrastructure": "alicloud",
			"name": "bosh-alicloud-kvm-ubuntu-jammy-go_agent",
			"os_distro": "ubuntu",
			"os_type": "linux",
			"version": "1.0"
		}
	]`

	It("calls EnableNvmeSupport on the copied image when encryption is enabled", func() {
		cfg := mustParseEncryptedConfig()
		spy := &nvmeStemcellManagerSpy{
			StemcellManager: mock.NewStemcellManagerMock(mock.NewTestContext(cfg)),
		}
		encryptedCaller := newEncryptedCallerWithSpy(spy, cfg)

		r := encryptedCaller.Run([]byte(`{
			"method": "create_stemcell",
			"arguments": ` + stemcellArgs + `,
			"context": {"director_uuid": "073eac6e-7a35-4a49-8c42-68988ea16ca7"}
		}`))

		Expect(r.GetError()).NotTo(HaveOccurred())
		Expect(spy.nvmeCalledWith).To(HaveLen(1),
			"EnableNvmeSupport should be called exactly once on the copied image")
	})

	It("does not call EnableNvmeSupport when encryption is disabled", func() {
		spy := &nvmeStemcellManagerSpy{
			StemcellManager: mock.NewStemcellManagerMock(mockContext),
		}

		// Use the global caller (no encryption) but swap in our spy via a local caller.
		localCaller := newCallerWithSpy(spy)

		r := localCaller.Run([]byte(`{
			"method": "create_stemcell",
			"arguments": ` + stemcellArgs + `,
			"context": {"director_uuid": "073eac6e-7a35-4a49-8c42-68988ea16ca7"}
		}`))

		Expect(r.GetError()).NotTo(HaveOccurred())
		Expect(spy.nvmeCalledWith).To(BeEmpty(),
			"EnableNvmeSupport should not be called when encryption is disabled")
	})

	It("propagates an error from EnableNvmeSupport", func() {
		cfg := mustParseEncryptedConfig()
		spy := &nvmeStemcellManagerSpy{
			StemcellManager: mock.NewStemcellManagerMock(mock.NewTestContext(cfg)),
			nvmeErr:         fmt.Errorf("ModifyImageAttribute: forbidden"),
		}
		encryptedCaller := newEncryptedCallerWithSpy(spy, cfg)

		r := encryptedCaller.Run([]byte(`{
			"method": "create_stemcell",
			"arguments": ` + stemcellArgs + `,
			"context": {"director_uuid": "073eac6e-7a35-4a49-8c42-68988ea16ca7"}
		}`))

		Expect(r.GetError()).To(HaveOccurred())
		Expect(r.GetError().Error()).To(ContainSubstring("ModifyImageAttribute: forbidden"))
	})
})

// newCallerWithSpy builds a non-encrypted Caller with the given stemcell manager.
func newCallerWithSpy(spy *nvmeStemcellManagerSpy) Caller {
	logger := boshlog.NewWriterLogger(boshlog.LevelWarn, os.Stderr)
	config, err := alicloud.NewConfigFromBytes(configForTest)
	Expect(err).NotTo(HaveOccurred())
	mc := mock.NewTestContext(config)
	services := Services{
		Stemcells: spy,
		Osses:     mock.NewOssManagerMock(mc),
		Instances: mock.NewInstanceManagerMock(mc),
		Disks:     mock.NewDiskManagerMock(mc),
		Networks:  mock.NewNetworkManagerMock(mc),
		Registry:  mock.NewRegistryMock(),
	}
	return NewCallerWithServices(config, logger, services)
}

func mustParseEncryptedConfig() alicloud.Config {
	encryptedConfigBytes := []byte(`{
		"cloud": {
			"plugin": "alicloud",
			"properties": {
				"alicloud": {
					"region": "cn-beijing",
					"availability_zone": "cn-beijing-a",
					"access_key_id": "---",
					"access_key_secret": "---",
					"encrypted": true
				},
				"registry": {
					"user": "registry", "password": "registry",
					"protocol": "http", "address": "127.0.0.1", "port": 25777
				},
				"agent": {
					"ntp": ["0.pool.ntp.org"],
					"mbus": "http://mbus:mbus@0.0.0.0:6868",
					"blobstore": {
						"provider": "dav",
						"options": {"endpoint": "http://10.0.0.2:25250", "user": "agent", "password": "agent-password"}
					}
				}
			}
		}
	}`)
	config, err := alicloud.NewConfigFromBytes(encryptedConfigBytes)
	Expect(err).NotTo(HaveOccurred())
	Expect(config.Validate()).To(Succeed())
	return config
}
