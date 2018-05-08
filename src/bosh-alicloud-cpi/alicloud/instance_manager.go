/*
 * Copyright (C) 2017-2018 Alibaba Group Holding Limited
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
	"github.com/google/uuid"
)

var DeleteInstanceCatcher = Catcher{"IncorrectInstanceStatus.Initializing", 10, 15}
var CreateInstanceCatcher_IpUsed = Catcher{"InvalidPrivateIpAddress.Duplicated", 10, 15}
var CreateInstanceCatcher_IpUsed2 = Catcher{"InvalidIPAddress.AlreadyUsed", 10, 15}
var CreateInstanceCatcher_TokenProcessing = Catcher{"LastTokenProcessing", 10, 15}

const (
	ChangeInstanceStatusTimeout       = time.Duration(300) * time.Second
	ChangeInstanceStatusSleepInterval = time.Duration(5) * time.Second
)

type InstanceManager interface {
	GetInstance(cid string) (*ecs.Instance, error)

	CreateInstance(args *ecs.CreateInstanceRequest) (string, error)
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

func (a InstanceManagerImpl) GetInstance(cid string) (*ecs.Instance, error) {
	client, err := a.config.NewEcsClient()
	if err != nil {
		return nil, err
	}

	args := ecs.CreateDescribeInstancesRequest()
	args.RegionId = a.region
	args.InstanceIds = "[\"" + cid + "\"]"

	var insts *ecs.DescribeInstancesResponse
	invoker := NewInvoker()
	err = invoker.Run(func() error {
		r, e := client.DescribeInstances(args)
		insts = r
		return e
	})

	if err != nil {
		return nil, err
	}

	if insts == nil || len(insts.Instances.Instance) <= 0 {
		return nil, nil
	}

	return &insts.Instances.Instance[0], nil
}

func (a InstanceManagerImpl) CreateInstance(args *ecs.CreateInstanceRequest) (string, error) {
	client, err := a.config.NewEcsClient()
	if err != nil {
		return "", err
	}

	invoker := NewInvoker()
	invoker.AddCatcher(CreateInstanceCatcher_IpUsed)
	invoker.AddCatcher(CreateInstanceCatcher_IpUsed2)
	invoker.AddCatcher(CreateInstanceCatcher_TokenProcessing)

	args.RegionId = a.config.OpenApi.GetRegion()
	token := strings.Replace(fmt.Sprintf("bosh-cpi-%s-%s", time.Now().String(), uuid.New().String()), " ", "", -1)
	args.ClientToken = token
	if len(token) > 64 {
		args.ClientToken = token[0:64]
	}

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

func (a InstanceManagerImpl) ModifyInstanceAttribute(cid string, name string, description string) error {
	client, err := a.config.NewEcsClient()
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
		a.log("ModifyInstanceAttributes", e, args, "ok")
		return e
	})
}

func (a InstanceManagerImpl) DeleteInstance(cid string) error {
	client, err := a.config.NewEcsClient()
	if err != nil {
		return err
	}

	invoker := NewInvoker()
	invoker.AddCatcher(DeleteInstanceCatcher)

	args := ecs.CreateDeleteInstanceRequest()
	args.InstanceId = cid

	return invoker.Run(func() error {
		_, e := client.DeleteInstance(args)
		a.log("DeleteInstance", e, cid, "ok")
		return e
	})
}

func (a InstanceManagerImpl) StartInstance(cid string) error {
	client, err := a.config.NewEcsClient()
	if err != nil {
		return err
	}

	invoker := NewInvoker()

	args := ecs.CreateStartInstanceRequest()
	args.InstanceId = cid

	return invoker.Run(func() error {
		_, err := client.StartInstance(args)
		a.log("StartInstance", err, cid, "ok")
		return err
	})
}

func (a InstanceManagerImpl) StopInstance(cid string) error {
	client, err := a.config.NewEcsClient()
	if err != nil {
		return err
	}

	invoker := NewInvoker()

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
	client, err := a.config.NewEcsClient()
	if err != nil {
		return err
	}

	invoker := NewInvoker()

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
	client, err := a.config.NewEcsClient()
	if err != nil {
		return err
	}

	args := ecs.CreateAddTagsRequest()
	args.RegionId = a.config.OpenApi.GetRegion()
	args.ResourceId = cid
	count := 1
	for k, v := range tags {
		switch count {

		case 1:
			args.Tag1Key = k
			args.Tag1Value = v
		case 2:
			args.Tag2Key = k
			args.Tag2Value = v
		case 3:
			args.Tag3Key = k
			args.Tag3Value = v
		case 4:
			args.Tag4Key = k
			args.Tag4Value = v
		case 5:
			args.Tag4Key = k
			args.Tag4Value = v
		default:
			break
		}
		count++
	}

	if strings.HasPrefix(cid, "i-") {
		args.ResourceType = string(TagResourceInstance)
	} else if strings.HasPrefix(cid, "d-") {
		args.ResourceType = string(TagResourceDisk)
	} else {
		return fmt.Errorf("unexpect resource type id=%s", cid)
	}

	invoker := NewInvoker()
	return invoker.Run(func() error {
		_, err := client.AddTags(args)
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
