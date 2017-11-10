/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package alicloud

import (
	"github.com/denverdino/aliyungo/ecs"
	"github.com/denverdino/aliyungo/common"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"encoding/json"
	"time"
	"fmt"
	"github.com/google/uuid"
)

var DeleteInstanceCatcher = Catcher {"IncorrectInstanceStatus.Initializing", 10, 15}
var CreateInstanceCatcher = Catcher {"InvalidPrivateIpAddress.Duplicated", 10, 15}

const (
	ChangeInstanceStatusTimeout = time.Duration(300) * time.Second
	ChangeInstanceStatusSleepInterval = time.Duration(5) * time.Second
)

type InstanceManager interface {
	GetInstance(cid string) (*ecs.InstanceAttributesType, error)

	CreateInstance(args ecs.CreateInstanceArgs) (string, error)
	ModifyInstanceAttribute(cid string, name string, description string) (error)

	DeleteInstance(cid string) (error)

	StartInstance(cid string) (error)
	StopInstance(cid string) (error)
	RebootInstance(cid string) (error)

	GetInstanceStatus(cid string) (ecs.InstanceStatus, error)

	// WaitForInstanceStatus(cid string, toStatus ecs.InstanceStatus) (ecs.InstanceStatus, error)
	ChangeInstanceStatus(cid string, toStatus ecs.InstanceStatus, checkFunc func(status ecs.InstanceStatus) (bool, error)) (error)
}

type InstanceManagerImpl struct {
	config Config
	logger boshlog.Logger
	region string
}

func NewInstanceManager(config Config, logger boshlog.Logger) (InstanceManager) {
	return InstanceManagerImpl {
		config: config,
		logger: logger,
		region: config.OpenApi.RegionId,
	}
}

func (a InstanceManagerImpl) log(action string, err error, args interface{}, result string) {
	s, _ := json.Marshal(args)
	if err != nil {
		a.logger.Error("InstanceManager", "%s failed args=%s err=%s", action, s, err)
	} else {
		a.logger.Info("InstanceManager", "%s done! args=%s result=%s", action, s, result)
	}
}

func (a InstanceManagerImpl) GetInstance(cid string) (*ecs.InstanceAttributesType, error) {
	client := a.config.NewEcsClient()

	var args ecs.DescribeInstancesArgs
	args.RegionId = common.Region(a.region)
	args.InstanceIds = "[\"" + cid + "\"]"

	var insts []ecs.InstanceAttributesType
	invoker := NewInvoker()
	err := invoker.Run(func() (error) {
		r, _, e := client.DescribeInstances(&args)
		insts = r
		return e
	})

	if err != nil {
		return nil, err
	}

	if len(insts) == 0 {
		return nil, nil
	}

	return &insts[0], nil
}

func (a InstanceManagerImpl) CreateInstance(args ecs.CreateInstanceArgs) (string, error) {
	client := a.config.NewEcsClient()

	invoker := NewInvoker()
	invoker.AddCatcher(CreateInstanceCatcher)
	args.ClientToken = uuid.New().String()

	var cid string
	err := invoker.Run(func() (error) {
		a2 := args // copy args to avoid base64 again
		c2, e := client.CreateInstance(&a2)
		cid = c2
		a.log("CreateInstance", e, a2, c2)
		return e
	})
	return cid, err
}

func (a InstanceManagerImpl) ModifyInstanceAttribute(cid string, name string, description string) (error) {
	client := a.config.NewEcsClient()

	var args ecs.ModifyInstanceAttributeArgs
	args.InstanceId = cid
	args.InstanceName = name
	args.Description = description

	invoker := NewInvoker()
	invoker.AddCatcher(CreateInstanceCatcher)
	return invoker.Run(func() (error) {
		e := client.ModifyInstanceAttribute(&args)
		a.log("ModifyInstanceAttributes", e, args, "ok")
		return e
	})
}

func (a InstanceManagerImpl) DeleteInstance(cid string) (error) {
	client := a.config.NewEcsClient()

	invoker := NewInvoker()
	invoker.AddCatcher(DeleteInstanceCatcher)

	return invoker.Run(func() (error) {
		e := client.DeleteInstance(cid)
		a.log("DeleteInstance", e, cid, "ok")
		return e
	})
}

func (a InstanceManagerImpl) StartInstance(cid string) error {
	client := a.config.NewEcsClient()
	invoker := NewInvoker()
	return invoker.Run(func() (error) {
		err := client.StartInstance(cid)
		a.log("StartInstance", err, cid, "ok")
		return err
	})
}

func (a InstanceManagerImpl) StopInstance(cid string) error {
	client := a.config.NewEcsClient()
	invoker := NewInvoker()
	return invoker.Run(func() (error) {
		err := client.StopInstance(cid, UseForceStop)
		a.log("StopInstance", err, cid, "ok")
		return err
	})
}

func (a InstanceManagerImpl) RebootInstance(cid string) error {
	client := a.config.NewEcsClient()
	invoker := NewInvoker()
	return invoker.Run(func() (error) {
		err := client.RebootInstance(cid, UseForceStop)
		a.log("RebootInstance", err, cid, "ok")
		return err
	})
}

func (a InstanceManagerImpl) GetInstanceStatus(cid string) (ecs.InstanceStatus, error) {
	inst, err := a.GetInstance(cid)

	if err != nil {
		return "", err
	}

	if inst == nil {
		return ecs.Deleted, nil
	}
	return inst.Status, nil
}

func (a InstanceManagerImpl) WaitForInstanceStatus(cid string, toStatus ecs.InstanceStatus) (ecs.InstanceStatus ,error) {
	invoker := NewInvoker()

	var status ecs.InstanceStatus

	ok, err := invoker.RunUntil(WaitTimeout, WaitInterval, func() (bool, error) {
		status, e := a.GetInstanceStatus(cid)
		a.logger.Info("InstanceManager", "Waiting Instance %s from %s to %s", cid, status, toStatus)
		return status == toStatus, e
	})

	if err != nil {
		return status, err
	}

	if !ok {
		return status, bosherr.Errorf("WaitForInstance %s to %s timeout", cid, toStatus)
	}

	return status, nil
}

func (a InstanceManagerImpl) ChangeInstanceStatus(cid string, toStatus ecs.InstanceStatus, checkFunc func(status ecs.InstanceStatus) (bool, error)) (error) {
	timeout := ChangeInstanceStatusTimeout
	for {
		status, err := a.GetInstanceStatus(cid)
		if err != nil {
			return err
		}

		ok, err := checkFunc(status)

		if err != nil {
			a.logger.Error("InstanceManager", "change %s from %s to %s failed %s", cid, status, toStatus, err.Error())
			return err
		}

		if ok {
			a.logger.Info("InstanceManager", "change %s to %s done!", cid, toStatus)
			return nil
		} else {
			a.logger.Info("InstanceManager", "changing %s from %s to %s ...", cid, status, toStatus)
		}

		timeout -= ChangeInstanceStatusSleepInterval
		time.Sleep(ChangeInstanceStatusSleepInterval)
		if timeout < 0 {
			return fmt.Errorf("change instance %s to %s timeout", cid, toStatus)
		}
	}
}
