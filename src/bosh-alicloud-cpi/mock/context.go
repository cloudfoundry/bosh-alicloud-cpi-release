/*
 * Copyright (C) 2017-2018 Alibaba Group Holding Limited
 */
package mock

import (
	"bosh-alicloud-cpi/alicloud"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

const (
	defaultDiskSize = 40
)

type TestContext struct {
	config     alicloud.Config
	Disks      map[string]*ecs.Disk
	Instances  map[string]*ecs.Instance
	Stemcells  map[string]*ecs.Image
	Buckets    map[string]*oss.Bucket
	OssObjects map[string]string
	Snapshots  map[string]string
}

func NewTestContext(config alicloud.Config) TestContext {
	return TestContext{
		config:     config,
		Disks:      make(map[string]*ecs.Disk),
		Instances:  make(map[string]*ecs.Instance),
		Stemcells:  make(map[string]*ecs.Image),
		Buckets:    make(map[string]*oss.Bucket),
		OssObjects: make(map[string]string),
		Snapshots:  make(map[string]string),
	}
}

func (c TestContext) NewDisk(instCid string) (string, *ecs.Disk) {
	d := ecs.Disk{
		DiskId:     NewDiskId(),
		RegionId:   c.config.OpenApi.GetRegion(),
		ZoneId:     c.config.OpenApi.ZoneId,
		Size:       defaultDiskSize,
		Status:     alicloud.DiskStatusAvailable,
		Category:   alicloud.DiskCategoryCloudEfficiency,
		InstanceId: instCid,
	}
	c.Disks[d.DiskId] = &d
	return d.DiskId, &d
}

func (c TestContext) NewInstance() (string, *ecs.Instance) {
	i := ecs.Instance{
		InstanceId: NewInstanceId(),
		RegionId:   c.config.OpenApi.GetRegion(),
		ZoneId:     c.config.OpenApi.ZoneId,
		Status:     alicloud.Stopped,
	}
	c.Instances[i.InstanceId] = &i
	return i.InstanceId, &i
}

func (c TestContext) NewStemcell() (string, *ecs.Image) {
	d := ecs.DiskDeviceMapping{
		Format:          string(alicloud.RAW),
		ImportOSSBucket: NewOssBucketName(),
		ImportOSSObject: NewOssObejctName(),
	}
	m := ecs.Image{
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

func (c TestContext) NewObject(name, path string) string {
	c.OssObjects[name] = path
	return name
}

func (c TestContext) NewSnapshot(diskCid string) string {
	ssid := NewSnapshotId()
	c.Snapshots[ssid] = diskCid
	return ssid
}
