/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package mock

import (
	"bosh-alicloud-cpi/alicloud"
	"fmt"
	"strconv"
	"strings"
	"time"

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
	if a.mc.Flags["failGetDisks"] {
		return nil, fmt.Errorf("GetDisks injected failure (mock) for %s", instCid)
	}
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

func (a DiskManagerMock) CreateDisk(region string, args *ecs.CreateDiskRequest) (string, error) {
	size, err := strconv.Atoi(string(args.Size))
	if err != nil {
		return "", fmt.Errorf("convert size %s got an error: %#v.", args.Size, err)
	}
	if size < 5 || size > MaxDiskSizeGB {
		return "", fmt.Errorf("CreateDisk size too small or large %d", size)
	}
	if args.ZoneId == "" {
		return "", fmt.Errorf("CreateDisk zone can't be empty")
	}
	if args.DiskCategory == "" {
		return "", fmt.Errorf("CreateDisk category empty")
	}
	id, d := a.mc.NewDisk("")
	d.Size = size
	d.RegionId = args.RegionId
	d.ZoneId = args.ZoneId
	d.Status = string(alicloud.DiskStatusAvailable)
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

	if disk.Status != string(alicloud.DiskStatusAvailable) {
		return fmt.Errorf("AttachDisk unexcepted Disk.Status: %s", disk.Status)
	}

	if disk.InstanceId != "" {
		return fmt.Errorf("AttachDisk, disk %s alrady attach to %s", diskCid, disk.InstanceId)
	}

	disk.InstanceId = instCid
	disk.Status = string(alicloud.DiskStatusInUse)
	return nil
}

func (a DiskManagerMock) DetachDisk(instCid string, diskCid string) error {
	disk, ok := a.mc.Disks[diskCid]
	if !ok {
		return fmt.Errorf("DetachDisk disk not exists id=%s", diskCid)
	}

	if disk.Status != string(alicloud.DiskStatusInUse) {
		return fmt.Errorf("DetachDisk unexcepted Disk.Status: %s", disk.Status)
	}

	if disk.InstanceId != instCid {
		return fmt.Errorf("DetachDisk, disk %s instance except %s but get %s", diskCid, instCid, disk.InstanceId)
	}

	disk.InstanceId = ""
	disk.Status = string(alicloud.DiskStatusAvailable)
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

func (a DiskManagerMock) ModifyDiskCategory(diskCid string, category alicloud.DiskCategory, performanceLevel string) error {
	disk, ok := a.mc.Disks[diskCid]
	if !ok {
		return fmt.Errorf("ModifyDiskCategory disk %s not exists", diskCid)
	}

	// Simulate AliCloud refusing downgrades — real API returns
	// InvalidDiskCategory.NotSupported, which the CPI translates to
	// Bosh::Clouds::NotSupported for the director to handle.
	if isCategoryDowngrade(alicloud.DiskCategory(disk.Category), category) {
		return alicloud.NewProviderError(
			alicloud.InvalidDiskCategoryNotSupported,
			fmt.Sprintf("AliCloud refused in-place category change on disk %s to %s", diskCid, category))
	}

	disk.Category = string(category)
	if performanceLevel != "" {
		disk.PerformanceLevel = performanceLevel
	}
	// Simulate an async conversion that hasn't settled: leave the disk
	// non-Available so the post-ModifyDiskSpec WaitForDiskSpec has to wait.
	if a.mc.Flags["stallModifyDiskSpec"] {
		disk.Status = string(alicloud.DiskStatusCreating)
	}
	return nil
}

// categoryTier ranks disk categories from lowest (0) to highest performance.
// AliCloud's ModifyDiskSpec only permits forward (upgrade) transitions; any
// same-tier or lower target is rejected in the real API.
func categoryTier(c alicloud.DiskCategory) int {
	switch c {
	case alicloud.DiskCategoryCloud, alicloud.DiskCategoryEphemeral, alicloud.DiskCategoryEphemeralSSD:
		return 0
	case alicloud.DiskCategoryCloudEfficiency:
		return 1
	case alicloud.DiskCategoryCloudSSD:
		return 2
	case alicloud.DiskCategoryCloudESSD, alicloud.DiskCategoryCloudAuto:
		return 3
	default:
		return -1
	}
}

func isCategoryDowngrade(from, to alicloud.DiskCategory) bool {
	fromTier := categoryTier(from)
	toTier := categoryTier(to)
	if fromTier < 0 || toTier < 0 {
		return false
	}
	return toTier < fromTier
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

func (a DiskManagerMock) WaitForDiskStatus(diskCid string, toStatus alicloud.DiskStatus, opts ...time.Duration) (string, error) {
	disk, ok := a.mc.Disks[diskCid]
	if !ok {
		return "", fmt.Errorf("WaitForDiskStatus disk not exists id=%s", diskCid)
	}
	if disk.Status != string(toStatus) {
		return "", fmt.Errorf("WaitForDiskStatus %s disk.Status expect %s but get %s", diskCid, toStatus, disk.Status)
	}
	return disk.Device, nil
}

// WaitForDiskSpec records opts for test assertions, then checks that the disk is
// Available with the target category and (when non-empty) the target PL.
func (a DiskManagerMock) WaitForDiskSpec(diskCid, targetCategory, targetPL string, opts ...time.Duration) error {
	*a.mc.WaitForDiskSpecOpts = opts
	disk, ok := a.mc.Disks[diskCid]
	if !ok {
		return fmt.Errorf("WaitForDiskSpec disk not exists id=%s", diskCid)
	}
	if disk.Status != string(alicloud.DiskStatusAvailable) {
		return fmt.Errorf("WaitForDiskSpec %s: status=%s is not Available", diskCid, disk.Status)
	}
	if disk.Category != targetCategory {
		return fmt.Errorf("WaitForDiskSpec %s: category=%s does not match target %s", diskCid, disk.Category, targetCategory)
	}
	if targetPL != "" && disk.PerformanceLevel != targetPL {
		return fmt.Errorf("WaitForDiskSpec %s: PL=%s does not match target %s", diskCid, disk.PerformanceLevel, targetPL)
	}
	return nil
}

func (a DiskManagerMock) ChangeDiskStatus(cid string, toStatus alicloud.DiskStatus, checkFunc func(disk *ecs.Disk) (bool, error), opts ...time.Duration) error {
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
	if status == string(toStatus) {
		return nil
	} else {
		return fmt.Errorf("<MOCK> expect instance %s status is %s but get %s", cid, toStatus, status)
	}
}

// GetDiskPath mimics the real resolution: NVMe-capable instance types (the mock
// treats any type containing "9i", e.g. ecs.c9i.*, as NVMe) resolve to the
// nvme-... by-id path; everything else to the virtio-... by-id path. The
// "failDiskPath" flag forces a resolution error, mirroring an EcsTeaClient /
// DescribeInstanceTypes failure.
func (a DiskManagerMock) GetDiskPath(path, diskId, instanceType string, category alicloud.DiskCategory) (string, error) {
	if instanceType == "" || diskId == "" {
		return path, nil
	}
	if a.mc.Flags["failDiskPath"] {
		return path, fmt.Errorf("GetDiskPath injected failure (mock) for %s", instanceType)
	}
	suffix := diskId
	if parts := strings.SplitN(diskId, "-", 2); len(parts) == 2 {
		suffix = parts[1]
	}
	if strings.Contains(instanceType, "9i") {
		return "/dev/disk/by-id/nvme-Alibaba_Cloud_Elastic_Block_Storage_" + suffix, nil
	}
	return "/dev/disk/by-id/virtio-" + suffix, nil
}
