/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package mock

import (
	"bosh-alicloud-cpi/alicloud"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

const (
	defaultDiskSize = 40
)

type TestContext struct {
	config            alicloud.Config
	Disks             map[string]*ecs.Disk
	Instances         map[string]*ecs.Instance
	Stemcells         map[string]*ecs.Image
	Buckets           map[string]*oss.Bucket
	OssObjects        map[string]string
	Snapshots         map[string]string
	NetworkInterfaces map[string]*ecs.NetworkInterface

	// Flags is a shared failure-injection map for tests. It's a map (reference
	// type) so mutations on the suite-level TestContext are visible to the mocks,
	// which each hold a by-value copy of the struct but share the same map header.
	// Recognized keys: "failGetDisks", "failDiskPath".
	Flags map[string]bool

	// WaitForDiskSpecOpts records the opts passed to the most recent WaitForDiskSpec
	// call, so tests can assert that update_disk passes the expected timeout/interval.
	// Stored as a pointer so the mock's copy of TestContext writes through to the original.
	WaitForDiskSpecOpts *[]time.Duration
}

func NewTestContext(config alicloud.Config) TestContext {
	opts := make([]time.Duration, 0)
	return TestContext{
		config:              config,
		Disks:               make(map[string]*ecs.Disk),
		Instances:           make(map[string]*ecs.Instance),
		Stemcells:           make(map[string]*ecs.Image),
		Buckets:             make(map[string]*oss.Bucket),
		OssObjects:          make(map[string]string),
		Snapshots:           make(map[string]string),
		Flags:               make(map[string]bool),
		WaitForDiskSpecOpts: &opts,
	}
}

func (c TestContext) NewDisk(instCid string) (string, *ecs.Disk) {
	d := ecs.Disk{
		DiskId:     NewDiskId(),
		RegionId:   c.config.OpenApi.GetRegion(""),
		ZoneId:     c.config.OpenApi.GetAvailabilityZone(),
		Size:       defaultDiskSize,
		Status:     string(alicloud.DiskStatusAvailable),
		Category:   string(alicloud.DiskCategoryCloudEfficiency),
		InstanceId: instCid,
	}
	c.Disks[d.DiskId] = &d
	return d.DiskId, &d
}

func (c TestContext) NewNetworkInterface() *ecs.NetworkInterface {
	return &ecs.NetworkInterface{
		NetworkInterfaceId: NewNetworkInterfaceId(),
	}
}

func (c TestContext) NewInstance() (string, *ecs.Instance) {
	i := ecs.Instance{
		InstanceId: NewInstanceId(),
		RegionId:   c.config.OpenApi.GetRegion(""),
		ZoneId:     c.config.OpenApi.GetAvailabilityZone(),
		Status:     string(alicloud.Stopped),
	}
	i.NetworkInterfaces.NetworkInterface = []ecs.NetworkInterface{*c.NewNetworkInterface()}
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
