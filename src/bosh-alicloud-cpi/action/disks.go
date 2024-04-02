/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package action

import (
	"bosh-alicloud-cpi/alicloud"
	"bosh-alicloud-cpi/registry"
	"fmt"
	"math"
	"strconv"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
)

const (
	DefaultDiskCategory     = alicloud.DiskCategoryCloudEfficiency
	DefaultSystemDiskSizeGB = 40
	AmendSmallDiskSize      = true
)

type Disks struct {
	SystemDisk      DiskInfo
	EphemeralDisk   DiskInfo
	PersistentDisks []PersistentDisk
}

type DiskInfo struct {
	SizeRaw            interface{}       `json:"size"`
	Category           string            `json:"category"`
	Encrypted          *bool             `json:"encrypted,omitempty"`
	KMSKeyId           string            `json:"kms_key_id,omitempty"`
	DeleteWithInstance *bool             `json:"delete_with_instance,omitempty"`
	Tags               map[string]string `json:"tags"`
	sizeGB             int
	path               string
	ecsCategory        alicloud.DiskCategory
}

type PersistentDisk struct {
	Cid      string
	VolumeId string
	Path     string
}

func NewDiskInfo() DiskInfo {
	return DiskInfo{}
}

func NewDiskInfoWithSize(size int, props apiv1.DiskCloudProps) (DiskInfo, error) {
	d := NewDiskInfo()
	err := props.As(&d)
	if err != nil {
		return d, fmt.Errorf("bad format for DiskCloudProps %v", props)
	}
	d.SizeRaw = size
	return d.Validate(false)
}

func NewDisksWithProps(systemDisk DiskInfo, ephemeralDisk DiskInfo) (Disks, error) {
	r := Disks{systemDisk, ephemeralDisk, []PersistentDisk{}}

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
		r.EphemeralDisk.path = "/dev/vdb"
	}

	return r, nil
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
		c = alicloud.DiskCategory(a.Category)
	}

	if isSystem {
		a.path = "/dev/vda"
		if a.sizeGB == 0 {
			a.sizeGB = DefaultSystemDiskSizeGB
		}
	} else {
		a.path = "/dev/vdb"
	}
	a.ecsCategory = c

	//
	// `Alibaba Cloud` supported disk size is a range for each category in GB
	// cloud: [5, 2000]
	// cloud_efficiency: [20, 32768]
	// cloud_ssd: [20, 32768]
	if AmendSmallDiskSize {
		if a.ecsCategory == alicloud.DiskCategoryCloud {
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

func ConvertToGB(size float64) int {
	return int(math.Ceil(size / float64(1024)))
}

func (a DiskInfo) GetSizeGB() int {
	return a.sizeGB
}

func (a DiskInfo) GetCategory() alicloud.DiskCategory {
	return a.ecsCategory
}

func (a DiskInfo) GetPath() string {
	return a.path
}

func (a Disks) FillCreateInstanceArgs(globalEncrypt *bool, globalKmsKeyId string, request map[string]interface{}) {
	request["SystemDisk.Size"] = requests.NewInteger(a.SystemDisk.sizeGB)
	request["SystemDisk.Category"] = string(a.SystemDisk.ecsCategory)
	if v := a.SystemDisk.Encrypted; v != nil {
		request["SystemDisk.Encrypted"] = strconv.FormatBool(*v)
	}
	if v := a.SystemDisk.KMSKeyId; v != "" {
		request["SystemDisk.KMSKeyId"] = v
	}

	encrypt := a.EphemeralDisk.Encrypted
	if encrypt == nil {
		encrypt = globalEncrypt
		if encrypt == nil {
			encrypt = new(bool)
			*encrypt = false
		}
	}
	kmsKeyId := a.EphemeralDisk.KMSKeyId
	if kmsKeyId == "" {
		kmsKeyId = globalKmsKeyId
	}
	deleteWithInstance := a.EphemeralDisk.DeleteWithInstance
	if deleteWithInstance == nil {
		deleteWithInstance = new(bool)
		*deleteWithInstance = true
	}
	if a.EphemeralDisk.sizeGB > 0 {
		request["DataDisk.1.Size"] = strconv.Itoa(a.EphemeralDisk.sizeGB)
		request["DataDisk.1.Category"] = string(a.EphemeralDisk.GetCategory())
		request["DataDisk.1.Encrypted"] = strconv.FormatBool(*encrypt)
		request["DataDisk.1.DeleteWithInstance"] = strconv.FormatBool(*deleteWithInstance)
	}
	if kmsKeyId != "" {
		request["DataDisk.1.KMSKeyId"] = kmsKeyId
	}
}

func (a Disks) AssociatePersistentDisk(cid string, path string) {
	a.PersistentDisks = append(a.PersistentDisks, PersistentDisk{
		Cid:  cid,
		Path: path,
	})
}

func (a Disks) getPersistentDiskMap() map[string]interface{} {
	r := map[string]interface{}{}
	for _, pd := range a.PersistentDisks {
		r[pd.Cid] = pd.Path
	}
	return r
}

func (a Disks) AsRegistrySettings() registry.DisksSettings {
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
