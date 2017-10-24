/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package action

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"bosh-alicloud-cpi/alicloud"
	"bosh-alicloud-cpi/registry"
	"fmt"
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

func (c CallContext) Errorf(msg string, args... interface{}) (error) {
	s := "input=`" + c.Input + "` " + fmt.Sprintf(msg, args...)
	return bosherr.Error(s)
}

func (c CallContext) WrapError(err error, msg string) (error) {
	return bosherr.WrapErrorf(err, "input=`%s` message=%s", c.Input, msg)
}

func (c CallContext) WrapErrorf(err error, msg string, args... interface{}) (error) {
	s := "input=`" + c.Input + "` " + fmt.Sprintf(msg, args...)
	return bosherr.WrapError(err, s)
}