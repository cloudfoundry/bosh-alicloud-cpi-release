/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package action

import (
	"bosh-alicloud-cpi/alicloud"

	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
)

type UpdateDiskMethod struct {
	CallContext
	disks alicloud.DiskManager
}

func NewUpdateDiskMethod(cc CallContext, disks alicloud.DiskManager) UpdateDiskMethod {
	return UpdateDiskMethod{cc, disks}
}

// UpdateDisk applies category and/or size changes to an existing disk in-place.
func (a UpdateDiskMethod) UpdateDisk(diskCID apiv1.DiskCID, newSize int, cloudProps apiv1.DiskCloudProps) (apiv1.DiskCID, error) {
	diskCid := diskCID.AsString()

	disk, err := a.disks.GetDisk(diskCid)
	if err != nil {
		return diskCID, bosherr.WrapErrorf(err, "UpdateDisk GetDisk failed %s", diskCid)
	}
	if disk == nil {
		return diskCID, bosherr.Errorf("UpdateDisk disk not found id=%s", diskCid)
	}

	var props DiskInfo
	if err := cloudProps.As(&props); err != nil {
		return diskCID, bosherr.WrapErrorf(err, "UpdateDisk failed to parse cloud_properties for disk %s", diskCid)
	}
	props, err = props.Validate(false)
	if err != nil {
		return diskCID, bosherr.WrapErrorf(err, "UpdateDisk invalid cloud_properties for disk %s", diskCid)
	}

	targetCategory := props.GetCategory()
	currentCategory := alicloud.DiskCategory(disk.Category)

	newSizeGB := ConvertToGB(float64(newSize))
	if newSizeGB < disk.Size {
		return diskCID, bosherr.Errorf("UpdateDisk cannot shrink disk %s: requested %d GB < current %d GB",
			diskCid, newSizeGB, disk.Size)
	}

	// Category change — try in-place via ModifyDiskSpec. Runs before the resize
	// so any subsequent resize operates on the disk in its target category.
	if currentCategory != targetCategory {
		err := a.disks.ModifyDiskCategory(diskCid, targetCategory)
		if err != nil && alicloud.IsCategoryRefusedError(err) {
			return diskCID, alicloud.NewNotSupportedError(
				"UpdateDisk: AliCloud refused in-place category change on disk %s (%s -> %s): %s",
				diskCid, currentCategory, targetCategory, err.Error())
		}
		if err != nil {
			return diskCID, bosherr.WrapErrorf(err, "UpdateDisk ModifyDiskCategory failed for disk %s (%s -> %s)",
				diskCid, currentCategory, targetCategory)
		}

		// Wait for Available after the category change before attempting resize.
		if _, err := a.disks.WaitForDiskStatus(diskCid, alicloud.DiskStatusAvailable); err != nil {
			return diskCID, bosherr.WrapErrorf(err, "UpdateDisk WaitForDiskStatus failed for disk %s after category change", diskCid)
		}
	}

	if newSizeGB > disk.Size {
		if err := a.disks.ResizeDisk(diskCid, newSizeGB); err != nil {
			return diskCID, bosherr.WrapErrorf(err, "UpdateDisk ResizeDisk failed for disk %s", diskCid)
		}
	}

	return diskCID, nil
}
