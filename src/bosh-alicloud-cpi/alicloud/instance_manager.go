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
)

const (
	UseForceStop = true
	ForceStopDelay = 2			//
	DefaultTimeoutMs = 120
	DefaultWaitInterval = 2
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
	return client.CreateInstance(&args)
}

func (a InstanceManagerImpl) DeleteInstance(cid string) (error) {
	client := a.config.NewEcsClient()
	err := client.DeleteInstance(cid)
	if err != nil {
		a.logger.Error("DELETE", "DeleteInstance got error %s", err.Error())
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
	return client.StartInstance(cid)
}

func (a InstanceManagerImpl) StopInstance(cid string) error {
	client := a.config.NewEcsClient()
	a.logger.Info("Instances", "Stopping vm %s", cid)
	err := client.StopInstance(cid, UseForceStop)
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
	return client.RebootInstance(cid, UseForceStop)
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
	timeout := DefaultTimeoutMs
	for {
		status, err := a.GetInstanceStatus(cid)
		a.logger.Info("WaitForInstance %s from %s to %s", cid, status, toStatus)

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
			timeout -= DefaultWaitInterval
			time.Sleep(time.Duration(DefaultWaitInterval) * time.Second)
		} else {
			return bosherr.Error("WaitForInstanceStatus timeout")
		}
	}
}
