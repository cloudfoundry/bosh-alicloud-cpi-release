/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package alicloud

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"

	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

var InstanceInvalidOperationConflictCatcher = Catcher{"InvalidOperation.Conflict", 20, 10}
var DeleteInstanceCatcher = Catcher{"IncorrectInstanceStatus.Initializing", 20, 15}
var CreateInstanceCatcher_IpUsed = Catcher{"InvalidPrivateIpAddress.Duplicated", 30, 10}
var CreateInstanceCatcher_IpUsed2 = Catcher{"InvalidIPAddress.AlreadyUsed", 30, 10}

const (
	ChangeInstanceStatusTimeout       = time.Duration(360) * time.Second
	ChangeInstanceStatusSleepInterval = time.Duration(5) * time.Second
)

type InstanceManager interface {
	GetInstance(cid string) (*ecs.DescribeInstanceAttributeResponse, error)

	CreateInstance(region string, args *ecs.CreateInstanceRequest) (string, error)
	ModifyInstanceAttribute(cid string, name string, description string) error
	AddTags(cid string, tags map[string]string) error

	DeleteInstance(cid string) error

	StartInstance(cid string) error
	StopInstance(cid string) error
	RebootInstance(cid string) error

	GetInstanceStatus(cid string) (InstanceStatus, error)

	// WaitForInstanceStatus(cid string, toStatus ecs.InstanceStatus) (ecs.InstanceStatus, error)
	ChangeInstanceStatus(cid string, toStatus InstanceStatus, checkFunc func(status InstanceStatus) (bool, error)) error
}

type InstanceManagerImpl struct {
	config Config
	logger boshlog.Logger
	region string
}

func NewInstanceManager(config Config, logger boshlog.Logger) InstanceManager {
	return InstanceManagerImpl{
		config: config,
		logger: logger,
		region: config.OpenApi.GetRegion(""),
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

func (a InstanceManagerImpl) GetInstance(cid string) (inst *ecs.DescribeInstanceAttributeResponse, err error) {
	client, err := a.config.NewEcsClient("")
	if err != nil {
		return nil, err
	}

	args := ecs.CreateDescribeInstanceAttributeRequest()
	args.InstanceId = cid

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		r, e := client.DescribeInstanceAttribute(args)
		if e != nil && IsExceptedErrors(e, EcsInstanceNotFound) {
			return nil
		}
		inst = r
		return e
	})

	return
}

func (a InstanceManagerImpl) CreateInstance(region string, args *ecs.CreateInstanceRequest) (string, error) {
	client, err := a.config.NewEcsClient(region)
	if err != nil {
		return "", err
	}

	invoker := NewInvoker()
	invoker.AddCatcher(CreateInstanceCatcher_IpUsed)
	invoker.AddCatcher(CreateInstanceCatcher_IpUsed2)

	args.ClientToken = buildClientToken(args.GetActionName())

	var cid string
	err = invoker.Run(func() error {
		resp, e := client.CreateInstance(args)
		if resp != nil {
			cid = resp.InstanceId
		}
		return e
	})
	return cid, err
}

func (a InstanceManagerImpl) ModifyInstanceAttribute(cid, name, description string) error {
	client, err := a.config.NewEcsClient("")
	if err != nil {
		return err
	}

	args := ecs.CreateModifyInstanceAttributeRequest()
	args.InstanceId = cid
	args.InstanceName = name
	args.Description = description

	invoker := NewInvoker()
	return invoker.Run(func() error {
		_, e := client.ModifyInstanceAttribute(args)
		a.log("ModifyInstanceAttribute", e, args, "ok")
		return e
	})
}

func (a InstanceManagerImpl) DeleteInstance(cid string) error {
	client, err := a.config.NewEcsClient("")
	if err != nil {
		return err
	}

	invoker := NewInvoker()
	invoker.AddCatcher(DeleteInstanceCatcher)
	invoker.AddCatcher(InstanceInvalidOperationConflictCatcher)

	args := ecs.CreateDeleteInstanceRequest()
	args.InstanceId = cid

	return invoker.Run(func() error {
		_, e := client.DeleteInstance(args)
		a.log("DeleteInstance", e, cid, "ok")
		if e != nil && IsExceptedErrors(e, EcsInstanceNotFound) {
			return nil
		}
		return e
	})
}

