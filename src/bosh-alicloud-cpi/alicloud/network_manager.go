/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package alicloud

import (
	"github.com/denverdino/aliyungo/ecs"
	"fmt"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/denverdino/aliyungo/slb"
)

type NetworkManager interface {
	BindElasticIp(instanceId string, eip string) (error)
	BindSLB(instanceId string, slbId string, weight int) (error)
}

type NetworkManagerImpl struct {
	config Config
	logger boshlog.Logger
}

func NewNetworkManager(config Config, logger boshlog.Logger) (NetworkManager) {
	// config.NewSlbClient().AddBackendServers()
	return NetworkManagerImpl {
		config: config,
		logger: logger,
	}
}

func (a NetworkManagerImpl) BindElasticIp(instanceId string, eip string) (error) {
	if eip == "" {
		return fmt.Errorf("BindElasticIp input can't be empty")
	}

	client := a.config.NewEcsClient()

	var arg ecs.DescribeEipAddressesArgs
	arg.EipAddress = eip
	arg.RegionId = a.config.OpenApi.Region

	eipAddresses, _, err := client.DescribeEipAddresses(&arg)

	if err != nil {
		return bosherr.WrapErrorf(err, "DescribeEipAddress failed")
	}

	var allocationId string
	for _, address := range eipAddresses {
		allocationId = address.AllocationId
	}

	if err := client.AssociateEipAddress(allocationId, instanceId); err != nil {
		return bosherr.WrapErrorf(err, "AssociateEipAddress %s to %s failed", eip, instanceId)
	}

	if err := client.WaitForEip(arg.RegionId, allocationId, ecs.EipStatusInUse, DefaultEipWaitSeconds); err != nil {
		return bosherr.WrapErrorf(err, "Wait Eip status to EipStatusInUse failed")
	}
	return nil
}

func (a NetworkManagerImpl) BindSLB(instanceId string, slbId string, weight int) (error) {
	client := a.config.NewSlbClient()

	if weight == 0 {
		weight = DefaultSlbWeight
	}

	servers := []slb.BackendServerType{
		{ServerId:instanceId, Weight: weight},
	}

	servers, err := client.AddBackendServers(slbId, servers)

	if err != nil {
		a.logger.Error("NetworkManager", "BindSLB %s to %s failed %v", instanceId, slbId, err)
		return bosherr.WrapErrorf(err, "BindSLB %s to %s failed", instanceId, slbId)
	}

	a.logger.Info("NetworkManager", "BindSLB %s to %s, after bind server=%v", instanceId, slbId, servers)
	return err
}
