package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cppforlife/bosh-cpi-go/apiv1"

	"bosh-alicloud-cpi/alicloud"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/denverdino/aliyungo/common"
	"go/types"
)

type HasDiskMethod struct {
	runner alicloud.Runner
}

func NewHasDiskMethod(runner alicloud.Runner) HasDiskMethod {
	return HasDiskMethod{runner}
}

func (a HasDiskMethod) HasDisk(cid apiv1.DiskCID) (bool, error) {
	client := a.runner.NewClient()
	instid := cid.AsString()

	var args ecs.DescribeDisksArgs
	args.RegionId = common.Region(a.runner.Config.RegionId)
	args.DiskIds = []string {cid.AsString()}

	disks, _, err := client.DescribeDisks(&args)

	if err != nil {
		return false, bosherr.WrapErrorf(err, "DescribeDisks failed cid=%s", instid)
	}

	if len(disks) == 0 {
		return false, nil
	} else {
		return true, nil
	}
}
