/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package mock

import (
	"bosh-alicloud-cpi/alicloud"
)

func NewNetworkManagerMock(context TestContext) (alicloud.NetworkManager) {
	return NetworkManagerMock{}
}

type NetworkManagerMock struct {
}

func (a NetworkManagerMock) BindElasticIp(instanceId string, eip string) (error) {
	return nil
}

func (a NetworkManagerMock) BindSLB(instanceId string, slbId string, weight int) (error) {
	return nil
}