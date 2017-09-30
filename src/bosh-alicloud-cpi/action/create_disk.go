package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	"bosh-alicloud-cpi/alicloud"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/denverdino/aliyungo/common"
)

type CreateDiskMethod struct {
	runner alicloud.Runner
}

func NewCreateDiskMethod(runner alicloud.Runner) CreateDiskMethod {
	return CreateDiskMethod{runner}
}

func (a CreateDiskMethod) CreateDisk(size int, props apiv1.DiskCloudProps, vmcid *apiv1.VMCID) (apiv1.DiskCID, error) {
	client := a.runner.NewClient()
	instid := vmcid.AsString()

	inst, err := a.runner.GetInstance(instid)

	if err != nil {
		return apiv1.DiskCID{}, bosherr.WrapError(err, "GetInstance Failed")
	}

	if inst == nil {
		return apiv1.DiskCID{}, bosherr.WrapErrorf(err, "Missing Vm cid = %s", instid)
	}

	zoneId := inst.ZoneId

	var args = ecs.CreateDiskArgs {
		RegionId: common.Region(a.runner.Config.OpenApi.RegionId),
		ZoneId: zoneId,
		DiskName: "",			//TODO
		Description: "",		//TODO
 		DiskCategory: ecs.DiskCategoryCloudEfficiency,		//TODO
		Size:     size,			//TODO
		SnapshotId:   "",		//TODO
		ClientToken:  "",		//TODO
	}

	diskId, err := client.CreateDisk(&args)

	if err != nil {
		return apiv1.DiskCID{}, bosherr.WrapErrorf(err, "Creating disk of size '%d'", size)
	}

	diskcid := apiv1.NewDiskCID(diskId)

	// TODO? need Attact?
	// NewAttachDiskMethod(a.runner).AttachDisk(*vmcid, diskcid)

	return diskcid, nil
}
