package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	"alibaba/bosh-alicloud-cpi/src/bosh-alicloud-cpi/alicloud"
)

type DeleteDiskMethod struct {
	config alicloud.AlicloudConfig
}

func NewDeleteDiskMethod(config alicloud.AlicloudConfig) DeleteDiskMethod {
	return DeleteDiskMethod{config:config}
}

func (a DeleteDiskMethod) DeleteDisk(cid apiv1.DiskCID) error {
	//disk, err := a.diskFinder.Find(cid)
	//if err != nil {
	//	return bosherr.WrapErrorf(err, "Finding disk '%s'", cid)
	//}
	//
	//err = disk.Delete()
	//if err != nil {
	//	return bosherr.WrapErrorf(err, "Deleting disk '%s'", cid)
	//}

	return bosherr.Error("NOT IMPLEMENTED")
}
