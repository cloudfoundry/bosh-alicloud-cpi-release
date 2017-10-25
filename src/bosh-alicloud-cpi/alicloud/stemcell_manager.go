/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package alicloud

type StemcellManager interface {
//	FindStemcellId() (string, error)
//	DeleteStemcell() (string, error)
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
}

func NewStemcellManager(config Config) (StemcellManager) {
	return StemcellManagerImpl{
	}
}