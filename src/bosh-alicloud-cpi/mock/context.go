/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package mock

import (
	"github.com/denverdino/aliyungo/ecs"
	"bosh-alicloud-cpi/alicloud"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

const (
	defaultDiskSize = 40
)

type TestContext struct {
	config     alicloud.Config
	Disks      map[string]*ecs.DiskItemType
	Instances  map[string]*ecs.InstanceAttributesType
	Stemcells  map[string]*ecs.ImageType
	Buckets    map[string]*oss.Bucket
	OssObjects map[string]string
	Snapshots  map[string]string
}

func NewTestContext(config alicloud.Config) TestContext {
	return TestContext{
		config:     config,
		Disks:      make(map[string]*ecs.DiskItemType),
		Instances:  make(map[string]*ecs.InstanceAttributesType),
		Stemcells:  make(map[string]*ecs.ImageType),
		Buckets:    make(map[string]*oss.Bucket),
		OssObjects: make(map[string]string),
		Snapshots:  make(map[string]string),
	}
}

func (c TestContext) NewDisk(instCid string) (string, *ecs.DiskItemType) {
	d := ecs.DiskItemType{
		DiskId:     NewDiskId(),
		RegionId:   c.config.OpenApi.GetRegion(),
		ZoneId:     c.config.OpenApi.ZoneId,
		Size:       defaultDiskSize,
		Status:     ecs.DiskStatusAvailable,
		Category:   ecs.DiskCategoryCloudEfficiency,
		InstanceId: instCid,
	}
	c.Disks[d.DiskId] = &d
	return d.DiskId, &d
}

func (c TestContext) NewInstance() (string, *ecs.InstanceAttributesType) {
	i := ecs.InstanceAttributesType{
		InstanceId: NewInstanceId(),
		RegionId:   c.config.OpenApi.GetRegion(),
		ZoneId:     c.config.OpenApi.ZoneId,
		Status:     ecs.Stopped,
	}
	c.Instances[i.InstanceId] = &i
	return i.InstanceId, &i
}

func (c TestContext) NewStemcell() (string, *ecs.ImageType) {
	d := ecs.DiskDeviceMapping{
		Format:    string(ecs.RAW),
		OSSBucket: NewOssBucketName(),
		OSSObject: NewOssObejctName(),
	}
	m := ecs.ImageType{
		ImageId: NewStemcellId(),
	}
	m.DiskDeviceMappings.DiskDeviceMapping = []ecs.DiskDeviceMapping{
		d,
	}

	c.Stemcells[m.ImageId] = &m
	return m.ImageId, &m
}

func (c TestContext) NewBucket(name string) (string, *oss.Bucket) {
	b := oss.Bucket{
		BucketName: NewOssBucketName(),
	}

	c.Buckets[name] = &b
	return b.BucketName, &b
}

func (c TestContext) NewObject(name, path string) (string) {
	c.OssObjects[name] = path
	return name
}

func (c TestContext) NewSnapshot(diskCid string) (string) {
	ssid := NewSnapshotId()
	c.Snapshots[ssid] = diskCid
	return ssid
}
