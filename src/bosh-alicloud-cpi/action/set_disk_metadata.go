/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package action

import (
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	"bosh-alicloud-cpi/alicloud"
	"fmt"
)

type SetDiskMetadataMethod struct {
	CallContext
	disks alicloud.DiskManager
}

func NewSetDiskMetadataMethod(cc CallContext, disks alicloud.DiskManager) SetDiskMetadataMethod {
	return SetDiskMetadataMethod{cc, disks}
}

func (a SetDiskMetadataMethod) SetDiskMetadata(diskCID apiv1.DiskCID, meta apiv1.DiskMeta) error {
	md, err := convertMetaData(meta)
	if err != nil {
		return a.WrapErrorf(err, "convert meta %v failed", meta)
	}

	diskCid := diskCID.AsString()

	name := ""
	if s, ok := md["instance_name"]; ok {
		name = s.(string)
	}
	if s, ok := md["instance_index"]; ok {
		name = name + "_" + fmt.Sprintf("%v", s)
	}

	desc := ""
	if s, ok := md["director"]; ok {
		desc = "director: " + s.(string) + "\n"
	}
	if s, ok := md["deployment"]; ok {
		desc = "deployment: " + s.(string) + "\n"
	}
	if s, ok := md["job"]; ok {
		desc = "job: " + s.(string) + "\n"
	}

	name = normalizeName(name, "d_")
	err = a.disks.ModifyDiskAttribute(diskCid, name, desc)
	if err != nil {
		return a.WrapErrorf(err, "modify disk %s attribute failed", diskCid)
	}
	return nil
}