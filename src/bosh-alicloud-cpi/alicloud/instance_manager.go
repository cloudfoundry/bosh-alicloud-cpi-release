/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package alicloud

import (
	"github.com/denverdino/aliyungo/ecs"
	"github.com/denverdino/aliyungo/common"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"time"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"strings"
	"encoding/json"
)

const (
	UseForceStop         = false
	ForceStopDelay       = 2			//
	DefaultTimeoutSecond = 120
	DefaultWaitSecond    = 2
)

type InstanceManager interface {
	GetInstance(cid string) (*ecs.InstanceAttributesType, error)

	CreateInstance(args ecs.CreateInstanceArgs) (string, error)
	DeleteInstance(cid string) (error)

	StartInstance(cid string) (error)
	StopInstance(cid string) (error)
	RebootInstance(cid string) (error)

	GetInstanceStatus(cid string) (ecs.InstanceStatus, error)
	WaitForInstanceStatus(cid string, toStatus ecs.InstanceStatus) (error)
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
		a.logger.Info("InstanceManager", "%s done args=%s result=%s", s, result)
	}
}

func (a InstanceManagerImpl) GetInstance(cid string) (*ecs.InstanceAttributesType, error) {
	client := a.config.NewEcsClient()

	var args ecs.DescribeInstancesArgs
	args.RegionId = common.Region(a.region)
	args.InstanceIds = "[\"" + cid + "\"]"

	insts, _, err := client.DescribeInstances(&args)

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
	cid, err := client.CreateInstance(&args)
	a.log("CreateInstance", err, args, cid)
	return cid, err
}

func (a InstanceManagerImpl) DeleteInstance(cid string) (error) {
	client := a.config.NewEcsClient()
	err := client.DeleteInstance(cid)
	a.log("DeleteInstance", err, cid, "ok")

	if err != nil {
		for i := 0; i < 10; i++ {
			if strings.Contains(err.Error(), "IncorrectInstanceStatus.Initializing") {
				time.Sleep(time.Duration(5) * time.Second)
				err := client.DeleteInstance(cid)
				a.logger.Error("DELETE", "DeleteInstance try %d time again got error %v", i, err)
				if err == nil {
					break
				}
			}
		}
	}
	return nil
}

func (a InstanceManagerImpl) StartInstance(cid string) error {
	client := a.config.NewEcsClient()
	err := client.StartInstance(cid)
	a.log("StartInstance", err, cid, "ok")
	return err
}

func (a InstanceManagerImpl) StopInstance(cid string) error {
	client := a.config.NewEcsClient()
	err := client.StopInstance(cid, UseForceStop)
	a.log("StopInstance", err, cid, "ok")
	if err != nil {
		return err
	}
	if UseForceStop {
		a.logger.Info("Instances", "when ForceStop sleep %d seconds...", ForceStopDelay)
		time.Sleep(time.Duration(ForceStopDelay) * time.Second)
		a.logger.Info("Instances", "when ForceStop sleep %d seconds done", ForceStopDelay)
	}
	return nil
}

func (a InstanceManagerImpl) RebootInstance(cid string) error {
	client := a.config.NewEcsClient()
	err := client.RebootInstance(cid, UseForceStop)
	a.log("RebootInstance", err, cid, "ok")
	return err
}

func (a InstanceManagerImpl) GetInstanceStatus(cid string) (ecs.InstanceStatus, error) {
	inst, err := a.GetInstance(cid)

	if err != nil {
		return "", err
	}

	if inst == nil {
		return ecs.Deleted, bosherr.Error("Missing Instance: id=" + cid)
	}

	return inst.Status, nil
}

func (a InstanceManagerImpl) WaitForInstanceStatus(cid string, toStatus ecs.InstanceStatus) (error) {
	timeout := DefaultTimeoutSecond
	for {
		status, err := a.GetInstanceStatus(cid)
		a.logger.Info("InstanceManager", "Waiting Instance %s from %s to %s", cid, status, toStatus)

		if err != nil {
			if toStatus == ecs.Deleted && status == ecs.Deleted {
				return nil
			}
			return err
		}

		if status == toStatus {
			return nil
		}

		if timeout > 0 {
			timeout -= DefaultWaitSecond
			time.Sleep(time.Duration(DefaultWaitSecond) * time.Second)
		} else {
			return bosherr.Error("WaitForInstanceStatus timeout")
		}
	}
}
