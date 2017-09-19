package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	"alibaba/bosh-alicloud-cpi/src/bosh-alicloud-cpi/alicloud"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/denverdino/aliyungo/common"
)

type GetDisksMethod struct {
	runner alicloud.Runner
}

func NewGetDisksMethod(runner alicloud.Runner) GetDisksMethod {
	return GetDisksMethod{runner}
}

func (a GetDisksMethod) GetDisks(cid apiv1.VMCID) ([]apiv1.DiskCID, error) {

	client := a.runner.NewClient()
	instid := cid.AsString()

	var args ecs.DescribeDisksArgs
	args.RegionId = common.Region(a.runner.Config.RegionId)
	args.InstanceId = instid

	disks, _, err := client.DescribeDisks(&args)

	if err != nil {
		return nil, bosherr.WrapErrorf(err, "DescribeDisks failed cid=%s", instid)
	}

	var results []apiv1.DiskCID
	for i, v := range disks {
		results[i] = apiv1.NewDiskCID(v.DiskId)
	}

	//TODO? go array not IndexOutOfRange

	return results, nil
}
