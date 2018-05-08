/*
 * Copyright (C) 2017-2018 Alibaba Group Holding Limited
 */
package mock

import (
	"bosh-alicloud-cpi/alicloud"
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
)

const (
	MaxDiskSizeGB = 32768
)

type DiskManagerMock struct {
	mc *TestContext
}

func NewDiskManagerMock(mc TestContext) alicloud.DiskManager {
	return DiskManagerMock{&mc}
}

func (a DiskManagerMock) GetDisks(instCid string) ([]ecs.Disk, error) {
	r := []ecs.Disk{}
	for _, d := range a.mc.Disks {
		if d.InstanceId == instCid {
			r = append(r, *d)
		}
	}
	return r, nil
}

func (a DiskManagerMock) GetDisk(diskCid string) (*ecs.Disk, error) {
	if d, ok := a.mc.Disks[diskCid]; ok {
		return d, nil
	} else {
		return nil, nil
	}
}

func (a DiskManagerMock) CreateDisk(args ecs.CreateDiskRequest) (string, error) {
	if args.Size < 5 || args.Size > MaxDiskSizeGB {
		return "", fmt.Errorf("CreateDisk size too small or large %d", args.Size)
	}
	if args.ZoneId == "" {
		return "", fmt.Errorf("CreateDisk zone can't be empty")
	}
	if args.DiskCategory == "" {
		return "", fmt.Errorf("CreateDisk category empty")
	}
	id, d := a.mc.NewDisk("")
	d.Size = args.Size
	d.RegionId = args.RegionId
	d.ZoneId = args.ZoneId
	d.Status = ecs.DiskStatusAvailable
	d.Category = args.DiskCategory
	return id, nil
}

func (a DiskManagerMock) DeleteDisk(diskCid string) error {
	if _, ok := a.mc.Disks[diskCid]; ok {
		delete(a.mc.Disks, diskCid)
		return nil
	} else {
		return fmt.Errorf("DeleteDisk disk not exists id=%s", diskCid)
	}
}

func (a DiskManagerMock) AttachDisk(instCid string, diskCid string) error {
	disk, ok := a.mc.Disks[diskCid]
	if !ok {
		return fmt.Errorf("AttachDisk disk not exists id=%s", diskCid)
	}

	if disk.Status != ecs.DiskStatusAvailable {
		return fmt.Errorf("AttachDisk unexcepted Disk.Status: %s", disk.Status)
	}

	if disk.InstanceId != "" {
		return fmt.Errorf("AttachDisk, disk %s alrady attach to %s", diskCid, disk.InstanceId)
	}

	disk.InstanceId = instCid
	disk.Status = ecs.DiskStatusInUse
	return nil
}

func (a DiskManagerMock) DetachDisk(instCid string, diskCid string) error {
	disk, ok := a.mc.Disks[diskCid]
	if !ok {
		return fmt.Errorf("DetachDisk disk not exists id=%s", diskCid)
	}

	if disk.Status != ecs.DiskStatusInUse {
		return fmt.Errorf("DetachDisk unexcepted Disk.Status: %s", disk.Status)
	}

	if disk.InstanceId != instCid {
		return fmt.Errorf("DetachDisk, disk %s instance except %s but get %s", diskCid, instCid, disk.InstanceId)
	}

	disk.InstanceId = ""
	disk.Status = ecs.DiskStatusAvailable
	return nil
}

func (a DiskManagerMock) ResizeDisk(diskCid string, sizeGB int) error {
	disk, ok := a.mc.Disks[diskCid]
	if !ok {
		return fmt.Errorf("ResizeDisk disk %s not exists", diskCid)
	}

	if sizeGB < disk.Size {
		return fmt.Errorf("ResizeDisk %s to %d can not smaller than before %d", diskCid, sizeGB, disk.Size)
	}

	disk.Size = sizeGB
	return nil
}

func (a DiskManagerMock) ModifyDiskAttribute(diskCid string, name string, description string) error {
	disk, ok := a.mc.Disks[diskCid]
	if !ok {
		return fmt.Errorf("ModifyDiskAttribute disk %s not exists", diskCid)
	}

	disk.DiskName = name
	disk.Description = description
	return nil
}

func (a DiskManagerMock) CreateSnapshot(diskCid string, snapshotName string) (string, error) {
	_, ok := a.mc.Disks[diskCid]
	if !ok {
		return "", fmt.Errorf("CreateSnapshot disk %s not exists", diskCid)
	}
	ssid := NewSnapshotId()
	a.mc.Snapshots[ssid] = diskCid
	return ssid, nil
}

func (a DiskManagerMock) DeleteSnapshot(snapshotCid string) error {
	_, ok := a.mc.Snapshots[snapshotCid]
	if !ok {
		return fmt.Errorf("DeleteSnapshot %s not found", snapshotCid)
	}
	delete(a.mc.Snapshots, snapshotCid)
	return nil
}

func (a DiskManagerMock) WaitForDiskStatus(diskCid string, toStatus ecs.DiskStatus) (string, error) {
	disk, ok := a.mc.Disks[diskCid]
	if !ok {
		return "", fmt.Errorf("WaitForDiskStatus disk not exists id=%s", diskCid)
	}
	if disk.Status != toStatus {
		return "", fmt.Errorf("WaitForDiskStatus %s disk.Status expect %s but get %s", diskCid, toStatus, disk.Status)
	}
	return disk.Device, nil
}

func (a DiskManagerMock) ChangeDiskStatus(cid string, toStatus ecs.DiskStatus, checkFunc func(*ecs.DiskItemType) (bool, error)) error {
	disk, err := a.GetDisk(cid)
	if err != nil {
		return err
	}

	ok, err := checkFunc(disk)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}

	disk, err = a.GetDisk(cid)
	if err != nil {
		return err
	}

	if disk == nil {
		return fmt.Errorf("missing disk %s", cid)
	}

	status := disk.Status
	if status == toStatus {
		return nil
	} else {
		return fmt.Errorf("<MOCK> expect instance %s status is %s but get %s", cid, toStatus, status)
	}
}
