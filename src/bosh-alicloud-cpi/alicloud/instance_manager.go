/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package alicloud

import (
	"github.com/denverdino/aliyungo/ecs"
	"github.com/denverdino/aliyungo/common"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"time"
	"strings"
	"encoding/json"
	"strconv"
)

const (
	UseForceStop			= true
	ForceStopWaitSeconds	= 90
	DefaultWaitTimeout		= time.Duration(120) * time.Second
	DefaultWaitInterval		= time.Duration(2) * time.Second

	DeleteInstanceRetryCount	= 10
	DeleteInstanceRetryReason	= "IncorrectInstanceStatus.Initializing"
	DeleteInstanceRetryInterval	= time.Duration(5) * time.Second
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
		//
		// retry
		for i := 0; i < DeleteInstanceRetryCount; i++ {
			if strings.Contains(err.Error(), DeleteInstanceRetryReason) {
				time.Sleep(DeleteInstanceRetryInterval)
				err := client.DeleteInstance(cid)
				a.log("DeleteInstance retry:" + strconv.Itoa(i), err, cid, "ok")
				if err == nil {
					break
				}
			} else {
				return err
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
	if !UseForceStop {
		a.log("StopInstance", err, cid, "ok")
		return err
	} else {
		//
		// if use force stop, some ECS resource is not released,
		// for run DeleteInstance, need wait for a while
		if err != nil {
			a.log("StopInstance(Force)", err, cid, "ok")
			return err
		} else {
			a.logger.Info("InstanceManager", "StopInstance(Force) %s done, waiting for %d seconds...", cid, ForceStopWaitSeconds)
			time.Sleep(time.Duration(ForceStopWaitSeconds) * time.Second)
			a.logger.Info("InstanceManager", "StopInstance(Force) wait done.")
			return nil
		}
	}
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
	timeout := DefaultWaitTimeout
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
			timeout -= DefaultWaitInterval
			time.Sleep(DefaultWaitInterval)
		} else {
			return bosherr.Error("WaitForInstanceStatus timeout")
		}
	}
}
