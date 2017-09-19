package action

import (
	"testing"
	"encoding/json"
	"fmt"
	"alibaba/bosh-alicloud-cpi/src/bosh-alicloud-cpi/alicloud"
	"github.com/cppforlife/bosh-cpi-go/rpc"
	"os"
	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	"bytes"
)

var testCreateStemcellInput = []byte(`
{
    "method": "create_stemcell",
    "arguments": [
        "/var/vcap/data/tmp/director/stemcell20170918-7815-se258c/image",
        {
            "architecture": "x86_64",
            "container_format": "",
            "disk": "50",
            "disk_format": "rawdisk",
            "hypervisor": "kvm",
            "image_id": {
                "cn-beijing": "m-2zehhdtfg22hq46reabf ",
                "cn-hangzhou": "m-bp1bidv1aeiaynlyhmu9"
            },
            "infrastructure": "alicloud",
            "name": "bosh-alicloud-kvm-hubuntu-trusty-go_agent",
            "os_distro": "ubuntu",
            "os_type": "linux",
            "root_device_name": "/dev/vda1",
            "source_url": "",
            "version": "1001"
        }
    ],
    "context": {
        "director_uuid": "911133bb-7d44-4811-bf8a-b215608bf084"
    }
}`)


var testConfig = []byte(`{
"RegionId": "cn-beijing",
  "AccessKeyId": "YourKeyId",
  "AccessKeySecret": "YouSecret",
  "Regions": [
    { "Name": "cn-beijing", "ImageId": "m-2zeggz4i4n2z510ajcvw" },
    { "Name": "cn-hangzhou", "ImageId": "m-bp1bidv1aeiaynlyhmu9" }
  ]
}`)


func TestConfigLoad(t *testing.T) {

}

func TestCreateStemcell(t *testing.T) {
	var config alicloud.AlicloudConfig
	json.Unmarshal(testConfig, &config)
	fmt.Println(config)
	t.Log(config)

	logger := boshlog.NewWriterLogger(boshlog.LevelDebug, os.Stderr)

	cpiFactory := NewFactory(config)

	reader := bytes.NewReader(testCreateStemcellInput)

	cli := rpc.NewFactory(logger).NewCLIWithInOut(reader, os.Stdout, cpiFactory)

	err := cli.ServeOnce()

	if err != nil {
		logger.Error("main", "Serving once %s", err)
		os.Exit(1)
	}
}
//
//wardenConn := wrdnconn.New(config.Warden.ConnectNetwork, config.Warden.ConnectAddress)
//wardenClient := wrdnclient.New(wardenConn)
//
//cpiFactory := bwcaction.NewFactory(
//wardenClient, fs, cmdRunner, uuidGen, config.Actions, logger)
//
//cli := rpc.NewFactory(logger).NewCLI(cpiFactory)
//
//err = cli.ServeOnce()
//if err != nil {
//logger.Error("main", "Serving once %s", err)
//os.Exit(1)
//}
