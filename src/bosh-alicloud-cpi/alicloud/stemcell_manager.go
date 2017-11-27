/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package alicloud

import (
	"github.com/denverdino/aliyungo/ecs"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"github.com/google/uuid"
	"encoding/json"
)

type StemcellManager interface {
	//	FindStemcellId() (string, error)
	//	DeleteStemcell() (string, error)
	ImportImage(string, args ecs.ImportImageArgs) (string, error)
}

//func (a Runner) FindStemcellId() (string, error) {
//	c := a.Config
//	for _, region := range c.OpenApi.Regions {
//		if strings.Compare(region.Name, c.OpenApi.RegionId) == 0 {
//			return region.ImageId, nil
//		}
//	}
//	return "", fmt.Errorf("Unknown Region")
//}

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

func (a StemcellManagerImpl) ImportImage(string, args ecs.ImportImageArgs) (string, error) {
	client := a.config.NewEcsClient()

	args.RegionId = a.config.OpenApi.GetRegion()
	args.ClientToken = uuid.New().String()

	var imageId string
	imageId, err := client.ImportImage(&args)
	a.log("CreateInstance", err, args, imageId)
	return imageId, err
}
