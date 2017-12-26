/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package mock

import (
	"bosh-alicloud-cpi/alicloud"
	"github.com/denverdino/aliyungo/ecs"
)

func NewNetworkManagerMock(context TestContext) (alicloud.NetworkManager) {
	return NetworkManagerMock{}
}

type NetworkManagerMock struct {
}

func (a NetworkManagerMock) DescribeEip(eip string) (ecs.EipAddressSetType, error) {
	return ecs.EipAddressSetType{}, nil
}

func (a NetworkManagerMock) BindEip(instanceId string, eip string) (error) {
	return nil
}

func (a NetworkManagerMock) WaitForEipStatus(eip string, toStatus ecs.EipStatus) (error) {
	return nil
}

func (a NetworkManagerMock) BindSLB(instanceId string, slbId string, weight int) error {
	return nil
}

func (a NetworkManagerMock) DescribeSecurityGroupAttribute(groupId string) (ecs.DescribeSecurityGroupAttributeResponse, error){
	return ecs.DescribeSecurityGroupAttributeResponse{}, nil
}

func (a NetworkManagerMock) JoinSecurityGroup(instanceId string, groupId string) error {
	return nil
}