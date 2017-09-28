package action

import (
	"bosh-alicloud-cpi/alicloud"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	"github.com/denverdino/aliyungo/ecs"
)

type AttachDiskMethod struct {
	runner alicloud.Runner
}

func NewAttachDiskMethod(runner alicloud.Runner) AttachDiskMethod {
	return AttachDiskMethod{runner}
}

func (a AttachDiskMethod) AttachDisk(vmcid apiv1.VMCID, diskCID apiv1.DiskCID) error {

	client := a.runner.NewClient()

	var args ecs.AttachDiskArgs

	args.InstanceId = vmcid.AsString()
	args.DiskId = diskCID.AsString()

	err := client.AttachDisk(&args)

	if err != nil {
		return bosherr.WrapErrorf(err, "Attaching disk '%s' to VM '%s'", diskCID, vmcid)
	}

	// client.DescribeDisks()
	registryClient := a.runner.GetHttpRegistryClient()
	agentSettings, _ := registryClient.Fetch(args.InstanceId)	//TODO, process error
	agentSettings.AttachPersistentDisk(vmcid.AsString(), diskCID.AsString(), "dev/vdc") 	//TODO temp solution
	err = registryClient.Update(vmcid.AsString(), agentSettings)
	if err != nil {
		return bosherr.WrapErrorf(err, "UpdateRegistry failed %s %s", diskCID, vmcid)
	}

	err = a.runner.WaitForDiskStatus(args.DiskId, ecs.DiskStatusAvailable)

	if err != nil {
		return bosherr.WrapErrorf(err, "WaitForDiskStatus failed %s", diskCID)
	}
	return nil
}

