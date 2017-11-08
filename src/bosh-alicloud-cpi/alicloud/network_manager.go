/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package alicloud

import (
	"encoding/json"
	"fmt"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/denverdino/aliyungo/slb"
)

type NetworkManager interface {
	DescribeEip(eip string) (ecs.EipAddressSetType, error)
	BindEip(instanceId string, eip string) error
	WaitForEipStatus(eip string, toStatus ecs.EipStatus) error

	BindSLB(instanceId string, slbId string, weight int) error
}

type NetworkManagerImpl struct {
	config Config
	logger boshlog.Logger
}

func NewNetworkManager(config Config, logger boshlog.Logger) NetworkManager {
	return NetworkManagerImpl{
		config: config,
		logger: logger,
	}
}

func (a NetworkManagerImpl) log(action string, err error, args interface{}, result string) {
	s, _ := json.Marshal(args)
	if err != nil {
		a.logger.Error("NetworkManager", "%s failed args=%s err=%s", action, s, err)
	} else {
		a.logger.Info("NetworkManager", "%s done! args=%s result=%s", action, s, result)
	}
}

func (a NetworkManagerImpl) DescribeEip(eip string) (ecs.EipAddressSetType, error) {
	client := a.config.NewEcsClient()
	invoker := NewInvoker()

	var args ecs.DescribeEipAddressesArgs
	args.EipAddress = eip
	args.RegionId = a.config.OpenApi.GetRegion()

	var eipAddress ecs.EipAddressSetType
	err := invoker.Run(func() error {
		r, _, err := client.DescribeEipAddresses(&args)
		if err != nil {
			return err
		}
		if len(r) != 1 {
			return fmt.Errorf("expect 1 EIP(%s) but get %d", eip, len(r))
		}
		eipAddress = r[0]
		return nil
	})

	if err != nil {
		return eipAddress, bosherr.WrapErrorf(err, "DescribeEipAddress(%v) failed", args)
	}
	return eipAddress, nil
}

func (a NetworkManagerImpl) BindEip(instanceId string, eip string) error {
	eipAddress, err := a.DescribeEip(eip)

	if err != nil {
		return bosherr.WrapErrorf(err, "DescribeEip(%s) failed", eip)
	}

	if eipAddress.Status != ecs.EipStatusAvailable {
		return bosherr.WrapErrorf(err, "BindEip(%s) status %s failed", eip, eipAddress.Status)
	}

	client := a.config.NewEcsClient()
	invoker := NewInvoker()

	allocationId := eipAddress.AllocationId
	err = invoker.Run(func() error {
		return client.AssociateEipAddress(allocationId, instanceId)
	})
	a.log("BingEip("+eip+")", err, instanceId, "")

	if err != nil {
		return bosherr.WrapErrorf(err, "AssociateEipAddress %s to %s failed", eip, instanceId)
	}

	err = a.WaitForEipStatus(eip, ecs.EipStatusInUse)
	if err != nil {
		return bosherr.WrapErrorf(err, "WaitForEipStatus %s to InUse failed", eip)
	}

	return nil
}

func (a NetworkManagerImpl) WaitForEipStatus(eip string, toStatus ecs.EipStatus) error {
	invoker := NewInvoker()
	ok, err := invoker.RunUntil(WaitTimeout, WaitInterval, func() (bool, error) {
		r, e := a.DescribeEip(eip)
		a.log("WaitForEipStatus", e, r.Status, "")
		return r.Status == toStatus, e
	})

	if err != nil {
		return bosherr.WrapErrorf(err, "WaitForEip")
	}

	if !ok {
		return fmt.Errorf("WaitForEipStatus %s time out", eip)
	}
	return nil
}

func (a NetworkManagerImpl) BindSLB(instanceId string, slbId string, weight int) error {
	client := a.config.NewSlbClient()

	if weight == 0 {
		weight = DefaultSlbWeight
	}

	servers := []slb.BackendServerType{
		{ServerId: instanceId, Weight: weight},
	}

	servers, err := client.AddBackendServers(slbId, servers)

	if err != nil {
		a.logger.Error("NetworkManager", "BindSLB %s to %s failed %v", instanceId, slbId, err)
		return bosherr.WrapErrorf(err, "BindSLB %s to %s failed", instanceId, slbId)
	}

	a.logger.Info("NetworkManager", "BindSLB %s to %s, after bind server=%v", instanceId, slbId, servers)
	return err
}
