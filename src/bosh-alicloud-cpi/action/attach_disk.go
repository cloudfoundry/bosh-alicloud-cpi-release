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
	"github.com/cloudfoundry/bosh-cpi-go/apiv1"
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
	_, err := a.attach(vmCID, diskCID, "v1")
	return err
}

func (a AttachDiskMethod) AttachDiskV2(vmCID apiv1.VMCID, diskCID apiv1.DiskCID) (apiv1.DiskHint, error) {
	diskPath, err := a.attach(vmCID, diskCID, "v2")
	if err != nil {
		return apiv1.DiskHint{}, err
	}
	return apiv1.NewDiskHintFromString(fmt.Sprintf("%s", diskPath)), nil
}

func (a AttachDiskMethod) attach(vmCID apiv1.VMCID, diskCID apiv1.DiskCID, cpiVersion string) (interface{}, error) {
	instCid := vmCID.AsString()
	diskCid := diskCID.AsString()
	device := ""

	// The registry has been deprecated from the CPI V2, so it is always empty when V2
	if cpiVersion == "v1" {
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
				return device, bosherr.WrapError(err, "stop instance failed")
			}
		}
	}

	err := a.disks.ChangeDiskStatus(diskCid, alicloud.DiskStatusInUse, func(disk *ecs.Disk) (bool, error) {
		switch alicloud.DiskStatus(disk.Status) {
		case alicloud.DiskStatusInUse:
			inst, er := a.instances.GetInstance(instCid)
			var derr error
			if er == nil {
				device, derr = a.disks.GetDiskPath(disk.Device, diskCid, inst.InstanceType, alicloud.DiskCategory(disk.Category))
			} else {
				device, derr = a.disks.GetDiskPath(disk.Device, diskCid, "", alicloud.DiskCategory(disk.Category))
			}
			// The disk is already attached; a resolution failure is non-fatal here
			// (unlike create_vm's ephemeral path). Log it and proceed with the
			// best-effort path GetDiskPath returned.
			if derr != nil {
				a.Logger.Warn("AttachDisk", "resolve disk path for %s on %s failed, using best-effort path %s: %s", diskCid, instCid, device, derr)
			}
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
		return device, bosherr.WrapErrorf(err, "attach disk %s to %s failed", diskCid, instCid)
	}

	registryClient := a.registry
	agentSettings, _ := registryClient.Fetch(instCid)
	agentSettings.AttachPersistentDisk(diskCid, "", device)

	err = registryClient.Update(instCid, agentSettings)
	if err != nil {
		return device, bosherr.WrapErrorf(err, "update registry failed %s %s", diskCid, instCid)
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
			return device, bosherr.WrapError(err, "stop instance failed")
		}
	}

	return device, nil
}
