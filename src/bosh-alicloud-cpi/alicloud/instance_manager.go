/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package alicloud

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"

	util "github.com/alibabacloud-go/tea-utils/v2/service"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

var CreateInstanceCatcher_TokenProcessing = Catcher{"LastTokenProcessing", 15, 5}
var CreateInstanceCatcher_IdempotentProcessing = Catcher{"IdempotentProcessing", 15, 5}
var InstanceInvalidOperationConflictCatcher = Catcher{"InvalidOperation.Conflict", 20, 10}
var DeleteInstanceCatcher = Catcher{"IncorrectInstanceStatus.Initializing", 20, 15}
var IncorrectInstanceStatusCatcher = Catcher{"IncorrectInstanceStatus", 30, 10}
var CreateInstanceCatcher_IpUsed = Catcher{"InvalidPrivateIpAddress.Duplicated", 30, 10}
var CreateInstanceCatcher_IpUsed2 = Catcher{"InvalidIPAddress.AlreadyUsed", 30, 10}
var NetworkInterfaceInvalidOperationInvalidEniStateCacher = Catcher{"InvalidOperation.InvalidEniState", 60, 5}
var NlbBindServerCatcher_Conflict_Lock = Catcher{"Conflict.Lock", 15, 3}

const (
	ChangeInstanceStatusTimeout       = time.Duration(600) * time.Second
	ChangeInstanceStatusSleepInterval = time.Duration(5) * time.Second
)

type InstanceManager interface {
	GetInstance(cid string) (*ecs.Instance, error)

	CreateInstance(region string, request map[string]interface{}) (string, error)
	ModifyInstanceAttribute(cid string, name string, description string) error
	AddTags(cid string, tags map[string]string) error

	DeleteInstance(cid string) error

	StartInstance(cid string) error
	StopInstance(cid string) error
	RebootInstance(cid string) error

	GetInstanceStatus(cid string) (InstanceStatus, error)

	// WaitForInstanceStatus(cid string, toStatus ecs.InstanceStatus) (ecs.InstanceStatus, error)
	ChangeInstanceStatus(cid string, toStatus InstanceStatus, checkFunc func(status InstanceStatus) (bool, error)) error

	// Cleanup the left network interfaces
	GetAttachedNetworkInterfaceIds(cid string) []string
	CleanupInstanceNetworkInterfaces(cid string, eniIds []string) error
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

func (a InstanceManagerImpl) GetInstance(cid string) (inst *ecs.Instance, err error) {
	client, err := a.config.NewEcsClient("")
	if err != nil {
		return nil, err
	}

	args := ecs.CreateDescribeInstancesRequest()
	args.RegionId = a.config.OpenApi.GetRegion("")
	args.InstanceIds = fmt.Sprintf("[\"%s\"]", cid)

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		r, e := client.DescribeInstances(args)
		if e != nil {
			return e
		}
		if len(r.Instances.Instance) > 0 {
			inst = &r.Instances.Instance[0]
		}
		return e
	})

	return
}

func (a InstanceManagerImpl) CreateInstance(region string, request map[string]interface{}) (string, error) {
	conn, err := a.config.EcsTeaClient(region)
	if err != nil {
		return "", err
	}

	invoker := NewInvoker()
	invoker.AddCatcher(CreateInstanceCatcher_IdempotentProcessing)
	invoker.AddCatcher(CreateInstanceCatcher_TokenProcessing)
	invoker.AddCatcher(CreateInstanceCatcher_IpUsed)
	invoker.AddCatcher(CreateInstanceCatcher_IpUsed2)

	action := "CreateInstance"
	request["ClientToken"] = buildClientToken(action)
	runtime := util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	var cid string
	err = invoker.Run(func() error {
		resp, e := conn.DoRequest(StringPointer(action), nil, StringPointer("POST"), StringPointer("2014-05-26"), StringPointer("AK"), nil, request, &runtime)
		if e != nil {
			if IsExceptedErrors(e, []string{"IdempotentFailed"}) {
				// If the error is not 5xx, the client token should be updated
				request["ClientToken"] = buildClientToken(action)
			}
			return e
		}
		cid = fmt.Sprint(resp["InstanceId"])
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
	invoker.AddCatcher(IncorrectInstanceStatusCatcher)

	args := ecs.CreateDeleteInstanceRequest()
	args.InstanceId = cid
	args.Force = requests.NewBoolean(true)

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

func (a InstanceManagerImpl) GetAttachedNetworkInterfaceIds(cid string) []string {
	inst, _ := a.GetInstance(cid)
	eniIds := []string{}
	if inst != nil {
		for _, eni := range inst.NetworkInterfaces.NetworkInterface {
			eniIds = append(eniIds, eni.NetworkInterfaceId)
		}
	}
	return eniIds
}
func (a InstanceManagerImpl) CleanupInstanceNetworkInterfaces(cid string, eniIds []string) error {
	if len(eniIds) > 0 {
		return nil
	}
	client, err := a.config.NewEcsClient("")
	if err != nil {
		return err
	}
	invoker := NewInvoker()
	invoker.AddCatcher(NetworkInterfaceInvalidOperationInvalidEniStateCacher)

	req := ecs.CreateDeleteNetworkInterfaceRequest()
	req.RegionId = a.config.OpenApi.GetRegion("")

	for _, id := range eniIds {
		req.NetworkInterfaceId = id
		err = invoker.Run(func() error {
			_, err := client.DeleteNetworkInterface(req)
			a.log("DeleteNetworkInterface", err, id, "ok")
			return err
		})
		if err != nil && !IsExceptedErrors(err, []string{"InvalidEniId.NotFound"}) {
			return fmt.Errorf("After the instance %s is deleted, cleanup the network interface %s failed. Error:\n %#v.", cid, id, err)
		}
	}
	return nil
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
	} else if strings.HasPrefix(cid, "m-") {
		args.ResourceType = string(TagResourceImage)
	}
	return args

}

func StringPointer(s string) *string {
	return &s
}
