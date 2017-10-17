package action

import (
	"os"
	"bytes"
	"bosh-alicloud-cpi/alicloud"
	"github.com/cppforlife/bosh-cpi-go/rpc"
	"testing"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"encoding/json"
)

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
                "address": "10.0.0.2",
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

type CpiResponse struct {
	Result string		`json:"result"`
	Error CpiError	 	`json:"error"`
	Log string			`json:"log"`
}

type CpiError struct {
	Type string			`json:"type"`
	Message string		`json:"message"`
	OkToRetry bool		`json:"ok_to_retry"`
}

//
// FOR Unit Test
func CallTestCase(testConfig []byte, testInput []byte, t *testing.T) (CpiResponse) {
	config, _ := alicloud.NewConfigFromBytes(testConfig)

	logger := boshlog.NewWriterLogger(boshlog.LevelDebug, os.Stderr)
	runner := alicloud.NewRunner(logger, config)

	cpiFactory := NewFactory(runner)

	reader := bytes.NewReader(testInput)
	output := new(bytes.Buffer)

	cli := rpc.NewFactory(logger).NewCLIWithInOut(reader, output, cpiFactory)
	err := cli.ServeOnce()

	if err != nil {
		t.Errorf("CallTestCase() ServeOnce Failed %s", err)
	}

	var resp CpiResponse
	err = json.Unmarshal(output.Bytes(), &resp)

	if err != nil {
		t.Errorf("ServeOnce Output Failed!")
	}

	if resp.Error.Type != "" {
		t.Errorf("ServeOnce Output a Error: %s", resp.Error)
	}

	return resp
}