package alicloud

import (
	"github.com/denverdino/aliyungo/ecs"
	"strings"
	"fmt"
)

type Region struct {
	Name string;
	ImageId string;
}

type AlicloudConfig struct {
	RegionId string
	ZoneId string
	AccessKeyId string
	AccessKeySecret string
	Regions[] Region;
}

func (a AlicloudConfig) NewClient() (* ecs.Client) {
	return ecs.NewClient(a.AccessKeyId, a.AccessKeySecret);
}

func (this AlicloudConfig) FindStemcellId() (string, error) {
	for _, region := range this.Regions {
		if (strings.Compare(region.Name, this.RegionId) == 0) {
			return region.ImageId, nil;
		}
	}
	return "", fmt.Errorf("Unknown Region")
}