/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package alicloud

import (
	"github.com/denverdino/aliyungo/ecs"
	"github.com/denverdino/aliyungo/common"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"time"
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
	region string
}

func NewInstanceManager(config Config) (InstanceManager) {
	return InstanceManagerImpl {
		config: config,
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
	return "", nil
}

func (a InstanceManagerImpl) DeleteInstance(cid string) (error) {
	return nil
}

func (a InstanceManagerImpl) StartInstance(cid string) error {
	client := a.config.NewEcsClient()
	return client.StartInstance(cid)
}

func (a InstanceManagerImpl) StopInstance(cid string) error {
	client := a.config.NewEcsClient()
	return client.StopInstance(cid, USE_FORCE_STOP)
}

func (a InstanceManagerImpl) RebootInstance(cid string) error {
	client := a.config.NewEcsClient()
	return client.RebootInstance(cid, USE_FORCE_STOP)
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
	timeout := DEFAULT_TIMEOUT
	for {
		status, err := a.GetInstanceStatus(cid)

		if err != nil {
			if status == ecs.Deleted && toStatus == ecs.Deleted {
				return nil
			}

			return err
		}

		if status == toStatus {
			return nil
		}

		if timeout > 0 {
			timeout -= DEFAULT_WAIT_INTERVAL
			time.Sleep(time.Duration(DEFAULT_WAIT_INTERVAL) * time.Millisecond)
		} else {
			return bosherr.Error("WaitForInstanceStatus timeout")
		}
	}
}
