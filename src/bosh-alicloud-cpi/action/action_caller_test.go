package action

import (
"testing"
"bosh-alicloud-cpi/alicloud"
"os"
boshlog "github.com/cloudfoundry/bosh-utils/logger"
)

func TestActionCaller(t *testing.T) {
	config, err := alicloud.NewConfigFromBytes(TestConfig)
	if err != nil {
		t.Error(err)
	}

	logger := boshlog.NewWriterLogger(boshlog.LevelDebug, os.Stderr)
	runner := alicloud.NewRunner(logger, config)
	client := runner.NewClient()

	t.Log(client.DescribeRegions())
}