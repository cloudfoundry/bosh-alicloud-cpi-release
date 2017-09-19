package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cppforlife/bosh-cpi-go/apiv1"

	"alibaba/bosh-alicloud-cpi/src/bosh-alicloud-cpi/alicloud"
)

type HasDiskMethod struct {
	runner alicloud.Runner
}

func NewHasDiskMethod(runner alicloud.Runner) HasDiskMethod {
	return HasDiskMethod{runner}
}

func (a HasDiskMethod) HasDisk(cid apiv1.DiskCID) (bool, error) {
	//disk, err := a.diskFinder.Find(cid)
	//if err != nil {
	//	return false, bosherr.WrapErrorf(err, "Finding disk '%s'", cid)
	//}
	//
	//return disk.Exists()
	return false, bosherr.Error("NOT IMPLEMENTED")
}
