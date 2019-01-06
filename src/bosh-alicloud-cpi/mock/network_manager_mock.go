/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package mock

import (
	"bosh-alicloud-cpi/alicloud"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
)

func NewNetworkManagerMock(context TestContext) alicloud.NetworkManager {
	return NetworkManagerMock{}
}

type NetworkManagerMock struct {
}

func (a NetworkManagerMock) DescribeEip(eip string) (ecs.EipAddressInDescribeEipAddresses, error) {
	return ecs.EipAddressInDescribeEipAddresses{}, nil
}

func (a NetworkManagerMock) BindEip(instanceId string, eip string) error {
	return nil
}

func (a NetworkManagerMock) WaitForEipStatus(eip string, toStatus alicloud.EipStatus) error {
	return nil
}

func (a NetworkManagerMock) BindSLB(instanceId string, slbId string, weight int) error {
	return nil
}

func (a NetworkManagerMock) DescribeSecurityGroupAttribute(groupId string) (ecs.DescribeSecurityGroupAttributeResponse, error) {
	return ecs.DescribeSecurityGroupAttributeResponse{}, nil
}

func (a NetworkManagerMock) JoinSecurityGroup(instanceId string, groupId string) error {
	return nil
}
