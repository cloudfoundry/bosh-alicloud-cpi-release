/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package action

import (
	"bosh-alicloud-cpi/alicloud"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	"github.com/satori/go.uuid"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"strings"
	"fmt"
	"github.com/denverdino/aliyungo/ecs"
)

const alicloudImageNamePrefix = "stemcell"
const alicloudImageServiceTag = "AlicloudImageService"

type StemcellProps struct {
	Architecture    string `json:"architecture"`
	ContainerFormat string `json:"container_format"`
	//Disk string				`json:"disk"`
	//DiskFormat string 		`json:"disk_format"`
	Hypervisor string `json:"hypervisor"`
	Name       string `json:"name"`
	//OsDistro string			`json:"os_distro"`
	//OsType string 			`json:"os_type"`
	//RootDeviceName string 	`json:"root_device_name"`
	SourceUrl  string `json:"source_url"`
	SourceSha1 string `json:"raw_disk_sha1,omitempty"`
	Format string `json:"format,omitempty"`
	OSSBucket string `json:"oss_bucket"`
	OSSObject string `json:"oss_object"`
	DiskImageSize string `json:"disk_image_size,omitempty"`
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
	case props.SourceUrl != "":
		stemcellId, err = a.CreateFromURL(props)
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
	uuidStr := uuid.NewV4()
	imageName := fmt.Sprintf("%s-%s", alicloudImageNamePrefix, uuidStr)

	var device ecs.DiskDeviceMapping
	device.Format = "RAW"
	device.OSSBucket = ""

	var args ecs.ImportImageArgs

	args.RegionId = a.Config.OpenApi.GetRegion()
	args.ImageName = imageName
	args.Description = props.Name

	a.Logger.Debug(alicloudImageServiceTag, "Creating Alicloud Image with params: %#v", image)
	operation, err := i.computeService.Images.Insert(i.project, image).Do()
	if err != nil {
		return "", bosherr.WrapErrorf(err, "Failed to create Google Image")
	}

	if _, err = i.operationService.Waiter(operation, "", ""); err != nil {
		i.cleanUp(image.Name)
		return "", bosherr.WrapErrorf(err, "Failed to create Google Image")
	}

	return image.Name, nil
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
