/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package action

import (
	"bosh-alicloud-cpi/alicloud"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"

	"bytes"
	"os/exec"
	"path"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/cppforlife/bosh-cpi-go/apiv1"
	"github.com/google/uuid"
)

const (
	AlicloudImageNamePrefix    = "stemcell"
	MinImageDiskSize           = 5 //in GB
	OSS_BUCKET_NAME_MAX_LENGTH = 64
	PART_SIZE                  = 5 * 1024 * 1024 // 5MB
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
	OSSBucket   string                 `json:"oss_bucket"`
	OSSObject   string                 `json:"oss_object"`
	Description string                 `json:"description,omitempty"`
	Version     string                 `json:"version"`
	Images      map[string]interface{} `json:"image_id"`
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

// Image size minimum is 5GB. Refer to https://www.alibabacloud.com/help/doc-detail/25542.htm
func (a StemcellProps) GetDiskGB() int {
	if a.diskGB < MinImageDiskSize {
		return MinImageDiskSize
	}
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
		return apiv1.StemcellCID{}, bosherr.WrapErrorf(err, "BadInput for CreateStemcell %s", cloudProps)
	}

	props, err = props.Validate()
	if err != nil {
		return apiv1.StemcellCID{}, bosherr.WrapErrorf(err, "BadInput for CreateStemcell %3v", err)
	}

	switch {
	case len(props.Images) > 0:
		// find stemcell from manifest.MF
		stemcellId, err = props.FindStemcellId(a.Config.OpenApi.GetRegion(""))
	case props.OSSBucket != "" && props.OSSObject != "":
		stemcellId, err = a.CreateFromURL(props)
	default:
		stemcellId, err = a.CreateFromTarball(imagePath, props)
	}

	if err != nil {
		return apiv1.StemcellCID{}, bosherr.WrapErrorf(err, "Importing stemcell from '%s'", imagePath)
	}

	stemcellId, err = a.copyImage(stemcellId, props)
	if err != nil {
		return apiv1.StemcellCID{}, bosherr.WrapErrorf(err, "Copying stemcell from '%s'", imagePath)
	}
	return apiv1.NewStemcellCID(stemcellId), nil
}

func (a CreateStemcellMethod) CreateFromURL(props StemcellProps) (string, error) {
	image, err := a.importImage(props)
	if err != nil {
		return "", bosherr.WrapError(err, "Creating Alicloud Image from URL")
	}

	return image, nil
}

func (a CreateStemcellMethod) importImage(props StemcellProps) (string, error) {
	var device ecs.ImportImageDiskDeviceMapping
	device.Format = strings.ToUpper(props.DiskFormat)
	device.OSSBucket = props.OSSBucket
	device.OSSObject = props.OSSObject
	device.DiskImageSize = strconv.Itoa(props.GetDiskGB())

	args := ecs.CreateImportImageRequest()
	args.ImageName = props.OSSObject
	args.Architecture = getValueOrDefault("Architecture", &props, alicloud.AlicloudDefaultImageArchitecture)
	// OS type valid values: linux and windows
	args.OSType = strings.ToLower(getValueOrDefault("OsType", &props, alicloud.AlicloudDefaultImageOSType))
	// The bionic stemcell should using Other Linux to avoid opening ipv6 setting
	if strings.Contains(props.Name, "-bionic-") {
		args.Platform = "Others Linux"
	} else {
		args.Platform = formatImagePlatform(strings.ToLower(props.OsDistro))
	}
	args.Description = props.Description

	devices := []ecs.ImportImageDiskDeviceMapping{
		device,
	}
	args.DiskDeviceMapping = &devices

	a.Logger.Debug(alicloud.AlicloudImageServiceTag, "Creating Alicloud Image with params: %#v", args)
	imageId, err := a.stemcells.ImportImage(args)
	if err != nil {
		return "", bosherr.WrapError(err, "Failed to create Alicloud Image")
	}

	if err = a.stemcells.WaitForImageReady(imageId); err != nil {
		a.cleanUp(imageId)
		return "", bosherr.WrapError(err, "Failed to create Alicloud Image")
	}

	a.Logger.Debug(alicloud.AlicloudImageServiceTag, "Create Alicloud Image %s success", imageId)
	return imageId, nil
}