func (a InstanceManagerImpl) StartInstance(cid string) error {
	client, err := a.config.NewEcsClient("")
	if err != nil {
		return err
	}

	invoker := NewInvoker()
	invoker.AddCatcher(InstanceInvalidOperationConflictCatcher)

	args := ecs.CreateStartInstanceRequest()
	args.InstanceId = cid

	return invoker.Run(func() error {
		_, err := client.StartInstance(args)
		a.log("StartInstance", err, cid, "ok")
		return err
	})
}

func (a InstanceManagerImpl) StopInstance(cid string) error {
	client, err := a.config.NewEcsClient("")
	if err != nil {
		return err
	}

	invoker := NewInvoker()
	invoker.AddCatcher(InstanceInvalidOperationConflictCatcher)

	args := ecs.CreateStopInstanceRequest()
	args.InstanceId = cid
	args.ForceStop = requests.NewBoolean(UseForceStop)

	return invoker.Run(func() error {
		_, err := client.StopInstance(args)
		a.log("StopInstance", err, cid, "ok")
		return err
	})
}

func (a InstanceManagerImpl) RebootInstance(cid string) error {
	client, err := a.config.NewEcsClient("")
	if err != nil {
		return err
	}

	invoker := NewInvoker()
	invoker.AddCatcher(InstanceInvalidOperationConflictCatcher)

	args := ecs.CreateRebootInstanceRequest()
	args.InstanceId = cid
	args.ForceStop = requests.NewBoolean(UseForceStop)

	return invoker.Run(func() error {
		_, err := client.RebootInstance(args)
		a.log("RebootInstance", err, cid, "ok")
		return err
	})
}

func (a InstanceManagerImpl) AddTags(cid string, tags map[string]string) error {
	client, err := a.config.NewEcsClient("")
	if err != nil {
		return err
	}

	invoker := NewInvoker()
	return invoker.Run(func() error {
		_, err := client.AddTags(getTagsRequest(cid, tags))
		if err != nil && IsExceptedErrors(err, ResourceNotFound) {
			regions, e := a.config.GetCrossRegions()
			if e != nil {
				return e
			}
			if len(regions) > 0 {
				for _, r := range regions {
					client, e := a.config.NewEcsClient(r)
					if e != nil {
						return e
					}
					if _, e = client.AddTags(getTagsRequest(cid, tags)); e != nil && IsExceptedErrors(e, ResourceNotFound) {
						continue
					}
					return e
				}
			}
		}
		return err
	})
}

func (a InstanceManagerImpl) GetInstanceStatus(cid string) (InstanceStatus, error) {
	inst, err := a.GetInstance(cid)

	if err != nil {
		return "", err
	}

	if inst == nil {
		return Deleted, nil
	}
	return InstanceStatus(inst.Status), nil
}

func (a InstanceManagerImpl) WaitForInstanceStatus(cid string, toStatus InstanceStatus) (InstanceStatus, error) {
	invoker := NewInvoker()

	var status InstanceStatus

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

func (a InstanceManagerImpl) ChangeInstanceStatus(cid string, toStatus InstanceStatus, checkFunc func(status InstanceStatus) (bool, error)) error {
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

func (a InstanceManagerImpl) GetInstanceUserData() {
	//client := a.config.NewEcsClient()
	//client.AddTags()
	//client.RemoveTags()
	//client.
	// inst, err := client.DescribeUserdata()
	// inst.
}

func getTagsRequest(cid string, tags map[string]string) *ecs.AddTagsRequest {
	args := ecs.CreateAddTagsRequest()
	args.ResourceId = cid
	var argsTags []ecs.AddTagsTag

	for k, v := range tags {
		argsTags = append(argsTags, ecs.AddTagsTag{
			Key:   k,
			Value: v,
		})
	}
	args.Tag = &argsTags

	if strings.HasPrefix(cid, "i-") {
		args.ResourceType = string(TagResourceInstance)
	} else if strings.HasPrefix(cid, "d-") {
		args.ResourceType = string(TagResourceDisk)
	}
	return args

}
