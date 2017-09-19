package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	"bosh-alicloud-cpi/alicloud"
	"github.com/denverdino/aliyungo/ecs"
)

type DetachDiskMethod struct {
	runner alicloud.Runner
}

func NewDetachDiskMethod(runner alicloud.Runner) DetachDiskMethod {
	return DetachDiskMethod{runner}
}

func (a DetachDiskMethod) DetachDisk(vmCID apiv1.VMCID, diskCID apiv1.DiskCID) error {
	client := a.runner.NewClient()

	var args ecs.DetachDiskArgs

	args.InstanceId = vmCID.AsString()
	args.DiskId = diskCID.AsString()

	err := client.DetachDisk(args.InstanceId, args.DiskId)

	if err != nil {
		return bosherr.WrapErrorf(err, "Attaching disk '%s' to VM '%s'", diskCID, vmCID)
	}

	return err
}
