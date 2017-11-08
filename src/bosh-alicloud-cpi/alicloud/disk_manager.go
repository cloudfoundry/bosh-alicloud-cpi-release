/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package alicloud

import (
	"github.com/denverdino/aliyungo/ecs"
	"github.com/denverdino/aliyungo/common"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"fmt"
	"encoding/json"
	"strings"
)

type DiskManager interface {
	GetDisks(instCid string) ([]ecs.DiskItemType, error)
	GetDisk(diskCid string) (*ecs.DiskItemType, error)

	CreateDisk(sizeGB int, category ecs.DiskCategory, zone string) (string, error)
	DeleteDisk(diskCid string) (error)

	AttachDisk(instCid string, diskCid string) (error)
	DetachDisk(instCid string, diskCid string) (error)

	ModifyDiskAttribute(diskCid string, name string, description string) (error)

	CreateSnapshot(diskCid string, snapshotName string) (string, error)
	DeleteSnapshot(snapshotCid string) (error)

	WaitForDiskStatus(diskCid string, toStatus ecs.DiskStatus) (string, error)
}

type DiskManagerImpl struct {
	config Config
	logger boshlog.Logger
	region string
}

var DeleteDiskCatcher = Catcher{"IncorrectDiskStatus.Initializing", 10, 15}

func NewDiskManager(config Config, logger boshlog.Logger) (DiskManager) {
	return DiskManagerImpl{
		config: config,
		logger: logger,
		region: config.OpenApi.RegionId,
	}
}

func (a DiskManagerImpl) log(action string, err error, args interface{}, result string) {
	s, _ := json.Marshal(args)
	if err != nil {
		a.logger.Error("DiskManager", "%s failed args=%s err=%s", action, s, err)
	} else {
		a.logger.Info("DiskManager", "%s done! args=%s result=%s", action, s, result)
	}
}

func (a DiskManagerImpl) GetDisks(instCid string) ([]ecs.DiskItemType, error) {
	client := a.config.NewEcsClient()
	args := ecs.DescribeDisksArgs {
		RegionId: common.Region(a.config.OpenApi.RegionId),
		InstanceId: instCid,
	}

	invoker := NewInvoker()
	var disks []ecs.DiskItemType
	err := invoker.Run(func() (error){
		r, _, e := client.DescribeDisks(&args)
		disks = r
		return e
	})
	return disks, err
}

func (a DiskManagerImpl) GetDisk(diskCid string) (*ecs.DiskItemType, error) {
	client := a.config.NewEcsClient()
	args := ecs.DescribeDisksArgs {
		RegionId: common.Region(a.config.OpenApi.RegionId),
		DiskIds: []string { diskCid, },
	}

	invoker := NewInvoker()
	var disks []ecs.DiskItemType
	err := invoker.Run(func() (error) {
		r, _, e := client.DescribeDisks(&args)
		disks = r
		return e
	})
	if err != nil {
		return nil, bosherr.WrapErrorf(err, "GetDisk() Failed %s", args)
	}
	if len(disks) == 0 {
		return nil, nil
	}
	return &disks[0], nil
}

func (a DiskManagerImpl) CreateDisk(sizeGB int, category ecs.DiskCategory, zone string) (string, error) {
	var args = ecs.CreateDiskArgs {
		RegionId: common.Region(a.region),
		ZoneId: zone,
		DiskCategory: category,
		Size: sizeGB,
	}

	client := a.config.NewEcsClient()
	invoker := NewInvoker()
	var cid string
	err := invoker.Run(func() (error) {
		r, e := client.CreateDisk(&args)
		cid = r
		a.log("CreateDisk", e, args, cid)
		return e
	})
	return cid, err
}

func (a DiskManagerImpl) DeleteDisk(diskCid string) (error) {
	client := a.config.NewEcsClient()

	invoker := NewInvoker()
	invoker.AddCatcher(DeleteDiskCatcher)
	return invoker.Run(func() (error) {
		err := client.DeleteDisk(diskCid)
		a.log("DeleteDisk", err, diskCid, "ok")
		return err
	})

	return nil
}

func (a DiskManagerImpl) AttachDisk(instCid string, diskCid string) (error) {
	client := a.config.NewEcsClient()
	var args ecs.AttachDiskArgs
	args.InstanceId = instCid
	args.DiskId = diskCid

	invoker := NewInvoker()
	return invoker.Run(func() (error) {
		err := client.AttachDisk(&args)
		a.log("AttachDisk", err, diskCid + " to " + instCid, "ok")
		return err
	})
}

