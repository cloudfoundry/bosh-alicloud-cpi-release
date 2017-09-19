package alicloud

import (
	"testing"
	"encoding/json"
	"fmt"
)


var testConfig = []byte(`{
  "AccessKeyId": "YourKeyId",
  "AccessKeySecret": "YouSecret",
  "Regions": [
    { "Name": "cn-beijing", "ImageId": "m-2zeggz4i4n2z510ajcvw" },
    { "Name": "cn-hangzhou", "ImageId": "m-bp1bidv1aeiaynlyhmu9" }
  ]
}`)


func TestConfigLoad(t *testing.T) {
	var config AlicloudConfig
	json.Unmarshal(testConfig, &config)
	fmt.Println(config)
	t.Log(config)
}