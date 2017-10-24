/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package alicloud

import (
	"github.com/denverdino/aliyungo/ecs"
	"github.com/denverdino/aliyungo/common"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"time"
)

type DiskManager interface {
	GetDisks(instCid string) ([]ecs.DiskItemType, error)
	GetDisk(diskCid string) (*ecs.DiskItemType, error)

	CreateDisk(sizeGB int, category ecs.DiskCategory) (string, error)
	DeleteDisk(diskCid string) (error)

	AttachDisk(instCid string, diskCid string) (error)
	DetachDisk(instCid string, diskCid string) (error)

	WaitForDiskStatus(diskCid string, toStatus ecs.DiskStatus) (string, error)
}

type DiskManagerImpl struct {
	config Config
	region string
}

func NewDiskManager(config Config) (DiskManager) {
	return DiskManagerImpl{
		config: config,
		region: config.OpenApi.RegionId,
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

func (a DiskManagerImpl) CreateDisk(sizeGB int, category ecs.DiskCategory) (string, error) {
	var args = ecs.CreateDiskArgs {
		RegionId: common.Region(a.region),
		DiskCategory: category,
		Size: sizeGB,
	}

	client := a.config.NewEcsClient()
	cid, err := client.CreateDisk(&args)
	return cid, err
}

func (a DiskManagerImpl) DeleteDisk(diskCid string) (error) {
	client := a.config.NewEcsClient()
	return client.DeleteDisk(diskCid)
}

func (a DiskManagerImpl) AttachDisk(instCid string, diskCid string) (error) {
	client := a.config.NewEcsClient()
	var args ecs.AttachDiskArgs
	args.InstanceId = instCid
	args.DiskId = diskCid
	err := client.AttachDisk(&args)
	return err
}

func (a DiskManagerImpl) DetachDisk(instCid string, diskCid string) (error) {
	client := a.config.NewEcsClient()

	var args ecs.DetachDiskArgs

	args.InstanceId = instCid
	args.DiskId = diskCid

	err := client.DetachDisk(args.InstanceId, args.DiskId)
	return err
}


func (a DiskManagerImpl) WaitForDiskStatus(diskCid string, toStatus ecs.DiskStatus) (string, error) {
	timeout := DefaultTimeoutMs
	for {
		disk, err := a.GetDisk(diskCid)

		if err != nil {
			return "", err
		}

		if disk.Status == toStatus {
			path := AmendDiskPath(disk.Device, disk.Category)
			return path, nil
		}

		if timeout > 0 {
			timeout -= 1000
			time.Sleep(time.Duration(DefaultWaitInterval) * time.Millisecond)
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