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

const (
	ChangeDiskStatusTimeout       = time.Duration(300) * time.Second
	ChangeDiskStatusSleepInterval = time.Duration(5) * time.Second
)

type DiskManager interface {
	GetDisks(instCid string) ([]ecs.Disk, error)
	GetDisk(diskCid string) (*ecs.Disk, error)

	CreateDisk(args *ecs.CreateDiskRequest) (string, error)
	DeleteDisk(diskCid string) error

	AttachDisk(instCid string, diskCid string) error
	DetachDisk(instCid string, diskCid string) error

	ResizeDisk(diskCid string, sizeGB int) error
	ModifyDiskAttribute(diskCid string, name string, description string) error

	CreateSnapshot(diskCid string, snapshotName string) (string, error)
	DeleteSnapshot(snapshotCid string) error

	WaitForDiskStatus(diskCid string, toStatus DiskStatus) (string, error)
	ChangeDiskStatus(cid string, toStatus DiskStatus, checkFunc func(*ecs.Disk) (bool, error)) error
}

type DiskManagerImpl struct {
	config Config
	logger boshlog.Logger
	region string
}

var DeleteDiskCatcher = Catcher{"IncorrectDiskStatus.Initializing", 10, 15}

func NewDiskManager(config Config, logger boshlog.Logger) DiskManager {
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

func (a DiskManagerImpl) GetDisks(instCid string) (disks []ecs.Disk, err error) {
	client, err := a.config.NewEcsClient()
	if err != nil {
		return
	}
	args := ecs.CreateDescribeDisksRequest()
	args.RegionId = a.config.OpenApi.RegionId
	args.InstanceId = instCid

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		r, e := client.DescribeDisks(args)
		if r != nil {
			disks = r.Disks.Disk
		}
		return e
	})
	return
}

