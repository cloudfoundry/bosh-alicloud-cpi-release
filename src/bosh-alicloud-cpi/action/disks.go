package action

import (
	"github.com/denverdino/aliyungo/ecs"
	"strconv"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	"math"
	"strings"
	"fmt"
	"bosh-alicloud-cpi/registry"
)

const (
	DefaultEcsDiskCategory = ecs.DiskCategoryCloudEfficiency
)

type Disks struct {
	SystemDisk DiskInfo
	EphemeralDisks []DiskInfo
	PersistentDisks []PersistentDisk
}

type DiskInfo struct {
	SizeRaw interface{}		`json:"size,omitempty"`
	Type string				`json:"type,omitempty"`
	sizeGB int
	ecsCategory ecs.DiskCategory
}

type PersistentDisk struct {
	Cid string
	VolumeId string
	Path string
}

func NewDisks(systemDisk DiskInfo, ephemeralDisks []DiskInfo) (Disks, error) {
	a := Disks {}
	d, err := systemDisk.Validate()
	if err != nil {
		return a, err
	}
	a.SystemDisk = d

	if systemDisk.sizeGB == 0 {
		a.SystemDisk.sizeGB = 50
		a.SystemDisk.ecsCategory = DefaultEcsDiskCategory
	}

	for _, d := range ephemeralDisks {
		d, err := d.Validate()
		if err != nil {
			return a, err
		}
		if d.sizeGB == 0 {
			return a, fmt.Errorf("ephemeral disk size can't be 0")
		}
		a.EphemeralDisks = append(a.EphemeralDisks, d)
	}
	return a, nil
}

func NewDiskInfo(size int, props apiv1.DiskCloudProps) (DiskInfo, error){
	var d DiskInfo
	err := props.As(&d)
	if err != nil {
		return d, fmt.Errorf("bad format for DiskCloudProps %v", props)
	}

	d.sizeGB = ConvertToGB(float64(size))
	d.ecsCategory = GetDiskCategory(d.Type)
	return d, nil
}

func (a DiskInfo) Validate() (DiskInfo, error) {
	switch a.SizeRaw.(type) {
	case int:
		a.sizeGB = ConvertToGB(a.SizeRaw.(float64))
	case float64:
		a.sizeGB = ConvertToGB(float64(a.SizeRaw.(int)))
	case string:
		s, _ := a.SizeRaw.(string)
		if strings.Compare("", s) == 0 {
			a.sizeGB = 0
		}
		s = strings.Replace(s,"_", "", -1)
		n, err := strconv.Atoi(s)
		if err != nil {
			return a, fmt.Errorf("error Disk.size format %v", a.SizeRaw)
		}
		a.sizeGB = ConvertToGB(float64(n))
	default:
		return a, fmt.Errorf("error Disk.size type %s", a.SizeRaw)
	}

	a.ecsCategory = GetDiskCategory(a.Type)
	return a, nil
}

func ConvertToGB(size float64) (int) {
	return int(math.Ceil(size / float64(1024)))
}

func GetDiskCategory(s string) (ecs.DiskCategory) {
	if strings.Compare("", s) != 0 {
		ec := ecs.DiskCategory(s)
		// TODO: validate disk category
		return ec
	} else {
		return DefaultEcsDiskCategory
	}
}

func (a DiskInfo) GetSizeGB() int {
	return a.sizeGB
}

func (a DiskInfo) GetCategory() ecs.DiskCategory {
 	return a.ecsCategory
}

func (a Disks) FillCreateInstanceArgs(args *ecs.CreateInstanceArgs) {
	args.SystemDisk.Size = a.SystemDisk.sizeGB
	args.SystemDisk.Category = a.SystemDisk.ecsCategory

	for _, d := range a.EphemeralDisks {
		args.DataDisk = append(args.DataDisk, ecs.DataDiskType {
			Size: d.sizeGB,
			Category: d.ecsCategory,
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
	return registry.DisksSettings {
		System: "/dev/vda",
		Ephemeral: "/dev/vdb",
		Persistent: a.getPersistentDiskMap(),
	}
}