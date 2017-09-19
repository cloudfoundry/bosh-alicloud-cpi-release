package alicloud

import (
	"testing"
	"github.com/denverdino/aliyungo/ecs"
)


var testConfig = []byte(`{
  "RegionId" : "cn-beijing",
	"ZoneId" : "",
  "AccessKeyId": "LTAI4CkjuLxk2Adz",
  "AccessKeySecret": "WpPV1F3V1nurTAMbw7oS4SnGEWRyRe",
  "Regions": [
    { "Name": "cn-beijing", "ImageId": "m-2zeggz4i4n2z510ajcvw" },
    { "Name": "cn-hangzhou", "ImageId": "m-bp1bidv1aeiaynlyhmu9" }
  ]
}`)


func TestConfigLoad(t *testing.T) {
	config, err := NewConfigFromBytes(testConfig)
	if err != nil {
		t.Error(err, "NewConfigFromBytes Failed!")
	}
	if config.RegionId == "" {
		t.Error("missing RegionId")
	}
	client := ecs.NewClient(config.AccessKeyId, config.AccessKeySecret)

	regions, err := client.DescribeRegions()

	if err != nil {
		t.Error("Error:", err)
	} else {
		t.Log("Regions", regions)
	}
}