func (a DiskManagerImpl) GetDisk(diskCid string) (*ecs.Disk, error) {
	client, err := a.config.NewEcsClient()
	if err != nil {
		return nil, err
	}
	args := ecs.CreateDescribeDisksRequest()
	args.RegionId = a.config.OpenApi.RegionId
	bytes, _ := json.Marshal([]string{diskCid})
	args.DiskIds = string(bytes)

	invoker := NewInvoker()
	var disks []ecs.Disk
	err = invoker.Run(func() error {
		r, e := client.DescribeDisks(args)
		if r != nil {
			disks = r.Disks.Disk
		}
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

func (a DiskManagerImpl) CreateDisk(args *ecs.CreateDiskRequest) (cid string, err error) {
	args.RegionId = a.config.OpenApi.GetRegion()
	args.ClientToken = uuid.New().String()

	client, err := a.config.NewEcsClient()
	if err != nil {
		return
	}
	invoker := NewInvoker()
	err = invoker.Run(func() error {
		r, e := client.CreateDisk(args)
		if r != nil {
			cid = r.DiskId
		}
		a.log("CreateDisk", e, args, cid)
		return e
	})
	return
}

func (a DiskManagerImpl) DeleteDisk(diskCid string) error {
	client, err := a.config.NewEcsClient()
	if err != nil {
		return err
	}

	invoker := NewInvoker()
	invoker.AddCatcher(DeleteDiskCatcher)

	args := ecs.CreateDeleteDiskRequest()
	args.RegionId = a.config.OpenApi.RegionId
	args.DiskId = diskCid

	return invoker.Run(func() error {
		_, err := client.DeleteDisk(args)
		a.log("DeleteDisk", err, diskCid, "ok")
		return err
	})
}

func (a DiskManagerImpl) AttachDisk(instCid string, diskCid string) error {
	client, err := a.config.NewEcsClient()
	if err != nil {
		return err
	}

	args := ecs.CreateAttachDiskRequest()
	args.InstanceId = instCid
	args.DiskId = diskCid

	invoker := NewInvoker()
	return invoker.Run(func() error {
		_, err := client.AttachDisk(args)
		a.log("AttachDisk", err, diskCid+" to "+instCid, "ok")
		return err
	})
}

func (a DiskManagerImpl) DetachDisk(instCid string, diskCid string) error {
	client, err := a.config.NewEcsClient()
	if err != nil {
		return err
	}
	args := ecs.CreateDetachDiskRequest()
	args.InstanceId = instCid
	args.DiskId = diskCid

	invoker := NewInvoker()
	return invoker.Run(func() error {
		_, err := client.DetachDisk(args)
		a.log("DetachDisk", err, diskCid+" from "+instCid, "ok")
		return err
	})
}

func (a DiskManagerImpl) ResizeDisk(diskCid string, size int) error {
	client, err := a.config.NewEcsClient()
	if err != nil {
		return err
	}
	invoker := NewInvoker()

	args := ecs.CreateResizeDiskRequest()
	args.DiskId = diskCid
	args.NewSize = requests.NewInteger(size)

	return invoker.Run(func() error {
		_, err := client.ResizeDisk(args)
		a.log("ResizeDisk", err, diskCid, "ok")
		return err
	})
}

func (a DiskManagerImpl) ModifyDiskAttribute(diskCid string, name string, description string) error {
	client, err := a.config.NewEcsClient()
	if err != nil {
		return err
	}

	args := ecs.CreateModifyDiskAttributeRequest()
	args.DiskId = diskCid
	args.DiskName = name
	args.Description = description

	invoker := NewInvoker()
	return invoker.Run(func() error {
		_, e := client.ModifyDiskAttribute(args)
		a.log("ModifyDiskAttribute", e, diskCid, "ok")
		return e
	})
}

func (a DiskManagerImpl) CreateSnapshot(diskCid string, snapshotName string) (snapshotId string, err error) {
	client, err := a.config.NewEcsClient()
	if err != nil {
		return
	}

	args := ecs.CreateCreateSnapshotRequest()
	args.DiskId = diskCid
	args.SnapshotName = snapshotName
	args.ClientToken = uuid.New().String()

	invoker := NewInvoker()

	err = invoker.Run(func() error {
		r, e := client.CreateSnapshot(args)
		a.log("CreateSnapshot", e, diskCid, r.String())
		if r != nil {
			snapshotId = r.SnapshotId
		}
		return e
	})
	return
}

func (a DiskManagerImpl) DeleteSnapshot(snapshotCid string) error {
	client, err := a.config.NewEcsClient()
	if err != nil {
		return err
	}

	invoker := NewInvoker()

	args := ecs.CreateDeleteSnapshotRequest()
	args.SnapshotId = snapshotCid

	return invoker.Run(func() error {
		_, e := client.DeleteSnapshot(args)
		a.log("DeleteSnapshot", e, snapshotCid, "ok")
		return e
	})
}

func (a DiskManagerImpl) WaitForDiskStatus(diskCid string, toStatus DiskStatus) (string, error) {
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

		if DiskStatus(disk.Status) == toStatus {
			path := disk.Device
			a.logger.Info("DiskManager", "Waiting disk %s to %s DONE! path=%s", diskCid, toStatus, path)
			if len(path) > 0 {
				//
				// expect "/dev/vda" or "/dev/xvda"
				if len(path) >= 8 && strings.HasPrefix(path, "/dev/") {
					path = AmendDiskPath(disk.Device, DiskCategory(disk.Category))
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

func (a DiskManagerImpl) ChangeDiskStatus(cid string, toStatus DiskStatus, checkFunc func(*ecs.Disk) (bool, error)) error {
	timeout := ChangeDiskStatusTimeout
	for {
		disk, err := a.GetDisk(cid)
		if err != nil {
			return fmt.Errorf("get disk %s status failed %s", cid, err.Error())
		}

		ok, err := checkFunc(disk)
		status := "Deleted"
		if disk != nil {
			status = string(disk.Status)
		}

		if err != nil {
			a.logger.Error("DiskManager", "change %s from %s to %s failed %s", cid, status, toStatus, err.Error())
			return err
		}

		if ok {
			a.logger.Info("DiskManager", "change %s to %s done!", cid, toStatus)
			return nil
		} else {
			a.logger.Info("DiskManager", "changing %s from %s to %s ...", cid, status, toStatus)
		}

		timeout -= ChangeDiskStatusSleepInterval
		time.Sleep(ChangeDiskStatusSleepInterval)
		if timeout < 0 {
			return fmt.Errorf("change disk %s to %s timeout", cid, toStatus)
		}
	}
}

func AmendDiskPath(path string, category DiskCategory) string {
	//
	// 如果非普通云盘，需要去除x字母，如: xvdb -> vdb
	// if not normal Cloud need trim first x: xvdc -> vcd
	//
	// cloud:
	// cloud_efficiency:
	// cloud_ssd:
	// ephemeral_ssd:
	if category == DiskCategoryCloudEfficiency {
		if path[5] == 'x' {
			path = "/dev/" + string(path[6:])
		}
	}

	return path
}
