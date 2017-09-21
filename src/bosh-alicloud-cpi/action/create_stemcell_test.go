package action

import (
	"testing"
)

var createStemcellArgs = []byte(`
{
    "method": "create_stemcell",
    "arguments": [
        "/Users/ailan/.bosh/installations/7a96cc73-17ac-40f4-5de2-49ecf08d968f/tmp/stemcell-manager653505843/image",
        {
            "architecture": "x86_64",
            "container_format": "",
            "disk": "50",
            "disk_format": "rawdisk",
            "hypervisor": "kvm",
            "image_id": {
                "cn-beijing": "m-2zeggz4i4n2z510ajcvw",
                "cn-zhangjiakou": "m-8vbcsgb8bmh4iya739i8"
            },
            "infrastructure": "alicloud",
            "name": "bosh-alicloud-kvm-hubuntu-trusty-go_agent",
            "os_distro": "ubuntu",
            "os_type": "linux",
            "root_device_name": "/dev/vda1",
            "source_url": "",
            "version": "1002"
        }
    ],
    "context": {
        "director_uuid": "0670ad54-66b8-4055-57cc-788d48f0886a"
    }
}
`)

func TestConfigLoad(t *testing.T) {
}

func TestCreateStemcell(t *testing.T) {
	CallTestCase(TestConfig, createStemcellArgs, t)
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
