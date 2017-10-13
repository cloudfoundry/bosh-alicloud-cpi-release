package integration

import (
	"bytes"
	"bosh-alicloud-cpi/alicloud"
	"bosh-alicloud-cpi/action"
	"github.com/cppforlife/bosh-cpi-go/rpc"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"encoding/json"
	"log"
	"os"
)

type CpiResponse struct {
	Result string   `json:"result"`
	Error  CpiError `json:"error"`
	Log    string   `json:"log"`
}

type CpiError struct {
	Type      string `json:"type"`
	Message   string `json:"message"`
	OkToRetry bool   `json:"ok_to_retry"`
}

var TestConfig = []byte(`
{
    "cloud": {
        "plugin": "alicloud",
        "properties": {
            "alicloud": {
                "region_id": "cn-beijing",
                "zone_id": "cn-beijing-c",
                "access_key_id": "${ACCESS_KEY_ID}",
                "access_key_secret": "${ACCESS_KEY_CONFIG}"
            },
            "registry": {
                "user": "registry",
                "password": "2a57f7c0-7726-4e76-43aa-00b10b073229",
                "protocol": "http",
                "address": "172.16.0.3",
                "port": 6901
            },
            "agent": {
                "ntp": ["0.pool.ntp.org", "1.pool.ntp.org"],
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

func execCPI(request string) (CpiResponse, error) {
	var resp CpiResponse
	config, _ := alicloud.NewConfigFromBytes(TestConfig)

	logger := boshlog.NewWriterLogger(boshlog.LevelDebug, os.Stderr)
	runner := alicloud.NewRunner(logger, config)

	cpiFactory := action.NewFactory(runner)

	bs := []byte(request)
	reader := bytes.NewReader(bs)
	output := new(bytes.Buffer)

	cli := rpc.NewFactory(logger).NewCLIWithInOut(reader, output, cpiFactory)
	err := cli.ServeOnce()

	if err != nil {
		log.Printf("CallTestCase() ServeOnce Failed: %#v", err)
		return resp, err
	}

	err = json.Unmarshal(output.Bytes(), &resp)
	if err != nil {
		log.Printf("ServeOnce Output Failed!")
	}

	return resp, err
}
