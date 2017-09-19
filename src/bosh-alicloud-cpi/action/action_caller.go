package action

import (
	"os"
	"bytes"
	"bosh-alicloud-cpi/alicloud"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"github.com/cppforlife/bosh-cpi-go/rpc"

	"testing"
)

var TestConfig = []byte(`{
  "RegionId" : "cn-beijing",
	"ZoneId" : "",
  "AccessKeyId": "***REMOVED***",
  "AccessKeySecret": "***REMOVED***",
  "Regions": [
    { "Name": "cn-beijing", "ImageId": "m-2zeggz4i4n2z510ajcvw" },
    { "Name": "cn-hangzhou", "ImageId": "m-bp1bidv1aeiaynlyhmu9" }
  ]
}`)


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