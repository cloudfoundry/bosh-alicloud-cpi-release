package action

import (
	"alibaba/bosh-alicloud-cpi/src/bosh-alicloud-cpi/alicloud"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cppforlife/bosh-cpi-go/apiv1"

)

type AttachDiskMethod struct {
	config alicloud.AlicloudConfig
}

func NewAttachDiskMethod(config alicloud.AlicloudConfig) AttachDiskMethod {
	return AttachDiskMethod{config: config}
}

func (a AttachDiskMethod) AttachDisk(vmCID apiv1.VMCID, diskCID apiv1.DiskCID) error {
	//vm, found, err := a.vmFinder.Find(vmCID)
	//if err != nil {
	//	return bosherr.WrapErrorf(err, "Finding VM '%s'", vmCID)
	//}
	//
	//if !found {
	//	return bosherr.Errorf("Expected to find VM '%s'", vmCID)
	//}
	//
	//disk, err := a.diskFinder.Find(diskCID)
	//if err != nil {
	//	return bosherr.WrapErrorf(err, "Finding disk '%s'", diskCID)
	//}
	//
	//err = vm.AttachDisk(disk)
	//if err != nil {
	//	return bosherr.WrapErrorf(err, "Attaching disk '%s' to VM '%s'", diskCID, vmCID)
	//}

	return bosherr.Error("NOT IMPLEMENTED")
}
