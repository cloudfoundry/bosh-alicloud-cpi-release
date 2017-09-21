package alicloud

import (
	"testing"
	"github.com/denverdino/aliyungo/ecs"
	"github.com/denverdino/aliyungo/common"
)

var testConfigV2 = []byte(`
{
    "cloud": {
        "plugin": "alicloud",
        "properties": {
            "alicloud": {
                "region_id": "cn-beijing",
                "zone_id": "cn-beijing-a",
                "access_key_id": "OAtTkdQDHJyYGCmU",
                "access_key_secret": "PzosfNHKsG6sFQKNnxvMjqQNx8x2of"
            },
            "registry": {
                "user": "registry",
                "password": "2a57f7c0-7726-4e76-43aa-00b10b073229",
                "protocol": "http",
                "address": "10.0.0.2",
                "port": "6901"
            },
            "agent": {
                "ntp": "ntp",
                "mbus": "http://mbus:mbus@0.0.0.0:6868",
                "blobstore": {
                    "provider": "dav",
                    "options": {
                        "endpoint": "http://10.0.0.2:25250",
                        "user": "agent",
                        "password": "agent-password"
                    }
                }
            }
        }
    }
}
`)


func TestConfigLoad(t *testing.T) {
	config, err := alicloud.NewConfigFromBytes(testConfigV2)
	t.Log(config)

	if err != nil {
		t.Error(err, "NewConfigFromBytes Failed!")
	}
	if config.OpenApi.RegionId == "" {
		t.Error("missing RegionId")
	}
	client := ecs.NewClient(config.OpenApi.AccessKeyId, config.OpenApi.AccessKeySecret)

	regions, err := client.DescribeZones(common.Region("cn-beijing"))

	if err != nil {
		t.Error("Error:", err)
	} else {
		t.Log("Regions", regions)
	}
}