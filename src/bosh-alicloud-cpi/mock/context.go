/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package mock

import (
	"github.com/denverdino/aliyungo/ecs"
	"bosh-alicloud-cpi/alicloud"
)

type TestContext struct {
	config alicloud.Config
	Disks map[string]*ecs.DiskItemType
	Instances map[string]*ecs.InstanceAttributesType
	Snapshots map[string]string
}

func NewTestContext(config alicloud.Config) TestContext {
	return TestContext{
		config: config,
		Disks: make(map[string]*ecs.DiskItemType),
		Instances: make(map[string]*ecs.InstanceAttributesType),
		Snapshots: make(map[string]string),
	}
}

func (c TestContext) NewDisk(instCid string) (string, *ecs.DiskItemType) {
	d := ecs.DiskItemType {
		DiskId:NewDiskId(),
		RegionId:c.config.OpenApi.GetRegion(),
		ZoneId:c.config.OpenApi.ZoneId,
		Status:ecs.DiskStatusAvailable,
		Category:ecs.DiskCategoryCloudEfficiency,
		InstanceId:instCid,
	}
	c.Disks[d.DiskId] = &d
	return d.DiskId, &d
}

func (c TestContext) NewInstance() (string, *ecs.InstanceAttributesType) {
	i := ecs.InstanceAttributesType{
		InstanceId:NewInstanceId(),
		RegionId:c.config.OpenApi.GetRegion(),
		ZoneId:c.config.OpenApi.ZoneId,
		Status: ecs.Stopped,
	}
	c.Instances[i.InstanceId] = &i
	return i.InstanceId, &i
}




