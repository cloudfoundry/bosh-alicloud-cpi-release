/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package action

import (
	"bosh-alicloud-cpi/alicloud"
	"bosh-alicloud-cpi/registry"
	"fmt"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
)

type AttachDiskMethod struct {
	CallContext
	disks     alicloud.DiskManager
	instances alicloud.InstanceManager
	registry  registry.Client
}

func NewAttachDiskMethod(cc CallContext, disks alicloud.DiskManager, instances alicloud.InstanceManager, rc registry.Client) AttachDiskMethod {
	return AttachDiskMethod{cc, disks, instances, rc}
}

func (a AttachDiskMethod) AttachDisk(vmCID apiv1.VMCID, diskCID apiv1.DiskCID) error {
	instCid := vmCID.AsString()
	diskCid := diskCID.AsString()
	device := ""

	if a.Config.Registry.IsEmpty() {
		err := a.instances.ChangeInstanceStatus(instCid, alicloud.Stopped, func(status alicloud.InstanceStatus) (bool, error) {
			switch status {
			case alicloud.Stopped:
				return true, nil
			case alicloud.Running:
				return false, a.instances.StopInstance(instCid)
			case alicloud.Stopping:
				return false, nil
			default:
				return false, fmt.Errorf("unexpect %s for StopInstance", status)
			}
		})

		if err != nil {
			return bosherr.WrapError(err, "stop instance failed")
		}
	}

	err := a.disks.ChangeDiskStatus(diskCid, alicloud.DiskStatusInUse, func(disk *ecs.Disk) (bool, error) {
		switch alicloud.DiskStatus(disk.Status) {
		case alicloud.DiskStatusInUse:
			device = alicloud.AmendDiskPath(disk.Device, alicloud.DiskCategory(disk.Category))
			return true, nil
		case alicloud.DiskStatusAvailable:
			return false, a.disks.AttachDisk(instCid, diskCid)
		case alicloud.DiskStatusAttaching:
			return false, nil
		default:
			return false, fmt.Errorf("unexcepted disk %s status %s", diskCid, disk.Status)
		}
	})

	if err != nil {
		return bosherr.WrapErrorf(err, "attach disk %s to %s failed", diskCid, instCid)
	}

	registryClient := a.registry
	agentSettings, _ := registryClient.Fetch(instCid)
	agentSettings.AttachPersistentDisk(diskCid, "", device)

	err = registryClient.Update(instCid, agentSettings)
	if err != nil {
		return bosherr.WrapErrorf(err, "update registry failed %s %s", diskCid, instCid)
	}

	if a.Config.Registry.IsEmpty() {
		err := a.instances.ChangeInstanceStatus(instCid, alicloud.Running, func(status alicloud.InstanceStatus) (bool, error) {
			switch status {
			case alicloud.Stopped:
				return false, a.instances.StartInstance(instCid)
			case alicloud.Starting:
				return false, nil
			case alicloud.Running:
				return true, nil
			default:
				return false, fmt.Errorf("unexpect %s for StopInstance", status)
			}
		})

		if err != nil {
			return bosherr.WrapError(err, "stop instance failed")
		}
	}

	return nil
}
