/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package mock

import "github.com/denverdino/aliyungo/ecs"

type TestContext struct {
	Disks map[string]*ecs.DiskItemType
	Instances map[string]*ecs.InstanceAttributesType
}

func NewTestContext() TestContext {
	return TestContext{
		Disks: make(map[string]*ecs.DiskItemType),
		Instances: make(map[string]*ecs.InstanceAttributesType),
	}
}

func (tc TestContext) NewDisk() (string, ecs.DiskItemType) {
	d := ecs.DiskItemType {
		DiskId:NewDiskId(),
		Status:ecs.DiskStatusAvailable,
		Category:ecs.DiskCategoryCloudEfficiency,
		InstanceId:"",
	}
	tc.Disks[d.DiskId] = &d
	return d.DiskId, d
}

func (tc TestContext) NewInstance() (string, ecs.InstanceAttributesType) {
	i := ecs.InstanceAttributesType{
		InstanceId:NewInstanceId(),
		Status: ecs.Stopped,
	}
	tc.Instances[i.InstanceId] = &i
	return i.InstanceId, i
}




