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
// Category upgrades use ModifyDiskSpec (upgrades only), size changes use ResizeDisk;
// when combined, the category change runs first. The original DiskCID is always returned.
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

	// Category change — use ModifyDiskSpec (in-place, no snapshot needed on AliCloud).
	// Runs first so a subsequent resize operates on the disk in its target category.
	if currentCategory != targetCategory {
		if err := a.disks.ModifyDiskCategory(diskCid, targetCategory); err != nil {
			return diskCID, bosherr.WrapErrorf(err, "UpdateDisk ModifyDiskCategory failed for disk %s (%s -> %s)",
				diskCid, currentCategory, targetCategory)
		}

		// Wait for the disk to return to Available after the category change; ResizeDisk
		// (and any subsequent director-side attach) requires the disk out of Modifying.
		if _, err := a.disks.WaitForDiskStatus(diskCid, alicloud.DiskStatusAvailable); err != nil {
			return diskCID, bosherr.WrapErrorf(err, "UpdateDisk WaitForDiskStatus failed for disk %s after category change", diskCid)
		}
	}

	// Resize if the requested size is larger than the current disk size. Runs after the
	// category change (if any) so both mutations are applied when combined in one call.
	newSizeGB := ConvertToGB(float64(newSize))
	if newSizeGB > disk.Size {
		if err := a.disks.ResizeDisk(diskCid, newSizeGB); err != nil {
			return diskCID, bosherr.WrapErrorf(err, "UpdateDisk ResizeDisk failed for disk %s", diskCid)
		}
	}

	return diskCID, nil
}
