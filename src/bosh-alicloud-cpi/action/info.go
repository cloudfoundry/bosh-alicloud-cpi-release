/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package action

import (
	"github.com/cppforlife/bosh-cpi-go/apiv1"
)

type InfoMethod struct{}

func NewInfoMethod() InfoMethod {
	return InfoMethod{}
}

func (a InfoMethod) Info() (apiv1.Info, error) {
	return apiv1.Info{StemcellFormats: []string{"general-tar", "alicloud-tar", "alicloud-raw"}}, nil
}
