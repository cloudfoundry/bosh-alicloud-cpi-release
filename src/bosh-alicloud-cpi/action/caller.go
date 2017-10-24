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
	"strings"
)

type CpiResponse struct {
	Result interface{}		`json:"result"`
	Error interface{}		 `json:"error,omitempty"`
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
	if r.Error == nil {
		return nil
	} else {
		e := r.Error.(CpiError)
		return e.ToError()
	}
}

func (r CpiResponse) GetResultString() (string) {
	return r.Result.(string)
}

func (r CpiResponse) GetResult() interface{} {
	return r.Result
}

type CpiError struct {
	Type string			`json:"type,omitempty"`
	Message string		`json:"message"`
	OkToRetry bool		`json:"ok_to_retry"`
}

func (e CpiError) ToError() error {
	if e.OkToRetry {
		return fmt.Errorf("CpiError[%s,ok_to_retry] %s", e.Type, e.Message)
	} else {
		return fmt.Errorf("CpiError[%s] %s", e.Type, e.Message)
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
		Instances: alicloud.NewInstanceManager(config, logger),
		Disks: alicloud.NewDiskManager(config),
		Networks: alicloud.NewNetworkManager(config),
		Registry: config.GetHttpRegistryClient(logger),
	}
	return NewCallerWithServices(config, logger, services)
}

func NewCallerWithServices(config alicloud.Config, logger boshlog.Logger, services Services) (Caller) {
	return Caller {config, logger, services}
}

func (c Caller) Run(input []byte) (CpiResponse) {
	// json.Validate not support with golang 1.8.1
	//if !json.Validate(input) {
	//	err := fmt.Errorf("input json invalid %s", string(input))
	//	return WrapErrorResponse(err, "Run failed")
	//}

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
