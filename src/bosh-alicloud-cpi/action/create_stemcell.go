/*
 * Copyright (C) 2017-2018 Alibaba Group Holding Limited
 */
package action

import (
	"bosh-alicloud-cpi/alicloud"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	"github.com/google/uuid"
)

const (
	alicloudImageNamePrefix    = "stemcell"
	UUID_LENGTH                = 32
	OSS_BUCKET_NAME_MAX_LENGTH = 64
)

type StemcellProps struct {
	Architecture    string      `json:"architecture"`
	ContainerFormat string      `json:"container_format"`
	Disk            interface{} `json:"disk"`
	DiskFormat      string      `json:"disk_format"`
	diskGB          int
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
	osses     alicloud.OssManager
}

func NewCreateStemcellMethod(cc CallContext, stemcells alicloud.StemcellManager, osses alicloud.OssManager) CreateStemcellMethod {
	return CreateStemcellMethod{cc, stemcells, osses}
}

func (a StemcellProps) Validate() (StemcellProps, error) {
	if a.Disk != nil {
		switch a.Disk.(type) {
		case int:
			a.diskGB = ConvertToGB(float64(a.Disk.(int)))
			if a.diskGB < 1 {
				return a, fmt.Errorf("disk size too small %dMB < 1GB", a.Disk.(int))
			}
		case float64:
			a.diskGB = ConvertToGB(a.Disk.(float64))
			if a.diskGB < 1 {
				return a, fmt.Errorf("disk size too small %fMB < 1GB", a.Disk.(float64))
			}
		case string:
			s, _ := a.Disk.(string)
			if strings.Compare("", s) == 0 {
				a.diskGB = 0
			}
			s = strings.Replace(s, "_", "", -1)
			n, err := strconv.Atoi(s)
			if err != nil {
				return a, fmt.Errorf("error Disk.size format %v", a.Disk)
			}
			a.diskGB = ConvertToGB(float64(n))
			if a.diskGB < 1 {
				return a, fmt.Errorf("disk size too small %dMB < 1GB", n)
			}
		default:
			return a, fmt.Errorf("error Disk.size type %v", a.Disk)
		}
	} else {
		a.diskGB = 0
	}

	//maxBucketNameLength := OSS_BUCKET_NAME_MAX_LENGTH - UUID_LENGTH
	//if len(a.Name) >= maxBucketNameLength {
	//	return a, fmt.Errorf("bosh stemcell name max length is %d", maxBucketNameLength)
	//}
	return a, nil
}

func (a StemcellProps) GetDiskGB() int {
	return a.diskGB
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

	props, err = props.Validate()
	if err != nil {
		return apiv1.StemcellCID{}, a.WrapErrorf(err, "BadInput for CreateStemcell %3v", err)
	}

	switch {
	case len(props.Images) > 0:
		// find stemcell from manifest.MF
		stemcellId, err = props.FindStemcellId(a.Config.OpenApi.RegionId)
	case props.OSSBucket != "" && props.OSSObject != "":
		stemcellId, err = a.CreateFromURL(props)
	default:
		stemcellId, err = a.CreateFromTarball(imagePath, props)
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
	var device ecs.ImportImageDiskDeviceMapping
	device.Format = string(props.DiskFormat)
	device.OSSBucket = props.OSSBucket
	device.OSSObject = props.OSSObject
	device.DiskImageSize = strconv.Itoa(props.GetDiskGB())

	args := ecs.CreateImportImageRequest()
	args.RegionId = a.Config.OpenApi.GetRegion()
	args.ImageName = a.getUUIDName(props)
	args.Architecture = getValueOrDefault("Architecture", &props, alicloud.AlicloudDefaultImageArchitecture)
	args.OSType = getValueOrDefault("OsType", &props, alicloud.AlicloudDefaultImageOSType)
	args.Platform = props.OsDistro
	args.Description = props.Description

	devices := []ecs.ImportImageDiskDeviceMapping{
		device,
	}
	args.DiskDeviceMapping = &devices

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

func (a CreateStemcellMethod) CreateFromTarball(imagePath string, props StemcellProps) (string, error) {
	imageName := a.getUUIDName(props)
	if err := a.osses.CreateBucket(imageName, oss.ACL(oss.ACLPublicRead)); err != nil {
		return "", bosherr.WrapErrorf(err, "Creating Alicloud OSS Bucket")
	}
	defer a.osses.DeleteBucket(imageName)

	bucket, err := a.osses.GetBucket(imageName)
	if err != nil {
		return "", bosherr.WrapErrorf(err, "Geting oss bucket")
	}

	objectName := fmt.Sprintf("%s.raw", imageName)

	imageFile, err := a.stemcells.OpenLocalFile(imagePath)
	if err != nil {
		return "", bosherr.WrapErrorf(err, "Reading stemcell image file from local")
	}
	defer imageFile.Close()

	err = a.osses.UploadFile(*bucket, objectName, imagePath, 100*1024, oss.Routines(5))
	if err != nil {
		return "", bosherr.WrapErrorf(err, "Uploading stemcell image file to oss")
	}
	defer a.osses.DeleteObject(*bucket, objectName)

	props.OSSBucket = imageName
	props.OSSObject = objectName
	image, err := a.importImage(props)
	if err != nil {
		return "", bosherr.WrapErrorf(err, "Creating Alicloud Image from Tarball")
	}
	return image, err
}

// image name should be unique
// bucket name max length is 64bit, and random suffix length is 32
// so the user input image name should less than 32bit
func (a CreateStemcellMethod) getUUIDName(props StemcellProps) string {
	uuidStr := uuid.New().String()
	name := getValueOrDefault("Name", &props, alicloud.AlicloudDefaultImageName)
	imageName := fmt.Sprintf("%s-%s", name, uuidStr[0:UUID_LENGTH])
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
