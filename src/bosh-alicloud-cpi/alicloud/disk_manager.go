/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package alicloud

import (
	"github.com/denverdino/aliyungo/ecs"
	"github.com/denverdino/aliyungo/common"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"time"
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

	WaitForDiskStatus(diskCid string, toStatus ecs.DiskStatus) (string, error)
}

type DiskManagerImpl struct {
	config Config
	logger boshlog.Logger
	region string
}

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
	disks, _, err := client.DescribeDisks(&args)
	return disks, err
}

func (a DiskManagerImpl) GetDisk(diskCid string) (*ecs.DiskItemType, error) {
	client := a.config.NewEcsClient()
	args := ecs.DescribeDisksArgs {
		RegionId: common.Region(a.config.OpenApi.RegionId),
		DiskIds: []string { diskCid, },
	}
	disks, _, err := client.DescribeDisks(&args)
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
	cid, err := client.CreateDisk(&args)
	a.log("CreateDisk", err, args, cid)
	return cid, err
}

func (a DiskManagerImpl) DeleteDisk(diskCid string) (error) {
	client := a.config.NewEcsClient()
	err := client.DeleteDisk(diskCid)
	a.log("DeleteDisk", err, diskCid, "ok")

	if err != nil {
		//
		// retry
		for i := 1; i <= DeleteDiskRetryCount; i++ {
			if strings.Contains(err.Error(), DeleteDiskRetryReason) {
				time.Sleep(DeleteDiskRetryInterval)
				err := client.DeleteDisk(diskCid)
				if err == nil {
					a.logger.Info("DiskManager", "DeleteDisk %s Done! (after %d retry)", diskCid, i)
					break
				}
				a.logger.Info("DiskManager", "DeleteDisk %s retry=%d", diskCid, i)
			} else {
				return err
			}
		}
	}
	return nil
}

func (a DiskManagerImpl) AttachDisk(instCid string, diskCid string) (error) {
	client := a.config.NewEcsClient()
	var args ecs.AttachDiskArgs
	args.InstanceId = instCid
	args.DiskId = diskCid
	err := client.AttachDisk(&args)
	a.log("AttachDisk", err, diskCid + " to " + instCid, "ok")
	return err
}

func (a DiskManagerImpl) DetachDisk(instCid string, diskCid string) (error) {
	client := a.config.NewEcsClient()
	var args ecs.DetachDiskArgs
	args.InstanceId = instCid
	args.DiskId = diskCid
	err := client.DetachDisk(args.InstanceId, args.DiskId)
	a.log("DetachDisk", err, diskCid + " from " + instCid, "ok")
	return err
}


func (a DiskManagerImpl) WaitForDiskStatus(diskCid string, toStatus ecs.DiskStatus) (string, error) {
	timeout := WaitTimeout
	for {
		disk, err := a.GetDisk(diskCid)

		if err != nil {
			return "", err
		}

		if disk == nil {
			return "", fmt.Errorf("disk missing id=%s", diskCid)
		}

		a.logger.Info("DiskManager", "Waiting disk %s from %v to %v", diskCid, disk.Status, toStatus)
		if disk.Status == toStatus {
			path := disk.Device
			if len(path) > 0 {
				// expect "/dev/vda" or "/dev/xvda"
				if len(path) >= 8 && strings.HasPrefix(path, "/dev/") {
					path = AmendDiskPath(disk.Device, disk.Category)
				} else {
					return "", fmt.Errorf("WaitForDiskStatus unexcepted disk.Device=%s", path)
				}
			}
			a.logger.Info("DiskManager", "Waiting disk %s to %s DONE! path=%s", diskCid, toStatus, path)
			return path, nil
		}

		if timeout > 0 {
			timeout -= WaitInterval
			time.Sleep(WaitInterval)
		} else {
			return "", bosherr.Error("WaitForInstanceStatus timeout")
		}
	}
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