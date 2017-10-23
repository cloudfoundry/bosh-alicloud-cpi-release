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
	"fmt"
	"bosh-alicloud-cpi/mock"
	"strings"
)

type CpiResponse struct {
	Result interface{}		`json:"result"`
	Error CpiError		 	`json:"error"`
	Log string				`json:"log"`
}

func WrapErrorResponse(err error, format string, args... interface{}) (CpiResponse) {
	return CpiResponse{
		Result: json.RawMessage{},
		Error:CpiError{
			"CpiError",
			err.Error(),
			false,
		},
		Log:fmt.Sprintf(format, args),
	}
}

func (r CpiResponse) GetError() error {
	return r.Error.ToError()
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

func NewTestCaller(config alicloud.Config, logger boshlog.Logger, mc mock.TestContext) (Caller) {
	services := Services {
		Stemcells: mock.NewStemcellManagerMock(mc),
		Instances: mock.NewInstanceManagerMock(mc),
		Disks: mock.NewDiskManagerMock(mc),
		Networks: mock.NewNetworkManagerMock(mc),
		Registry: mock.NewRegistryMock(),
	}
	return Caller {config, logger, services}
}

func (c Caller) Run(input []byte) (CpiResponse) {
	if !json.Valid(input) {
		err := fmt.Errorf("input json invalid %s", string(input))
		return WrapErrorResponse(err, "Run failed")
	}

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

func (c Caller) CallGeneric(method string, args ...interface{}) (interface{}, error) {
	arguments := ""
	for i, a := range args {
		if i > 0 {
			arguments += ","
		}
		switch a.(type) {
		case string:
			s := a.(string)
			if strings.HasPrefix(s, "{") && strings.HasSuffix(s, "}") {
				arguments += s
			} else {
				arguments += `"` + s + `"`
			}
		//case []interface{}:
		//	arguments += "["
		//	for j, s := range a.([]interface{}) {
		//		if j > 0 {
		//			arguments += ","
		//		}
		//		arguments += s.(string)
		//	}
		//	arguments += "]"
		default:
			j, _ := json.Marshal(a)
			arguments = arguments + string(j)
		}
	}

	in := fmt.Sprintf(`{
		"method": "%s",
		"arguments": [%s],
		"context": { "director_uuid": "%s" }
	}`, method, arguments, "911133bb-7d44-4811-bf8a-b215608bf084")

	r := c.Run([]byte(in))

	err := r.GetError()
	if err != nil {
		return "", err
	}
	return r.Result, nil
}

func (c Caller) Call(method string, args ...interface{}) (string, error) {
	r, err := c.CallGeneric(method, args...)
	if err != nil {
		return "", err
	}

	if r == nil {
		return "", err
	}

	s, ok := r.(string)
	if ok {
		return s, nil
	} else {
		return "", fmt.Errorf("result is not string %v", r)
	}
}
