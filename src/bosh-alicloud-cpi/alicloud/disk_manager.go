/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package alicloud

import (
	"strings"

	openapi "github.com/alibabacloud-go/darabonba-openapi/v2/client"
	openapiutil "github.com/alibabacloud-go/openapi-util/service"
	"github.com/alibabacloud-go/tea/tea"

	util "github.com/alibabacloud-go/tea-utils/v2/service"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"

	"encoding/json"
	"fmt"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

const (
	ChangeDiskStatusTimeout       = time.Duration(300) * time.Second
	ChangeDiskStatusSleepInterval = time.Duration(5) * time.Second
)

type DiskManager interface {
	GetDisks(instCid string) ([]ecs.Disk, error)
	GetDisk(diskCid string) (*ecs.Disk, error)

	CreateDisk(region string, args *ecs.CreateDiskRequest) (string, error)
	DeleteDisk(diskCid string) error

	AttachDisk(instCid, diskCid string) error
	DetachDisk(instCid, diskCid string) error

	ResizeDisk(diskCid string, sizeGB int) error
	ModifyDiskAttribute(diskCid, name, description string) error
	ModifyDiskCategory(diskCid string, category DiskCategory, performanceLevel string) error

	CreateSnapshot(diskCid string, snapshotName string) (string, error)
	DeleteSnapshot(snapshotCid string) error

	WaitForDiskStatus(diskCid string, toStatus DiskStatus, opts ...time.Duration) (string, error)
	WaitForDiskSpec(diskCid, targetCategory, targetPL string, opts ...time.Duration) error
	ChangeDiskStatus(cid string, toStatus DiskStatus, checkFunc func(*ecs.Disk) (bool, error)) error

	GetDiskPath(path, diskId, instanceType string, category DiskCategory) (string, error)
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
		region: config.OpenApi.GetRegion(""),
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
	region, err := a.config.GetInstanceRegion(instCid)
	if err != nil {
		return
	}
	client, err := a.config.NewEcsClient(region)
	if err != nil {
		return
	}
	args := ecs.CreateDescribeDisksRequest()
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

func (a DiskManagerImpl) GetDisk(diskCid string) (disk *ecs.Disk, err error) {
	client, err := a.config.NewEcsClient("")
	if err != nil {
		return nil, err
	}

	disk, err = DescribeDisks(client, diskCid)
	if err != nil {
		return
	}
	if disk == nil {
		regions, e := a.config.GetCrossRegions()
		if e != nil {
			return nil, e
		}
		if len(regions) <= 0 {
			return
		}
		for _, region := range regions {
			client, e := a.config.NewEcsClient(region)
			if e != nil {
				return nil, e
			}
			d, e := DescribeDisks(client, diskCid)
			if e != nil {
				return nil, e
			}
			if d != nil {
				disk = d
				return
			}
		}
	}
	return
}

func (a DiskManagerImpl) CreateDisk(region string, args *ecs.CreateDiskRequest) (cid string, err error) {
	args.ClientToken = buildClientToken(args.GetActionName())

	client, err := a.config.NewEcsClient(region)
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
	client, err := a.config.NewEcsClient("")
	if err != nil {
		return err
	}

	invoker := NewInvoker()
	invoker.AddCatcher(DeleteDiskCatcher)

	args := ecs.CreateDeleteDiskRequest()
	args.DiskId = diskCid

	return invoker.Run(func() error {
		_, err := client.DeleteDisk(args)
		a.log("DeleteDisk", err, diskCid, "ok")
		return err
	})
}

func (a DiskManagerImpl) AttachDisk(instCid string, diskCid string) error {
	client, err := a.config.NewEcsClient("")
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

func (a DiskManagerImpl) DetachDisk(instCid, diskCid string) error {
	client, err := a.config.NewEcsClient("")
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
	client, err := a.config.NewEcsClient("")
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

func (a DiskManagerImpl) ModifyDiskAttribute(diskCid, name, description string) error {
	client, err := a.config.NewEcsClient("")
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

// ModifyDiskCategory changes a disk's category and/or performance level in-place
// via ModifyDiskSpec. Category moves are forward-only (e.g. cloud_efficiency ->
// cloud_essd); performanceLevel applies to ESSD (PL0-PL3), "" leaves it unchanged.
// The disk must be detached first. Raw AliCloud errors are returned unchanged;
// callers can test refusal with IsCategoryRefusedError to fall back to snapshot+recreate.
func (a DiskManagerImpl) ModifyDiskCategory(diskCid string, category DiskCategory, performanceLevel string) error {
	client, err := a.config.NewEcsClient("")
	if err != nil {
		return err
	}

	args := ecs.CreateModifyDiskSpecRequest()
	args.DiskId = diskCid
	args.DiskCategory = string(category)
	if performanceLevel != "" {
		args.PerformanceLevel = performanceLevel
	}

	invoker := NewInvoker()
	return invoker.Run(func() error {
		_, e := client.ModifyDiskSpec(args)
		a.log("ModifyDiskCategory", e, diskCid, string(category))
		return e
	})
}

// IsCategoryRefusedError reports whether err is AliCloud's in-place category
// refusal from ModifyDiskSpec, used to decide on a snapshot+recreate fallback.
func IsCategoryRefusedError(err error) bool {
	return err != nil && IsExceptedErrors(err, DiskCategoryUnsupportedCodes)
}

func (a DiskManagerImpl) CreateSnapshot(diskCid, snapshotName string) (snapshotId string, err error) {
	client, err := a.config.NewEcsClient("")
	if err != nil {
		return
	}

	args := ecs.CreateCreateSnapshotRequest()
	args.DiskId = diskCid
	args.SnapshotName = snapshotName
	args.ClientToken = buildClientToken(args.GetActionName())

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
	client, err := a.config.NewEcsClient("")
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

// WaitForDiskStatus polls until the disk reaches toStatus. The timeout and
// interval are optional: WaitForDiskStatus(cid, status) uses the default
// WaitTimeout/WaitInterval; pass WaitForDiskStatus(cid, status, timeout) or
// (cid, status, timeout, interval) to override. Extra args beyond two are
// ignored.
func (a DiskManagerImpl) WaitForDiskStatus(diskCid string, toStatus DiskStatus, opts ...time.Duration) (string, error) {
	invoker := NewInvoker()

	timeout, interval := WaitTimeout, WaitInterval
	if len(opts) >= 1 {
		timeout = opts[0]
	}
	if len(opts) >= 2 {
		interval = opts[1]
	}

	result := ""
	ok, err := invoker.RunUntil(timeout, interval, func() (bool, error) {
		disk, e := a.GetDisk(diskCid)

		if e != nil {
			return false, e
		}

		if disk == nil {
			return false, fmt.Errorf("disk missing id=%s", diskCid)
		}

		if DiskStatus(disk.Status) != toStatus {
			a.logger.Info("DiskManager", "Waiting disk %s from %v to %v", diskCid, disk.Status, toStatus)
			return false, nil
		}
		return true, nil
	})

	if err != nil {
		return result, err
	}

	if !ok {
		return "", bosherr.Errorf("WaitForDisk %s to %s timeout", diskCid, toStatus)
	}
	return result, nil
}

// WaitForDiskSpec polls until the disk is Available with the target category and,
// when targetPL is non-empty, the target performance level. ModifyDiskSpec is
// asynchronous: the disk may still show Available before it transitions to
// Modifying, so checking category+PL prevents a spurious early return.
func (a DiskManagerImpl) WaitForDiskSpec(diskCid, targetCategory, targetPL string, opts ...time.Duration) error {
	invoker := NewInvoker()

	timeout, interval := WaitTimeout, WaitInterval
	if len(opts) >= 1 {
		timeout = opts[0]
	}
	if len(opts) >= 2 {
		interval = opts[1]
	}

	ok, err := invoker.RunUntil(timeout, interval, func() (bool, error) {
		disk, e := a.GetDisk(diskCid)
		if e != nil {
			return false, e
		}
		if disk == nil {
			return false, fmt.Errorf("disk missing id=%s", diskCid)
		}
		if DiskStatus(disk.Status) != DiskStatusAvailable {
			a.logger.Info("DiskManager", "WaitForDiskSpec %s: status=%s waiting for Available", diskCid, disk.Status)
			return false, nil
		}
		if disk.Category != targetCategory {
			a.logger.Info("DiskManager", "WaitForDiskSpec %s: category=%s waiting for %s", diskCid, disk.Category, targetCategory)
			return false, nil
		}
		if targetPL != "" && disk.PerformanceLevel != targetPL {
			a.logger.Info("DiskManager", "WaitForDiskSpec %s: PL=%s waiting for %s", diskCid, disk.PerformanceLevel, targetPL)
			return false, nil
		}
		return true, nil
	})

	if err != nil {
		return err
	}
	if !ok {
		return bosherr.Errorf("WaitForDiskSpec %s to category=%s PL=%s timeout", diskCid, targetCategory, targetPL)
	}
	return nil
}

func (a DiskManagerImpl) ChangeDiskStatus(cid string, toStatus DiskStatus, checkFunc func(*ecs.Disk) (bool, error)) error {
	timeout := ChangeDiskStatusTimeout
	for {
		disk, err := a.GetDisk(cid)
		if err != nil {
			return fmt.Errorf("get disk %s status failed %s", cid, err.Error())
		}
		if disk == nil {
			return fmt.Errorf("get disk %s failed and disk is nil", cid)
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
	// if not normal Cloud need trim first x, like: xvdc -> vdc
	//
	// cloud:
	// cloud_efficiency:
	// cloud_ssd:
	// ephemeral_ssd:
	if category == DiskCategoryCloudEfficiency || category == DiskCategoryCloudSSD ||
		category == DiskCategoryCloudESSD || category == DiskCategoryCloudAuto {
		if path[5] == 'x' {
			path = "/dev/" + string(path[6:])
		}
	}

	return path
}

// GetDiskPath resolves the device path the agent should use for a disk. For
// NVMe-capable instance types it returns the /dev/disk/by-id/nvme-... path; for
// other types the /dev/disk/by-id/virtio-... path. The returned string is always
// a usable best-effort path (the amended input path when resolution can't run),
// but a non-nil error means the NVMe-vs-virtio determination failed and the path
// may be wrong — callers on the create_vm hot path must treat that as fatal, since
// a wrong path leaves the agent hanging until the 600s ping timeout.
func (a DiskManagerImpl) GetDiskPath(path, diskId, instanceType string, category DiskCategory) (string, error) {
	amendPath := AmendDiskPath(path, category)

	if instanceType == "" || diskId == "" {
		return amendPath, nil
	}

	conn, err := a.config.EcsTeaClient("")
	if err != nil {
		a.log("EcsTeaClient", err, nil, "")
		return amendPath, err
	}

	invoker := NewInvoker()
	invoker.AddCatcher(CreateInstanceCatcher_IdempotentProcessing)
	invoker.AddCatcher(CreateInstanceCatcher_TokenProcessing)
	invoker.AddCatcher(CreateInstanceCatcher_IpUsed)
	invoker.AddCatcher(CreateInstanceCatcher_IpUsed2)

	action := "DescribeInstanceTypes"
	queries := map[string]interface{}{
		"InstanceTypes.1": instanceType,
		"NvmeSupport":     "required",
	}
	params := &openapi.Params{
		Action:      tea.String(action),
		Version:     tea.String("2014-05-26"),
		Protocol:    tea.String("HTTPS"),
		Method:      tea.String("POST"),
		AuthType:    tea.String("AK"),
		Style:       tea.String("RPC"),
		Pathname:    tea.String("/"),
		ReqBodyType: tea.String("json"),
		BodyType:    tea.String("json"),
	}
	runtime := &util.RuntimeOptions{}
	runtime.SetAutoretry(true)
	request := &openapi.OpenApiRequest{
		Query: openapiutil.Query(queries),
	}
	err = invoker.Run(func() error {
		resp, e := conn.CallApi(params, request, runtime)
		if e != nil {
			return e
		}
		if resp["body"] == nil {
			return fmt.Errorf("invoking DescribeInstanceTypes got an empty response: %v", resp)
		}
		instanceTypes := resp["body"].(map[string]interface{})["InstanceTypes"]
		if instanceTypes != nil && instanceTypes.(map[string]interface{})["InstanceType"] != nil &&
			len(instanceTypes.(map[string]interface{})["InstanceType"].([]interface{})) > 0 {
			amendPath = "/dev/disk/by-id/nvme-Alibaba_Cloud_Elastic_Block_Storage_" + strings.Split(diskId, "-")[1]
		} else {
			amendPath = "/dev/disk/by-id/virtio-" + strings.Split(diskId, "-")[1]
		}
		return e
	})
	if err != nil {
		a.log(action, err, queries, "")
		return amendPath, err
	}

	return amendPath, nil
}

func DescribeDisks(client *ecs.Client, diskId string) (disk *ecs.Disk, err error) {
	args := ecs.CreateDescribeDisksRequest()
	bytes, _ := json.Marshal([]string{diskId})
	args.DiskIds = string(bytes)

	invoker := NewInvoker()
	err = invoker.Run(func() error {
		r, e := client.DescribeDisks(args)
		if r != nil && len(r.Disks.Disk) > 0 {
			disk = &r.Disks.Disk[0]
		}
		return e
	})
	return
}
