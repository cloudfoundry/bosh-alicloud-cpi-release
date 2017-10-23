/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package action

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"bosh-alicloud-cpi/alicloud"
	"bosh-alicloud-cpi/registry"
)

type CallContext struct {
	Input string
	Logger boshlog.Logger
	Config alicloud.Config
	ConfigConcern bool
}

type Services struct {
	Stemcells alicloud.StemcellManager
	Instances alicloud.InstanceManager
	Disks alicloud.DiskManager
	Networks alicloud.NetworkManager
	Registry registry.Client
}

func NewCallContext(input []byte, logger boshlog.Logger, config alicloud.Config) (CallContext) {
	return CallContext{
		Input: string(input),
		Logger: logger,
		Config: config,
		ConfigConcern: false,
	}
}

func (c CallContext) WrapError(err error, msg string) (error) {
	// TODO add input to output
	return bosherr.WrapError(err, msg)
}

func (c CallContext) WrapErrorf(err error, fmt string, args... interface{}) (error) {
	// TODO add input json to error output
	return bosherr.WrapErrorf(err, fmt, args)
}