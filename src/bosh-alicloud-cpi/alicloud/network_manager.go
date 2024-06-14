/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package alicloud

import (
	"encoding/json"
	"fmt"
	"strings"

	util "github.com/alibabacloud-go/tea-utils/v2/service"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/slb"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

type NetworkManager interface {
	DescribeEip(region, eip string) (ecs.EipAddressInDescribeEipAddresses, error)
	BindEip(region, instanceId, eip string) error
	WaitForEipStatus(region, eip string, toStatus EipStatus) error

	BindSLB(region, instanceId, slbId string, weight int) error
	BindSlbServerGroup(region, instanceId, slbId string, weight int, port int) error
	BindNlbServerGroup(region, instanceId, nlbServerGroupId string, weight int, port int) error
	DescribeSecurityGroupAttribute(region, groupId string) (ecs.DescribeSecurityGroupAttributeResponse, error)
	JoinSecurityGroup(region, instanceId, groupId string) error
}

type NetworkManagerImpl struct {
	config Config
	logger boshlog.Logger
}

type BackendServerType struct {
	ServerId string
	Weight   int
}
type SlbServerGroupBackendServerType struct {
	ServerId string
	Weight   int
	Port     int
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

func (a NetworkManagerImpl) DescribeEip(region, eip string) (eipAddress ecs.EipAddressInDescribeEipAddresses, err error) {
	client, err := a.config.NewEcsClient(region)
	if err != nil {
		return
	}
	invoker := NewInvoker()

	args := ecs.CreateDescribeEipAddressesRequest()
	args.EipAddress = eip

	err = invoker.Run(func() error {
		r, err := client.DescribeEipAddresses(args)
		if err != nil {
			return err
		}
		if r == nil || len(r.EipAddresses.EipAddress) != 1 {
			return fmt.Errorf("expect 1 EIP(%s) but get %#v", eip, r)
		}
		eipAddress = r.EipAddresses.EipAddress[0]
		return nil
	})

	if err != nil {
		return eipAddress, bosherr.WrapErrorf(err, "DescribeEipAddress(%v) failed", args)
	}
	return
}

func (a NetworkManagerImpl) BindEip(region, instanceId, eip string) error {
	eipAddress, err := a.DescribeEip(region, eip)

	if err != nil {
		return bosherr.WrapErrorf(err, "DescribeEip(%s) failed", eip)
	}

	if eipAddress.Status != string(EipStatusAvailable) {
		return bosherr.WrapErrorf(err, "BindEip(%s) status %s failed", eip, eipAddress.Status)
	}

	client, err := a.config.NewEcsClient(region)
	if err != nil {
		return err
	}
	invoker := NewInvoker()

	args := ecs.CreateAssociateEipAddressRequest()
	args.AllocationId = eipAddress.AllocationId
	args.InstanceId = instanceId
	if strings.HasPrefix(instanceId, "i-") {
		args.InstanceType = "EcsInstance"
	}

	err = invoker.Run(func() error {
		_, e := client.AssociateEipAddress(args)
		return e
	})
	a.log("BingEip("+eip+")", err, instanceId, "")

	if err != nil {
		return bosherr.WrapErrorf(err, "AssociateEipAddress %s to %s failed", eip, instanceId)
	}

	err = a.WaitForEipStatus(region, eip, EipStatusInUse)
	if err != nil {
		return bosherr.WrapErrorf(err, "WaitForEipStatus %s to InUse failed", eip)
	}

	return nil
}

func (a NetworkManagerImpl) WaitForEipStatus(region, eip string, toStatus EipStatus) error {
	invoker := NewInvoker()
	ok, err := invoker.RunUntil(WaitTimeout, WaitInterval, func() (bool, error) {
		r, e := a.DescribeEip(region, eip)
		a.log("WaitForEipStatus", e, r.Status, "")
		return EipStatus(r.Status) == toStatus, e
	})

	if err != nil {
		return bosherr.WrapErrorf(err, "WaitForEip")
	}

	if !ok {
		return fmt.Errorf("WaitForEipStatus %s time out", eip)
	}
	return nil
}

func (a NetworkManagerImpl) BindSlbServerGroup(region, instanceId string, slbServerGroupId string, weight int, port int) error {
	client, err := a.config.NewSlbClient(region)
	if err != nil {
		return err
	}
	if weight == 0 {
		weight = DefaultSlbServerGroupWeight
	}
	bytes, _ := json.Marshal([]SlbServerGroupBackendServerType{
		{ServerId: instanceId, Weight: weight, Port: port},
	})
	args := slb.CreateAddVServerGroupBackendServersRequest()
	args.VServerGroupId = slbServerGroupId
	args.BackendServers = string(bytes)
	servers, err := client.AddVServerGroupBackendServers(args)
	if err != nil {
		a.logger.Error("NetworkManager", "BindSlbServerGroup %s to %s failed %v", instanceId, slbServerGroupId, err)
		return bosherr.WrapErrorf(err, "BindSlbServerGroup %s to %s failed", instanceId, slbServerGroupId)
	}

	a.logger.Info("NetworkManager", "BindSlbServerGroup %s to %s, after bind server=%v", instanceId, slbServerGroupId, servers)
	return err
}

func (a NetworkManagerImpl) BindNlbServerGroup(region, instanceId string, nlbServerGroupId string, weight int, port int) error {
	conn, err := a.config.NlbTeaClient(region)
	if err != nil {
		return err
	}
	request := map[string]interface{}{
		"Servers.1.Port":       port,
		"ServerGroupId":        nlbServerGroupId,
		"Servers.1.ServerId":   instanceId,
		"Servers.1.ServerType": "Ecs",
		"ClientToken":          buildClientToken("AddServersToServerGroup"),
	}
	if weight != 0 {
		request["Servers.1.Weight"] = weight
	}
	action := "AddServersToServerGroup"
	invoker := NewInvoker()
	invoker.AddCatcher(NlbBindServerCatcher_Conflict_Lock)
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	err = invoker.Run(func() error {
		_, e := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2022-04-30"), StringPointer("AK"), nil, request, &runtime)
		a.logger.Error("NetworkManager", "BindNlbServerGroup %s to %s failed %v. Retry...", instanceId, nlbServerGroupId, err)
		return e
	})
	return err
}

// TODO: add retry
func (a NetworkManagerImpl) BindSLB(region, instanceId string, slbId string, weight int) error {
	client, err := a.config.NewSlbClient(region)
	if err != nil {
		return err
	}
	if weight == 0 {
		weight = DefaultSlbWeight
	}
	bytes, _ := json.Marshal([]BackendServerType{
		{ServerId: instanceId, Weight: weight},
	})
	args := slb.CreateAddBackendServersRequest()
	args.LoadBalancerId = slbId
	args.BackendServers = string(bytes)

	servers, err := client.AddBackendServers(args)

	if err != nil {
		a.logger.Error("NetworkManager", "BindSLB %s to %s failed %v", instanceId, slbId, err)
		return bosherr.WrapErrorf(err, "BindSLB %s to %s failed", instanceId, slbId)
	}

	a.logger.Info("NetworkManager", "BindSLB %s to %s, after bind server=%v", instanceId, slbId, servers)
	return err
}

func (a NetworkManagerImpl) DescribeSecurityGroupAttribute(region, groupId string) (group ecs.DescribeSecurityGroupAttributeResponse, err error) {
	client, err := a.config.NewEcsClient(region)
	if err != nil {
		return
	}
	invoker := NewInvoker()

	args := ecs.CreateDescribeSecurityGroupAttributeRequest()
	args.SecurityGroupId = groupId

	err = invoker.Run(func() error {
		r, err := client.DescribeSecurityGroupAttribute(args)
		if r != nil {
			group = *r
		}
		return err
	})

	if err != nil {
		return group, bosherr.WrapErrorf(err, "DescribeSecurityGroupAttribute(%v) failed", args)
	}
	return
}

func (a NetworkManagerImpl) JoinSecurityGroup(region, instanceId, groupId string) error {
	if _, err := a.DescribeSecurityGroupAttribute(region, groupId); err != nil {
		return bosherr.WrapErrorf(err, "DescribeSecurityGroupAttribute(%s) failed", groupId)
	}

	client, err := a.config.NewEcsClient(region)
	if err != nil {
		return err
	}
	invoker := NewInvoker()
	invoker.AddCatcher(Catcher{"OperationConflict", 30, 1})

	args := ecs.CreateJoinSecurityGroupRequest()
	args.InstanceId = instanceId
	args.SecurityGroupId = groupId

	err = invoker.Run(func() error {
		_, err := client.JoinSecurityGroup(args)
		return err
	})
	a.log("JoinSecurityGroup("+groupId+")", err, instanceId, "")

	if err != nil {
		if e, ok := err.(*errors.ServerError); ok && e.ErrorCode() == "InvalidInstanceId.AlreadyExists" {
			return nil
		}
		return bosherr.WrapErrorf(err, "JoinSecurityGroup %s to %s failed, args: %#v.", instanceId, groupId, args)
	}

	return nil
}