func (a DiskManagerImpl) DetachDisk(instCid string, diskCid string) (error) {
	client := a.config.NewEcsClient()
	var args ecs.DetachDiskArgs
	args.InstanceId = instCid
	args.DiskId = diskCid

	invoker := NewInvoker()
	return invoker.Run(func() (error) {
		err := client.DetachDisk(args.InstanceId, args.DiskId)
		a.log("DetachDisk", err, diskCid+" from "+instCid, "ok")
		return err
	})
}

func (a DiskManagerImpl) ResizeDisk(diskCid string, size int) (error) {
	//client := a.config.NewEcsClient()
	//var args ecs.Resi
	//args.InstanceId = instCid
	//args.DiskId = diskCid
	//
	//invoker := NewInvoker()
	//return invoker.Run(func() (error) {
	//	a.log("DetachDisk", err, diskCid+" from "+instCid, "ok")
	//	return err
	//})
	return fmt.Errorf("Unsupported")
}

func (a DiskManagerImpl) ModifyDiskAttribute(diskCid string, name string, description string) (error) {
	client := a.config.NewEcsClient()
	var args ecs.ModifyDiskAttributeArgs
	args.DiskId = diskCid
	args.DiskName = name
	args.Description = description

	invoker := NewInvoker()
	return invoker.Run(func() (error) {
		e := client.ModifyDiskAttribute(&args)
		a.log("ModifyDiskAttribute", e, diskCid, "ok")
		return e
	})
}

func (a DiskManagerImpl) CreateSnapshot(diskCid string, snapshotName string) (string, error) {
	client := a.config.NewEcsClient()
	var args ecs.CreateSnapshotArgs
	args.DiskId = diskCid
	args.SnapshotName = snapshotName

	invoker := NewInvoker()
	var snapshotId string
	err := invoker.Run(func() (error) {
		id, e := client.CreateSnapshot(&args)
		a.log("CreateSnapshot", e, diskCid, id)
		snapshotId = id
		return e
	})
	return snapshotId, err
}

func (a DiskManagerImpl) DeleteSnapshot(snapshotCid string) (error) {
	client := a.config.NewEcsClient()

	invoker := NewInvoker()
	return invoker.Run(func() (error) {
		e := client.DeleteSnapshot(snapshotCid)
		a.log("DeleteSnapshot", e, snapshotCid, "ok")
		return e
	})
}

func (a DiskManagerImpl) WaitForDiskStatus(diskCid string, toStatus ecs.DiskStatus) (string, error) {

	invoker := NewInvoker()

	result := ""
	ok, err := invoker.RunUntil(WaitTimeout, WaitInterval, func() (bool, error) {
		disk, e := a.GetDisk(diskCid)

		if e != nil {
			return false, e
		}

		if disk == nil {
			return false, fmt.Errorf("disk missing id=%s", diskCid)
		}

		if disk.Status == toStatus {
			path := disk.Device
			a.logger.Info("DiskManager", "Waiting disk %s to %s DONE! path=%s", diskCid, toStatus, path)
			if len(path) > 0 {
				//
				// expect "/dev/vda" or "/dev/xvda"
				if len(path) >= 8 && strings.HasPrefix(path, "/dev/") {
					path = AmendDiskPath(disk.Device, disk.Category)
					result = path
					return true, nil
				} else {
					return false, fmt.Errorf("WaitForDiskStatus unexcepted disk.Device=%s", path)
				}
			} else {
				return true, nil
			}
		} else {
			a.logger.Info("DiskManager", "Waiting disk %s from %v to %v", diskCid, disk.Status, toStatus)
			return false, nil
		}
	})

	if err != nil {
		return result, err
	}

	if !ok {
		return "", bosherr.Errorf("WaitForDisk %s to %s timeout", diskCid, toStatus)
	}
	return result, nil
}


func AmendDiskPath(path string, category ecs.DiskCategory) (string) {
	//
	// 如果非普通云盘，需要去除x字母，如: xvdb -> vdb
	// if not normal Cloud need trim first x: xvdc -> vcd
	//
	// cloud:
	// cloud_efficiency:
	// cloud_ssd:
	// ephemeral_ssd:
	if category == ecs.DiskCategoryCloudEfficiency {
		if path[5] == 'x' {
			path = "/dev/" + string(path[6:])
		}
	}

	return path
}