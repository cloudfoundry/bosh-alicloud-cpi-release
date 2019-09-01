/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package mock

import (
	"bosh-alicloud-cpi/alicloud"
	"fmt"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
)

type InstanceManagerMock struct {
	mc *TestContext
}

func NewInstanceManagerMock(mc TestContext) alicloud.InstanceManager {
	return InstanceManagerMock{&mc}
}

func (a InstanceManagerMock) GetInstance(cid string) (*ecs.Instance, error) {
	i, ok := a.mc.Instances[cid]
	if !ok {
		return nil, nil
	} else {
		return i, nil
	}
}

func (a InstanceManagerMock) CreateInstance(region string, args *ecs.CreateInstanceRequest) (string, error) {
	id, inst := a.mc.NewInstance()

	inst.RegionId = args.RegionId
	inst.ZoneId = args.ZoneId
	// ...

	return id, nil
}

func (a InstanceManagerMock) ModifyInstanceAttribute(cid string, name string, description string) error {
	inst, ok := a.mc.Instances[cid]
	if !ok {
		return fmt.Errorf("ModifyInstanceAttribute instance not exists %s", cid)
	}
	inst.InstanceName = name
	inst.Description = description
	return nil
}

func (a InstanceManagerMock) AddTags(cid string, tags map[string]string) error {
	ok := true
	if strings.HasPrefix(cid, "i-") {
		_, ok = a.mc.Instances[cid]
	} else {
		_, ok = a.mc.Disks[cid]
	}
	if !ok {
		return fmt.Errorf("AddTags resource not exists %s", cid)
	}
	return nil
}

func (a InstanceManagerMock) DeleteInstance(cid string) error {
	inst, ok := a.mc.Instances[cid]
	if !ok {
		return fmt.Errorf("DeleteInstance instance not exists %s", cid)
	}
	if alicloud.InstanceStatus(inst.Status) != alicloud.Stopped {
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
	if alicloud.InstanceStatus(inst.Status) != alicloud.Stopped {
		return fmt.Errorf("StartInstance instance %s status %s is not Stopped", cid, inst.Status)
	}
	inst.Status = string(alicloud.Running)
	return nil
}

func (a InstanceManagerMock) StopInstance(cid string) error {
	inst, ok := a.mc.Instances[cid]
	if !ok {
		return fmt.Errorf("StopInstance instance not exists %s", cid)
	}
	if alicloud.InstanceStatus(inst.Status) != alicloud.Running {
		return fmt.Errorf("StopInstance instance %s status %s is not Running", cid, inst.Status)
	}
	inst.Status = string(alicloud.Stopped)
	a.mc.Instances[inst.InstanceId] = inst
	return nil
}

func (a InstanceManagerMock) RebootInstance(cid string) error {
	inst, ok := a.mc.Instances[cid]
	if !ok {
		return fmt.Errorf("StopInstance instance not exists %s", cid)
	}
	if alicloud.InstanceStatus(inst.Status) != alicloud.Running || alicloud.InstanceStatus(inst.Status) != alicloud.Stopped {
		return fmt.Errorf("RebootInstance instance %s status %s can't reboot", cid, inst.Status)
	}
	inst.Status = string(alicloud.Running)
	return nil
}

func (a InstanceManagerMock) GetInstanceStatus(cid string) (alicloud.InstanceStatus, error) {
	inst, ok := a.mc.Instances[cid]
	if !ok {
		return alicloud.Deleted, nil
	} else {
		return alicloud.InstanceStatus(inst.Status), nil
	}
}

func (a InstanceManagerMock) WaitForInstanceStatus(cid string, toStatus alicloud.InstanceStatus) (alicloud.InstanceStatus, error) {
	inst, ok := a.mc.Instances[cid]
	if !ok {
		if toStatus == alicloud.Deleted {
			return alicloud.Deleted, nil
		} else {
			return alicloud.Deleted, fmt.Errorf("WaitForInstanceStatus instance not exists %s", cid)
		}
	}
	if alicloud.InstanceStatus(inst.Status) != toStatus {
		return alicloud.InstanceStatus(inst.Status), fmt.Errorf("WaitForInstanceStatus instance %s excepted status %s but get %s", cid, toStatus, inst.Status)
	}
	return toStatus, nil
}

func (a InstanceManagerMock) ChangeInstanceStatus(cid string, toStatus alicloud.InstanceStatus, checkFunc func(status alicloud.InstanceStatus) (bool, error)) error {
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

func (a InstanceManagerMock) GetAttachedNetworkInterfaceIds (cid string) []string {
	inst, ok := a.mc.Instances[cid]
	eniIds := make([]string, 1,1)
	if !ok {
		return eniIds
	} else {
		for _, eni := range inst.NetworkInterfaces.NetworkInterface {
			eniIds = append(eniIds, eni.NetworkInterfaceId)
		}
	}
	return eniIds
}

func (a InstanceManagerMock) CleanupInstanceNetworkInterfaces(cid string, eniIds []string) error  {
	for _, id := range eniIds {
		delete(a.mc.NetworkInterfaces, id)
	}
	return nil
}