func (a CreateStemcellMethod) copyImage(stemcellId string, props StemcellProps) (string, error) {
	encryptImage := a.Config.OpenApi.Encrypted
	if encryptImage == nil || !*encryptImage {
		return stemcellId, nil
	}
	kmsKeyId := a.Config.OpenApi.KmsKeyId
	if kmsKeyId != "" {
		a.Logger.Debug(alicloud.AlicloudImageServiceTag, "Copying Alicloud Image with kms key id %s to encrypt the image.", kmsKeyId)
	} else {
		a.Logger.Debug(alicloud.AlicloudImageServiceTag, "Copying Alicloud Image with default kms key to encrypt the image.")
	}

	args := ecs.CreateCopyImageRequest()
	args.ImageId = stemcellId
	args.RegionId = a.Config.OpenApi.GetRegion("")
	args.DestinationRegionId = a.Config.OpenApi.GetRegion("")
	args.DestinationImageName = fmt.Sprintf("bosh-stemcell-%s-%s", props.Version, uuid.New().String())
	args.DestinationDescription = fmt.Sprintf("Copied from stemcell %s:%s", props.Name, props.Version)
	args.Encrypted = requests.NewBoolean(true)
	args.KMSKeyId = kmsKeyId

	imageId, err := a.stemcells.CopyImage(args)
	if err != nil {
		return "", bosherr.WrapError(err, "Failed to create Alicloud Image")
	}

	if err = a.stemcells.WaitForImageReady(imageId); err != nil {
		a.cleanUp(imageId)
		return "", bosherr.WrapError(err, "Failed to copy Alicloud Image")
	}

	a.Logger.Debug(alicloud.AlicloudImageServiceTag, "Copy Alicloud Image %s success", imageId)
	return imageId, nil
}

func (a CreateStemcellMethod) CreateFromTarball(imagePath string, props StemcellProps) (string, error) {
	imageName := fmt.Sprintf("%s-%s.raw", AlicloudImageNamePrefix, a.getUUIDName(props))
	bucketName := fmt.Sprintf("%s-%s", alicloud.AlicloudDefaultImageName, uuid.New().String())

	if len(bucketName) > OSS_BUCKET_NAME_MAX_LENGTH {
		bucketName = bucketName[0:OSS_BUCKET_NAME_MAX_LENGTH]
	}

	var out bytes.Buffer
	var stderr bytes.Buffer

	defer a.osses.DeleteBucket(bucketName)

	if err := a.osses.CreateBucket(bucketName, oss.ACL(oss.ACLPrivate)); err != nil {
		return "", bosherr.WrapError(err, "Creating Alicloud OSS Bucket")
	}

	bucket, err := a.osses.GetBucket(bucketName)

	if err != nil {
		return "", bosherr.WrapError(err, "Geting oss bucket")
	}

	cmd := exec.Command("tar", "-xf", imagePath)
	cmd.Dir = path.Dir(imagePath)
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err = cmd.Run()

	if err != nil {
		return "", bosherr.WrapErrorf(err, fmt.Sprintf("%s-(%s)-(%s)", "Unable to extract image", out.String(), stderr.String()))
	}
	// The root stemcell is 3 GB and using multipart uploading to avoid timeout error
	err = a.osses.MultipartUploadFile(*bucket, imageName, fmt.Sprintf("%s/%s", path.Dir(imagePath), "root.img"), PART_SIZE, oss.Routines(5))
	if err != nil {
		return "", bosherr.WrapError(err, "Uploading stemcell image file to oss")
	}
	defer a.osses.DeleteObject(*bucket, imageName)

	props.OSSBucket = bucketName
	props.OSSObject = imageName
	image, err := a.importImage(props)
	if err != nil {
		return "", bosherr.WrapError(err, "Creating Alicloud Image from Tarball")
	}

	return image, err
}

// image name should be unique and it comes from full stemcell and random suffix length
func (a CreateStemcellMethod) getUUIDName(props StemcellProps) string {
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

// Convert bosh stemcell platform to alibaba cloud api image platform format
// https://www.alibabacloud.com/help/doc-detail/25542.htm
func formatImagePlatform(platform string) string {
	switch platform {
	case "ubuntu":
		return "Ubuntu"
	case "centos":
		return "CentOS"
	case "opensuse":
		return "OpenSUSE"
	default:
		return "Others Linux"
	}
}
