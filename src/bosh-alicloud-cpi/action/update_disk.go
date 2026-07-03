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

// UpdateDisk changes the category of an existing persistent disk in-place using AliCloud's
// ModifyDiskSpec API. AliCloud supports forward (upgrade) transitions only, e.g.
// cloud_efficiency -> cloud_essd. The disk CID is unchanged after migration.
//
// If the disk already has the requested category, the call is a no-op.
// If only the size changes (same category), the disk is resized in-place.
//
// All AliCloud paths update the disk in place — the returned DiskCID is always the
// original diskCID. A different CID would signal that the disk was replaced (e.g.
// snapshot + recreate), which this implementation never does.
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
	if currentCategory != targetCategory {
		if err := a.disks.ModifyDiskCategory(diskCid, targetCategory); err != nil {
			return diskCID, bosherr.WrapErrorf(err, "UpdateDisk ModifyDiskCategory failed for disk %s (%s -> %s)",
				diskCid, currentCategory, targetCategory)
		}

		// Wait for the disk to return to Available after the category change.
		if _, err := a.disks.WaitForDiskStatus(diskCid, alicloud.DiskStatusAvailable); err != nil {
			return diskCID, bosherr.WrapErrorf(err, "UpdateDisk WaitForDiskStatus failed for disk %s after category change", diskCid)
		}

		return diskCID, nil
	}

	// Same category — resize only if the new size is larger.
	newSizeGB := ConvertToGB(float64(newSize))
	if newSizeGB > disk.Size {
		if err := a.disks.ResizeDisk(diskCid, newSizeGB); err != nil {
			return diskCID, bosherr.WrapErrorf(err, "UpdateDisk ResizeDisk failed for disk %s", diskCid)
		}
	}

	return diskCID, nil
}
