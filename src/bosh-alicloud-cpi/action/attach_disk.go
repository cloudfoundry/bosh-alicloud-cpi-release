/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package action

import (
	"bosh-alicloud-cpi/alicloud"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	"github.com/denverdino/aliyungo/ecs"
	"bosh-alicloud-cpi/registry"
	"fmt"
)

type AttachDiskMethod struct {
	CallContext
	disks alicloud.DiskManager
	instances alicloud.InstanceManager
	registry registry.Client
}

func NewAttachDiskMethod(cc CallContext,disks alicloud.DiskManager, instances alicloud.InstanceManager, rc registry.Client) AttachDiskMethod {
	return AttachDiskMethod{cc, disks, instances, rc}
}


func (a AttachDiskMethod) AttachDisk(vmCID apiv1.VMCID, diskCID apiv1.DiskCID) error {
	instCid := vmCID.AsString()
	diskCid := diskCID.AsString()
	device := ""

	if a.Config.Registry.IsEmpty() {
		err := a.instances.ChangeInstanceStatus(instCid, ecs.Stopped, func(status ecs.InstanceStatus) (bool, error) {
			switch status {
			case ecs.Stopped:
				return true, nil
			case ecs.Running:
				return false, a.instances.StopInstance(instCid)
			case ecs.Stopping:
				return false, nil
			default:
				return false, fmt.Errorf("unexpect %s for StopInstance", status)
			}
		})

		if err != nil {
			return a.WrapError(err, "stop instance failed")
		}
	}

	err := a.disks.ChangeDiskStatus(diskCid, ecs.DiskStatusInUse, func(disk *ecs.DiskItemType) (bool, error) {
		if disk == nil {
			return false, fmt.Errorf("missing disk %s", diskCid)
		}
		switch disk.Status {
		case ecs.DiskStatusInUse:
			device = alicloud.AmendDiskPath(disk.Device, disk.Category)
			return true, nil
		case ecs.DiskStatusAvailable:
			return false, a.disks.AttachDisk(instCid, diskCid)
		case ecs.DiskStatusAttaching:
			return false, nil
		default:
			return false, fmt.Errorf("unexcepted disk %s status %s", diskCid, disk.Status)
		}
	})

	if err != nil {
		return a.WrapErrorf(err, "attach disk %s to %s failed", diskCid, instCid)
	}

	registryClient := a.registry
	agentSettings, _ := registryClient.Fetch(instCid)
	agentSettings.AttachPersistentDisk(diskCid, "", device)

	err = registryClient.Update(instCid, agentSettings)
	if err != nil {
		return a.WrapErrorf(err, "update registry failed %s %s", diskCid, instCid)
	}

	if a.Config.Registry.IsEmpty() {
		err := a.instances.ChangeInstanceStatus(instCid, ecs.Running, func(status ecs.InstanceStatus) (bool, error) {
			switch status {
			case ecs.Stopped:
				return false, a.instances.StartInstance(instCid)
			case ecs.Starting:
				return false, nil
			case ecs.Running:
				return true, nil
			default:
				return false, fmt.Errorf("unexpect %s for StopInstance", status)
			}
		})

		if err != nil {
			return a.WrapError(err, "stop instance failed")
		}
	}

	return nil
}

