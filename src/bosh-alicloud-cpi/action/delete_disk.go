package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	"bosh-alicloud-cpi/alicloud"
	"github.com/denverdino/aliyungo/ecs"
)

type DeleteDiskMethod struct {
	runner alicloud.Runner
}

func NewDeleteDiskMethod(runner alicloud.Runner) DeleteDiskMethod {
	return DeleteDiskMethod{runner}
}

func (a DeleteDiskMethod) DeleteDisk(cid apiv1.DiskCID) error {
	client := a.runner.NewClient()

	disk, err := a.runner.GetDiskStatus(cid.AsString())
	if err != nil {
		return bosherr.WrapErrorf(err, "delete_disk/getStatus cid=%s", cid.AsString())
	}

	if disk == nil {
		return nil
	}

	switch disk.Status {
	case ecs.DiskStatusDetaching:
		_, err := a.runner.WaitForDiskStatus(cid.AsString(), ecs.DiskStatusAvailable)
		if err != nil {
			return bosherr.WrapErrorf(err, "delete_disk/wait for detaching failed=%s", cid.AsString())
		}
	case ecs.DiskStatusAvailable:
		break
	default:
		return bosherr.Errorf("delete_disk/unexcepted_DiskStatus: %s-%v", cid.AsString(), disk.Status)
	}

	err = client.DeleteDisk(cid.AsString())
	if err != nil {
		return bosherr.WrapErrorf(err, "DeleteDisk %s failed", cid)
	}

	return nil
}
