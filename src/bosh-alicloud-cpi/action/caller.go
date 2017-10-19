/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package action

import (
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshrpc "github.com/cppforlife/bosh-cpi-go/rpc"

	"bosh-alicloud-cpi/alicloud"
	"bytes"
	"encoding/json"
	"bosh-alicloud-cpi/registry"
	"fmt"
)

type CpiResponse struct {
	Result string		`json:"result"`
	Error CpiError	 	`json:"error"`
	Log string			`json:"log"`
}

func WrapErrorResponse(err error, msg string, args... interface{}) (CpiResponse) {
	return CpiResponse{}	//TODO
}

func (c CpiResponse) GetError() error {
	return c.Error.ToError()
}

type CpiError struct {
	Type string			`json:"type"`
	Message string		`json:"message"`
	OkToRetry bool		`json:"ok_to_retry"`
}

func (e CpiError) ToError() error {
	if e.Type != "" {
		return fmt.Errorf("%s %s retry=%v", e.Type, e.Message, e.OkToRetry)
	} else {
		return nil
	}
}

type Caller struct {
	Config alicloud.Config
	Logger boshlog.Logger
	Services
}

func NewCaller(config alicloud.Config, logger boshlog.Logger) (Caller) {
	services := Services {
		Stemcells: alicloud.NewStemcellManager(config),
		Instances: alicloud.NewInstanceManager(config),
		Disks: alicloud.NewDiskManager(config),
		Networks: alicloud.NewNetworkManager(config),
		Registry: config.GetHttpRegistryClient(logger),
	}
	return Caller {config, logger,services}
}

func NewTestCaller(config alicloud.Config, logger boshlog.Logger) (Caller) {
	services := Services {
		Stemcells: alicloud.NewStemcellManagerMock(),
		Instances: alicloud.NewInstanceManagerMock(),
		Disks: alicloud.NewDiskManagerMock(),
		Networks: alicloud.NewNetworkManagerMock(),
		Registry: registry.NewClientMock(),
	}
	return Caller {config, logger, services}
}

func (c Caller)RunTest(input []byte) (CpiResponse) {
//	logger := boshlog.NewWriterLogger(boshlog.LevelDebug, os.Stderr)

	reader := bytes.NewReader(input)
	output := new(bytes.Buffer)

	cc := NewCallContext(input, c.Logger, c.Config)

	cpiFactory := NewFactory(cc, c.Services)
	cli := boshrpc.NewFactory(c.Logger).NewCLIWithInOut(reader, output, cpiFactory)
	err := cli.ServeOnce()

	if err != nil {
		return WrapErrorResponse(err, "ServeOnce() Failed")
	}

	var resp CpiResponse
	err = json.Unmarshal(output.Bytes(), &resp)

	if err != nil {
		return WrapErrorResponse(err, "ServeOnce() result unmarshal failed %s", output.Bytes())
	}

	return resp
}

func (c Caller)Run(input []byte) (CpiResponse) {
	//	logger := boshlog.NewWriterLogger(boshlog.LevelDebug, os.Stderr)

	reader := bytes.NewReader(input)
	output := new(bytes.Buffer)

	cc := NewCallContext(input, c.Logger, c.Config)

	cpiFactory := NewFactory(cc, c.Services)
	cli := boshrpc.NewFactory(c.Logger).NewCLIWithInOut(reader, output, cpiFactory)
	err := cli.ServeOnce()

	if err != nil {
		return WrapErrorResponse(err, "ServeOnce() Failed")
	}

	var resp CpiResponse
	err = json.Unmarshal(output.Bytes(), &resp)

	if err != nil {
		return WrapErrorResponse(err, "ServeOnce() result unmarshal failed %s", output.Bytes())
	}

	return resp
}
