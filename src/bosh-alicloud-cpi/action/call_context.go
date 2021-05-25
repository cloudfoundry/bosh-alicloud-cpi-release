/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package action

import (
	"bosh-alicloud-cpi/alicloud"
	//"bosh-alicloud-cpi/registry"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

type CallContext struct {
	Input         string
	Logger        boshlog.Logger
	Config        alicloud.Config
	ConfigConcern bool
}

type Services struct {
	Stemcells alicloud.StemcellManager
	Osses     alicloud.OssManager
	Instances alicloud.InstanceManager
	Disks     alicloud.DiskManager
	Networks  alicloud.NetworkManager
	Registry  alicloud.RegistryManager
}

func NewCallContext(input []byte, logger boshlog.Logger, config alicloud.Config) CallContext {
	return CallContext{
		Input:         string(input),
		Logger:        logger,
		Config:        config,
		ConfigConcern: false,
	}
}
