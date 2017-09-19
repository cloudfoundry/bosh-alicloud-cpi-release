package action

import (
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	"alibaba/bosh-alicloud-cpi/src/bosh-alicloud-cpi/alicloud"
)

type DetachDiskMethod struct {
	config alicloud.AlicloudConfig
}

func NewDetachDiskMethod(config alicloud.AlicloudConfig) DetachDiskMethod {
	return DetachDiskMethod{config: config}
}

func (a DetachDiskMethod) DetachDisk(vmCID apiv1.VMCID, diskCID apiv1.DiskCID) error {
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
	//err = vm.DetachDisk(disk)
	//if err != nil {
	//	return bosherr.WrapErrorf(err, "Detaching disk '%s' to VM '%s'", diskCID, vmCID)
	//}

	return bosherr.Error("NOT IMPLEMENTED")
}
