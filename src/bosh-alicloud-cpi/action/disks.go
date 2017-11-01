/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/denverdino/aliyungo/ecs"
	"strconv"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	"math"
	"strings"
	"fmt"
	"bosh-alicloud-cpi/registry"
	"bosh-alicloud-cpi/alicloud"
)

const (
	DefaultDiskCategory     = ecs.DiskCategoryCloudEfficiency
	DefaultSystemDiskSizeGB = 40
	AmendSmallDiskSize      = true
)

type Disks struct {
	SystemDisk DiskInfo
	EphemeralDisk DiskInfo
	PersistentDisks []PersistentDisk
}

type DiskInfo struct {
	SizeRaw     interface{} `json:"size,omitempty"`
	Category    string      `json:"category,omitempty"`
	sizeGB      int
	path        string
	ecsCategory ecs.DiskCategory
}


type PersistentDisk struct {
	Cid string
	VolumeId string
	Path string
}

func NewDisks(systemDisk DiskInfo, ephemeralDisk DiskInfo) (Disks, error) {
	r := Disks {systemDisk, ephemeralDisk, []PersistentDisk{}}

	d, err := systemDisk.Validate(true)
	if err != nil {
		return r, bosherr.WrapError(err, "system_disk validate failed")
	}
	r.SystemDisk = d

	if r.SystemDisk.sizeGB == 0 {
		r.SystemDisk.sizeGB = DefaultSystemDiskSizeGB
	}

	r.EphemeralDisk, err = ephemeralDisk.Validate(false)
	if err != nil {
		return r, bosherr.WrapError(err, "ephemeral_disk validate failed")
	}

	if r.EphemeralDisk.sizeGB > 0 {
		r.EphemeralDisk.path = alicloud.AmendDiskPath("/dev/xvdb", r.EphemeralDisk.GetCategory())
	}

	return r, nil
}

func NewDiskInfo(size int, props apiv1.DiskCloudProps) (DiskInfo, error){
	var d DiskInfo
	err := props.As(&d)
	if err != nil {
		return d, fmt.Errorf("bad format for DiskCloudProps %v", props)
	}
	d.SizeRaw = size
	return d.Validate(false)
}

func (a DiskInfo) Validate(isSystem bool) (DiskInfo, error) {
	if a.SizeRaw != nil {
		switch a.SizeRaw.(type) {
		case int:
			a.sizeGB = ConvertToGB(float64(a.SizeRaw.(int)))
			if a.sizeGB < 1 {
				return a, fmt.Errorf("disk size too small %dMB < 1GB", a.SizeRaw.(int))
			}
		case float64:
			a.sizeGB = ConvertToGB(a.SizeRaw.(float64))
			if a.sizeGB < 1 {
				return a, fmt.Errorf("disk size too small %fMB < 1GB", a.SizeRaw.(float64))
			}
		case string:
			s, _ := a.SizeRaw.(string)
			if strings.Compare("", s) == 0 {
				a.sizeGB = 0
			}
			s = strings.Replace(s, "_", "", -1)
			n, err := strconv.Atoi(s)
			if err != nil {
				return a, fmt.Errorf("error Disk.size format %v", a.SizeRaw)
			}
			a.sizeGB = ConvertToGB(float64(n))
			if a.sizeGB < 1 {
				return a, fmt.Errorf("disk size too small %dMB < 1GB", n)
			}
		default:
			return a, fmt.Errorf("error Disk.size type %v", a.SizeRaw)
		}
	} else {
		a.sizeGB = 0
	}

	c := DefaultDiskCategory
	if len(strings.TrimSpace(a.Category)) > 0 {
		c = ecs.DiskCategory(a.Category)
	}

	if isSystem {
		if c != ecs.DiskCategoryCloudEfficiency && c != ecs.DiskCategoryCloudSSD {
			return a, fmt.Errorf("system disk only support: cloud_efficiency/cloud_ssd not %s", a.ecsCategory)
		}
		if a.sizeGB == 0 {
			a.sizeGB = DefaultSystemDiskSizeGB
		}
		a.path = "/dev/xvda"
	} else {
		if c != ecs.DiskCategoryCloud && c != ecs.DiskCategoryCloudEfficiency &&
			c != ecs.DiskCategoryCloudSSD && c != ecs.DiskCategoryEphemeralSSD {
			return a, fmt.Errorf("unsupported ephemeral disk type: %s", c)
		}
		a.path = "/dev/xvdb"
	}
	a.ecsCategory = c
	a.path = alicloud.AmendDiskPath(a.path, a.ecsCategory)

	//
	// `Alibaba Cloud` supported disk size is a range for each category in GB
	// cloud: [5, 2000]
	// cloud_efficiency: [20, 32768]
	// cloud_ssd: [20, 32768]
	if AmendSmallDiskSize {
		if a.ecsCategory == ecs.DiskCategoryCloud {
			if a.sizeGB < 5 {
				a.sizeGB = 5
			}
		} else {
			if a.sizeGB < 20 {
				a.sizeGB = 20
			}
		}
	}
	return a, nil
}

func ConvertToGB(size float64) (int) {
	return int(math.Ceil(size / float64(1024)))
}

func (a DiskInfo) GetSizeGB() int {
	return a.sizeGB
}

func (a DiskInfo) GetCategory() ecs.DiskCategory {
 	return a.ecsCategory
}

func (a DiskInfo) GetPath() (string) {
	return a.path
}

func (a Disks) FillCreateInstanceArgs(args *ecs.CreateInstanceArgs) {
	args.SystemDisk.Size = a.SystemDisk.sizeGB
	args.SystemDisk.Category = a.SystemDisk.ecsCategory

	if a.EphemeralDisk.sizeGB > 0 {
		args.DataDisk = append(args.DataDisk, ecs.DataDiskType{
			Size:               a.EphemeralDisk.sizeGB,
			Category:           a.EphemeralDisk.GetCategory(),
			DeleteWithInstance: true,
		})
	}
}

func (a Disks) AssociatePersistentDisk(cid string, path string) {
	a.PersistentDisks = append(a.PersistentDisks, PersistentDisk{
		Cid: cid,
		Path: path,
	})
}

func (a Disks) getPersistentDiskMap() (map[string]interface{}){
	r := map[string]interface{}{}
	for _, pd := range a.PersistentDisks {
		r[pd.Cid] = pd.Path
	}
	return r
}

func (a Disks) AsRegistrySettings() (registry.DisksSettings) {
	if a.EphemeralDisk.sizeGB > 0 {
		return registry.DisksSettings{
			System:     a.SystemDisk.path,
			Ephemeral:  a.EphemeralDisk.path,
			Persistent: a.getPersistentDiskMap(),
		}
	} else {
		return registry.DisksSettings{
			System:     a.SystemDisk.path,
			Persistent: a.getPersistentDiskMap(),
		}
	}
}