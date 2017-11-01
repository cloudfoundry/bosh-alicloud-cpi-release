/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package mock

import (
	"github.com/denverdino/aliyungo/ecs"
	"fmt"
	"bosh-alicloud-cpi/alicloud"
)

type InstanceManagerMock struct {
	mc *TestContext
}

func NewInstanceManagerMock(mc TestContext) (alicloud.InstanceManager) {
	return InstanceManagerMock{&mc}
}

func (a InstanceManagerMock) GetInstance(cid string) (*ecs.InstanceAttributesType, error) {
	i, ok := a.mc.Instances[cid]
	if !ok {
		return nil, nil
	}  else {
		return i, nil
	}
}

func (a InstanceManagerMock) CreateInstance(args ecs.CreateInstanceArgs) (string, error) {
	id, inst := a.mc.NewInstance()

	inst.RegionId = args.RegionId
	inst.ZoneId = args.ZoneId
	// ...

	return id, nil
}

func (a InstanceManagerMock) DeleteInstance(cid string) (error) {
	inst, ok := a.mc.Instances[cid]
	if !ok {
		return fmt.Errorf("DeleteInstance instance not exists %s", cid)
	}
	if inst.Status != ecs.Stopped {
		return fmt.Errorf("DeleteInstance instance %s status %s is not Stopped", cid, inst.Status)
	}
	delete(a.mc.Instances, cid)
	return nil
}

func (a InstanceManagerMock) StartInstance(cid string) error {
	inst, ok := a.mc.Instances[cid]
	if !ok {
		return fmt.Errorf("StartInstance instance not exists %s", cid)
	}
	if inst.Status != ecs.Stopped {
		return fmt.Errorf("StartInstance instance %s status %s is not Stopped", cid, inst.Status)
	}
	inst.Status = ecs.Running
	return nil
}

func (a InstanceManagerMock) StopInstance(cid string) error {
	inst, ok := a.mc.Instances[cid]
	if !ok {
		return fmt.Errorf("StopInstance instance not exists %s", cid)
	}
	if inst.Status != ecs.Running {
		return fmt.Errorf("StopInstance instance %s status %s is not Running", cid, inst.Status)
	}
	inst.Status = ecs.Stopped
	a.mc.Instances[inst.InstanceId] = inst
	return nil
}

func (a InstanceManagerMock) RebootInstance(cid string) error {
	inst, ok := a.mc.Instances[cid]
	if !ok {
		return fmt.Errorf("StopInstance instance not exists %s", cid)
	}
	if inst.Status != ecs.Running || inst.Status != ecs.Stopped {
		return fmt.Errorf("RebootInstance instance %s status %s can't reboot", cid, inst.Status)
	}
	inst.Status = ecs.Running
	return nil
}

func (a InstanceManagerMock) GetInstanceStatus(cid string) (ecs.InstanceStatus, error) {
	inst, ok := a.mc.Instances[cid]
	if !ok {
		return ecs.Deleted, nil
	} else {
		return inst.Status, nil
	}
}

func (a InstanceManagerMock) WaitForInstanceStatus(cid string, toStatus ecs.InstanceStatus) (ecs.InstanceStatus, error) {
	inst, ok := a.mc.Instances[cid]
	if !ok {
		if toStatus == ecs.Deleted {
			return ecs.Deleted, nil
		} else {
			return ecs.Deleted, fmt.Errorf("WaitForInstanceStatus instance not exists %s", cid)
		}
	}
	if inst.Status != toStatus {
		return inst.Status,fmt.Errorf("WaitForInstanceStatus instance %s excepted status %s but get %s", cid, toStatus, inst.Status)
	}
	return toStatus,nil
}

func (a InstanceManagerMock) ChangeInstanceStatus(cid string, toStatus ecs.InstanceStatus, checkFunc func(status ecs.InstanceStatus) (bool, error)) (error) {
	status, err := a.GetInstanceStatus(cid)
	if err != nil {
		return err
	}

	if status == toStatus {
		return nil
	}

	ok, err := checkFunc(status)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}

	status, err = a.GetInstanceStatus(cid)
	if err != nil {
		return err
	}

	if status == toStatus {
		return nil
	} else {
		return fmt.Errorf("<MOCK> expect instance %s status is %s but get %s", cid, toStatus, status)
	}
}