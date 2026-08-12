/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package action

import (
	"bosh-alicloud-cpi/alicloud"

	"github.com/cloudfoundry/bosh-cpi-go/apiv1"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
)

type UpdateDiskMethod struct {
	CallContext
	disks alicloud.DiskManager
}

func NewUpdateDiskMethod(cc CallContext, disks alicloud.DiskManager) UpdateDiskMethod {
	return UpdateDiskMethod{cc, disks}
}

// UpdateDisk applies category, PL and/or size changes to a disk in-place.
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

	targetPL := props.GetPerformanceLevel()
	currentPL := disk.PerformanceLevel

	newSizeGB := ConvertToGB(float64(newSize))
	if newSizeGB < disk.Size {
		// AliCloud can't shrink in-place. Return NotSupported (not a generic error) so
		// the director falls back to its copy/migrate path, as it did before update_disk.
		return diskCID, alicloud.NewNotSupportedError(
			"UpdateDisk cannot shrink disk %s in-place: requested %d GB < current %d GB",
			diskCid, newSizeGB, disk.Size)
	}

	// Apply category/PL change in-place before resize so the resize runs on the
	// target spec. PL only changes when it differs; empty target PL leaves it as-is.
	categoryChanged := currentCategory != targetCategory
	plChanged := targetPL != "" && targetPL != currentPL
	if categoryChanged || plChanged {
		err := a.disks.ModifyDiskCategory(diskCid, targetCategory, targetPL)
		if err != nil && alicloud.IsCategoryRefusedError(err) {
			return diskCID, alicloud.NewNotSupportedError(
				"UpdateDisk: AliCloud refused in-place spec change on disk %s (category %s -> %s, PL %s -> %s): %s",
				diskCid, currentCategory, targetCategory, currentPL, targetPL, err.Error())
		}
		if err != nil {
			return diskCID, bosherr.WrapErrorf(err, "UpdateDisk ModifyDiskCategory failed for disk %s (category %s -> %s, PL %s -> %s)",
				diskCid, currentCategory, targetCategory, currentPL, targetPL)
		}

		// Wait for Available before resizing.
		if _, err := a.disks.WaitForDiskStatus(diskCid, alicloud.DiskStatusAvailable); err != nil {
			return diskCID, bosherr.WrapErrorf(err, "UpdateDisk WaitForDiskStatus failed for disk %s after spec change", diskCid)
		}
	}

	if newSizeGB > disk.Size {
		if err := a.disks.ResizeDisk(diskCid, newSizeGB); err != nil {
			return diskCID, bosherr.WrapErrorf(err, "UpdateDisk ResizeDisk failed for disk %s", diskCid)
		}
	}

	return diskCID, nil
}
