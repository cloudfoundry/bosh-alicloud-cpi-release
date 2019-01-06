/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package action

import (
	"bosh-alicloud-cpi/alicloud"

	"github.com/cppforlife/bosh-cpi-go/apiv1"
)

type SetDiskMetadataMethod struct {
	CallContext
	disks     alicloud.DiskManager
	instances alicloud.InstanceManager
}

func NewSetDiskMetadataMethod(cc CallContext, disks alicloud.DiskManager, instances alicloud.InstanceManager) SetDiskMetadataMethod {
	return SetDiskMetadataMethod{cc, disks, instances}
}

func (a SetDiskMetadataMethod) SetDiskMetadata(diskCID apiv1.DiskCID, meta apiv1.DiskMeta) error {
	md, err := convertMetaData(meta)
	if err != nil {
		return a.WrapErrorf(err, "convert meta %v failed", meta)
	}

	diskCid := diskCID.AsString()

	tags := make(map[string]string)
	for k, v := range md {
		if k == "deployment" || k == "director" || k == "index" || k == "instance_group" || k == "job" {
			tk := normalizeTag(k)
			if tk != "" {
				tags[tk] = normalizeTag(v.(string))
			}
		}
	}

	err = a.instances.AddTags(diskCid, tags)
	if err != nil {
		return a.WrapErrorf(err, "AddTags %v to %s failed", tags, diskCid)
	}
	return nil
}
