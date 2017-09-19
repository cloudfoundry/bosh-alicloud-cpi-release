package alicloud

import (
	"testing"
	"github.com/denverdino/aliyungo/ecs"
)

var testConfigV2 = []byte(`
{
    "alicloud": {
        "region_id": "cn-beijing",
        "access_key_id": "${ACCESS_KEY_ID}",
        "access_key_secret": "${ACCESS_KEY_SECRET}",
        "regions": [
            {
                "name": "cn-beijing",
                "image_id": "m-2zeggz4i4n2z510ajcvw"
            },
            {
                "name": "cn-hangzhou",
                "image_id": "m-bp1bidv1aeiaynlyhmu9"
            }
        ]
    },
    "actions": {
        "agent": {
            "mbus": "http://mbus:mbus@0.0.0.0:6868",
            "blobstore": {
                "provider": "dav",
                "options": {
                    "endpoint": "http://10.0.0.2:25250",
                    "user": "agent",
                    "password": "agent-password"
                }
            }
        },
        "registry": {
            "user": "admin",
            "password": "admin",
            "protocol": "http",
            "host": "127.0.0.1",
            "port": "25777"
        }
    }
}
`)


func TestConfigLoad(t *testing.T) {
	config, err := NewConfigFromBytes(testConfigV2)
	t.Log(config)

	if err != nil {
		t.Error(err, "NewConfigFromBytes Failed!")
	}
	if config.OpenApi.RegionId == "" {
		t.Error("missing RegionId")
	}
	client := ecs.NewClient(config.OpenApi.AccessKeyId, config.OpenApi.AccessKeySecret)

	regions, err := client.DescribeRegions()

	if err != nil {
		t.Error("Error:", err)
	} else {
		t.Log("Regions", regions)
	}
}