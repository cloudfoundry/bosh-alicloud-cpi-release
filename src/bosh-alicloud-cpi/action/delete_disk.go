package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	"alibaba/bosh-alicloud-cpi/src/bosh-alicloud-cpi/alicloud"
)

type DeleteDiskMethod struct {
	runner alicloud.Runner
}

func NewDeleteDiskMethod(runner alicloud.Runner) DeleteDiskMethod {
	return DeleteDiskMethod{runner}
}

func (a DeleteDiskMethod) DeleteDisk(cid apiv1.DiskCID) error {
	client := a.runner.NewClient()

	err := client.DeleteDisk(cid.AsString())
	if err != nil {
		return bosherr.WrapErrorf(err, "DeleteDisk %s failed", cid)
	}

	return nil
}
