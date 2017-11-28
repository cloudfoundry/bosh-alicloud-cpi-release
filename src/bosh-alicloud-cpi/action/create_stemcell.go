/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package action

import (
	"bosh-alicloud-cpi/alicloud"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/google/uuid"
	"strings"
	"github.com/denverdino/aliyungo/ecs"
	"reflect"
	"fmt"
)

const alicloudImageNamePrefix = "stemcell"

type StemcellProps struct {
	Architecture    string `json:"architecture"`
	ContainerFormat string `json:"container_format"`
	Disk            string `json:"disk"`
	DiskFormat      string `json:"disk_format"`
	Hypervisor      string `json:"hypervisor"`
	Name            string `json:"name"`
	OsDistro        string `json:"os_distro"`
	OsType          string `json:"os_type"`
	//RootDeviceName string 	`json:"root_device_name"`
	SourceUrl string `json:"source_url"`
	//SourceSha1    string `json:"raw_disk_sha1,omitempty"`
	OSSBucket   string `json:"oss_bucket"`
	OSSObject   string `json:"oss_object"`
	Description string `json:"description,omitempty"`
	//	Version string 			`json:"version"`		TODO  sometimes string, and sometimes int
	Images map[string]interface{} `json:"image_id"`
}

type CreateStemcellMethod struct {
	CallContext
	stemcells alicloud.StemcellManager
}

func NewCreateStemcellMethod(cc CallContext, stemcells alicloud.StemcellManager) CreateStemcellMethod {
	return CreateStemcellMethod{cc, stemcells}
}

// create stemcell policy:
// 1. stemcell metadate set in cloudProps
// 2. we provide three ways to create stemcell
//    2.1 region-imageId format, return Stemcell directly
//    2.2 source URL, put the raw image on OSS, then create image by the oss url
//    2.3 local tarball, upload the tarball on OSS, the create image by the oss url
func (a CreateStemcellMethod) CreateStemcell(imagePath string, cloudProps apiv1.StemcellCloudProps) (apiv1.StemcellCID, error) {
	var props StemcellProps
	var stemcellId string
	err := cloudProps.As(&props)

	if err != nil {
		return apiv1.StemcellCID{}, a.WrapErrorf(err, "BadInput for CreateStemcell %s", cloudProps)
	}

	switch {
	case len(props.Images) > 0:
		// find stemcell from manifest.MF
		stemcellId, err = props.FindStemcellId(a.Config.OpenApi.RegionId)
	case props.OSSBucket != "" && props.OSSObject != "":
		stemcellId, err = a.CreateFromURL(props)
	default:
		// todo: support create image from tarball
		return apiv1.StemcellCID{}, a.WrapErrorf(err, "Do not support create image from local tarball '%s'", imagePath)
	}

	if err != nil {
		return apiv1.StemcellCID{}, a.WrapErrorf(err, "Importing stemcell from '%s'", imagePath)
	}

	return apiv1.NewStemcellCID(stemcellId), nil
}

func (a CreateStemcellMethod) CreateFromURL(props StemcellProps) (string, error) {
	image, err := a.importImage(props)
	if err != nil {
		return "", bosherr.WrapErrorf(err, "Creating Alicloud Image from URL")
	}

	return image, nil
}

func (a CreateStemcellMethod) importImage(props StemcellProps) (string, error) {
	var device ecs.DiskDeviceMapping
	device.Format = string(props.DiskFormat)
	device.OSSBucket = props.OSSBucket
	device.OSSObject = props.OSSObject
	device.DiskImageSize = props.Disk

	var args ecs.ImportImageArgs
	args.RegionId = a.Config.OpenApi.GetRegion()
	args.ImageName = a.getUUIDName(props)
	args.Architecture = getValueOrDefault("Architecture", &props, alicloud.AlicloudDefaultImageArchitecture)
	args.OSType = getValueOrDefault("OsType", &props, alicloud.AlicloudDefaultImageOSType)
	args.Platform = props.OsDistro
	args.Description = props.Description

	args.DiskDeviceMappings.DiskDeviceMapping = []ecs.DiskDeviceMapping{
		device,
	}

	a.Logger.Debug(alicloud.AlicloudImageServiceTag, "Creating Alicloud Image with params: %#v", args)
	imageId, err := a.stemcells.ImportImage(args)
	if err != nil {
		return "", bosherr.WrapErrorf(err, "Failed to create Alicloud Image")
	}

	if err = a.stemcells.WaitForImageReady(imageId); err != nil {
		a.cleanUp(imageId)
		return "", bosherr.WrapErrorf(err, "Failed to create Alicloud Image")
	}

	a.Logger.Debug(alicloud.AlicloudImageServiceTag, "Create Alicloud Image %s success", imageId)
	return imageId, nil
}

// image name should be unique
func (a CreateStemcellMethod) getUUIDName(props StemcellProps) (string) {
	uuidStr := uuid.New().String()
	name := getValueOrDefault("Name", &props, alicloud.AlicloudDefaultImageName)
	imageName := fmt.Sprintf("%s-%s", name, uuidStr)
	return imageName
}

func getValueOrDefault(key string, v *StemcellProps, defaultVal string) (val string) {
	r := reflect.ValueOf(v)
	f := reflect.Indirect(r).FieldByName(key)
	return string(f.String())
}

func (a CreateStemcellMethod) cleanUp(id string) {
	if err := a.stemcells.DeleteStemcell(id); err != nil {
		a.Logger.Debug(alicloud.AlicloudImageServiceTag, "Failed cleaning up Alicloud Image '%s': %#v", id, err)
	}
}

func (a StemcellProps) FindStemcellId(region string) (string, error) {
	for k, v := range a.Images {
		if strings.Compare(region, string(k)) == 0 {
			s, ok := v.(string)
			if !ok {
				return "", bosherr.Errorf("Type Failed!")
			}
			return s, nil
		}
	}
	return "", bosherr.Errorf("Can't find stemcell for Region: %s", region)
}
