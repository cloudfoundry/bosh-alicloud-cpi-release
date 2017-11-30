/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package alicloud

import (
	"github.com/denverdino/aliyungo/ecs"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	bosherr "github.com/cloudfoundry/bosh-utils/errors"
	"github.com/google/uuid"
	"encoding/json"
	"os"
)

const (
	AlicloudImageServiceTag          = "AlicloudImageService"
	AlicloudDefaultImageName         = "bosh-stemcell"
	AlicloudDefaultImageArchitecture = "x86_64"
	AlicloudDefaultImageOSType       = "linux"
	WaitForImageReadyTimeout         = 1800
)

type StemcellManager interface {
	FindStemcellById(id string) (*ecs.ImageType, error)
	DeleteStemcell(id string) (error)
	ImportImage(args ecs.ImportImageArgs) (string, error)
	OpenLocalFile(path string) (*os.File, error)
	WaitForImageReady(id string) (error)
}

type StemcellManagerImpl struct {
	config Config
	logger boshlog.Logger
	region string
}

func NewStemcellManager(config Config, logger boshlog.Logger) (StemcellManager) {
	return StemcellManagerImpl{
		config: config,
		logger: logger,
		region: config.OpenApi.RegionId,
	}
}

func (a StemcellManagerImpl) log(action string, err error, args interface{}, result string) {
	s, _ := json.Marshal(args)
	if err != nil {
		a.logger.Error("StemcellManager", "%s failed args=%s err=%s", action, s, err)
	} else {
		a.logger.Info("StemcellManager", "%s done! args=%s result=%s", action, s, result)
	}
}

func (a StemcellManagerImpl) FindStemcellById(id string) (*ecs.ImageType, error) {
	client := a.config.NewEcsClient()
	a.logger.Debug(AlicloudImageServiceTag, "Finding Alicloud Image '%s'", id)

	args := ecs.DescribeImagesArgs{
		RegionId: a.config.OpenApi.GetRegion(),
		ImageId:  id,
	}

	images, _, err := client.DescribeImages(&args)
	a.logger.Debug(AlicloudImageServiceTag, "Find Alicloud Images '%#v'", images)

	if err != nil {
		return nil, err
	}

	if len(images) == 0 {
		return nil, nil
	}

	return &images[0], nil
}

func (a StemcellManagerImpl) DeleteStemcell(id string) (error) {
	client := a.config.NewEcsClient()
	image, err := a.FindStemcellById(id)
	if err != nil {
		return err
	}
	if image == nil {
		return bosherr.WrapErrorf(err, "Alicloud Image '%s' does not exists", id)
	}

	a.logger.Debug(AlicloudImageServiceTag, "Deleting Alicloud Image '%s'", id)
	err = client.DeleteImage(a.config.OpenApi.GetRegion(), id)

	if err != nil {
		return bosherr.WrapErrorf(err, "Failed to delete Alicloud Image '%s'", id)
	}

	return nil
}

func (a StemcellManagerImpl) ImportImage(args ecs.ImportImageArgs) (string, error) {
	client := a.config.NewEcsClient()

	args.RegionId = a.config.OpenApi.GetRegion()
	args.ClientToken = uuid.New().String()

	var imageId string
	imageId, err := client.ImportImage(&args)
	a.log("Importing Image", err, args, imageId)
	return imageId, err
}

func (a StemcellManagerImpl) OpenLocalFile(path string) (*os.File, error) {
	return os.Open(path)
}

// import image from oss may take >=15min
// we set timeout value to 30min, if needed turn it up
func (a StemcellManagerImpl) WaitForImageReady(id string) (error) {
	client := a.config.NewEcsClient()
	region := a.config.OpenApi.GetRegion()
	return client.WaitForImageReady(region, id, WaitForImageReadyTimeout)
}
