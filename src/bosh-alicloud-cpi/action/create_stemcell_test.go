package action

import (
	"testing"
)

var createStemcellArgs = []byte(`
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
