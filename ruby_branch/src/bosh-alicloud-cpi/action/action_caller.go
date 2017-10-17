package action

import (
	"os"
	"bytes"
	"bosh-alicloud-cpi/alicloud"
	"github.com/cppforlife/bosh-cpi-go/rpc"
	"testing"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

var TestConfig = []byte(`{
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


//
// FOR Unit Test
func CallTestCase(testConfig []byte, testInput []byte, t *testing.T) {
	config, _ := alicloud.NewConfigFromBytes(testConfig)

	logger := boshlog.NewWriterLogger(boshlog.LevelDebug, os.Stderr)
	runner := alicloud.NewRunner(logger, config)

	cpiFactory := NewFactory(runner)

	reader := bytes.NewReader(testInput)

	cli := rpc.NewFactory(logger).NewCLIWithInOut(reader, os.Stdout, cpiFactory)

	err := cli.ServeOnce()

	if err != nil {
		t.Errorf("main", "Serving once %s", err)
	}
}