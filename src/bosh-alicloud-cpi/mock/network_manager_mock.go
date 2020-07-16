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

func (a NetworkManagerMock) DescribeEip(region, eip string) (ecs.EipAddressInDescribeEipAddresses, error) {
	return ecs.EipAddressInDescribeEipAddresses{}, nil
}

func (a NetworkManagerMock) BindEip(region string, instanceId string, eip string) error {
	return nil
}

func (a NetworkManagerMock) WaitForEipStatus(region, eip string, toStatus alicloud.EipStatus) error {
	return nil
}

func (a NetworkManagerMock) BindSlbServerGroup(region, instanceId, slbId string, weight int, port int) error {
	return nil
}

func (a NetworkManagerMock) BindSLB(region, instanceId, slbId string, weight int) error {
	return nil
}

func (a NetworkManagerMock) DescribeSecurityGroupAttribute(region, groupId string) (ecs.DescribeSecurityGroupAttributeResponse, error) {
	return ecs.DescribeSecurityGroupAttributeResponse{}, nil
}

func (a NetworkManagerMock) JoinSecurityGroup(region, instanceId string, groupId string) error {
	return nil
}
