package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	"alibaba/bosh-alicloud-cpi/src/bosh-alicloud-cpi/alicloud"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/denverdino/aliyungo/common"
)

type CreateDiskMethod struct {
	config alicloud.AlicloudConfig
}

func NewCreateDiskMethod(config alicloud.AlicloudConfig) CreateDiskMethod {
	return CreateDiskMethod{config: config}
}

func (this CreateDiskMethod) CreateDisk(size int, props apiv1.DiskCloudProps, vmcid *apiv1.VMCID) (apiv1.DiskCID, error) {
	var client = this.config.NewClient();

	var args = ecs.CreateDiskArgs {
		RegionId: common.Region(this.config.RegionId),
		ZoneId: this.config.ZoneId,
		DiskName: "",			//TODO
		Description: "",		//TODO
 		DiskCategory: "",		//TODO
		Size:     size,			//TODO
		SnapshotId:   "",		//TODO
		ClientToken:  "",		//TODO
	}

	diskId, err := client.CreateDisk(&args)

	if err != nil {
		return apiv1.DiskCID{}, bosherr.WrapErrorf(err, "Creating disk of size '%d'", size)
	}

	return apiv1.NewDiskCID(diskId), nil
}
