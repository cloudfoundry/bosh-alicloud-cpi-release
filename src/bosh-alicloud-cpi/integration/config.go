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
	"fmt"
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

var (
	// provider config
	regionId         = envOrDefault("BOSH_ALICLOUD_REGION_ID", "cn-beijing")
	zoneId           = envOrDefault("BOSH_ALICLOUD_ZONE_ID", "cn-beijing-c")
	registry_address = envOrDefault("BOSH_ALICLOUD_REGISTRY_ADDRESS", "172.16.0.3")

	// Configurable defaults
	boshStemcellImageId = envOrDefault("BOSH_ALICLOUD_STEMCELL_ID", "m-2zeggz4i4n2z510ajcvw")
	securityGroupId     = envOrDefault("BOSH_ALICLOUD_SECURITY_GROUP_ID", "sg-2ze7qg9qdmt1lt9lgvgt")
	vswitchId           = envOrDefault("BOSH_ALICLOUD_VSWITCH_ID", "vsw-2ze1oepoom33cdt6nsk88")
	externalIp          = envOrDefault("BOSH_ALICLOUD_EXTERNAL_IP", "47.94.216.146")

	cfgContent = fmt.Sprintf(`{
    "cloud": {
        "plugin": "alicloud",
        "properties": {
            "alicloud": {
                "region_id": "%v",
                "zone_id": "%v",
                "access_key_id": "${ACCESS_KEY_ID}",
                "access_key_secret": "${ACCESS_KEY_CONFIG}"
            },
            "registry": {
                "user": "registry",
                "password": "2a57f7c0-7726-4e76-43aa-00b10b073229",
                "protocol": "http",
                "address": "%v",
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
}`, regionId, zoneId, registry_address)
)

func execCPI(request string) (CpiResponse, error) {
	var resp CpiResponse
	config, _ := alicloud.NewConfigFromBytes([]byte(cfgContent))

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

func envOrDefault(key, defaultVal string) (val string) {
	if val = os.Getenv(key); val == "" {
		val = defaultVal
	}
	return
}
