package action

import (
	"os"
	"bytes"
	"bosh-alicloud-cpi/alicloud"
	"github.com/cppforlife/bosh-cpi-go/rpc"
	"testing"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

var TestConfig = []byte(`
{
    "cloud": {
        "plugin": "alicloud",
        "properties": {
            "alicloud": {
                "region_id": "cn-beijing",
                "zone_id": "cn-beijing-a",
                "access_key_id": "LTAI4CkjuLxk2Adz",
                "access_key_secret": "WpPV1F3V1nurTAMbw7oS4SnGEWRyRe"
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