/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package action

import (
	"testing"
)

var createStemcellArgs = []byte(`
{
    "method": "create_stemcell",
    "arguments": [
        "/var/vcap/data/tmp/director/stemcell20170926-9684-17ncxdz/image",
        {
            "architecture": "x86_64",
            "container_format": "",
            "disk": "50",
            "disk_format": "rawdisk",
            "hypervisor": "kvm",
            "image_id": {
                "cn-beijing": "m-2zehhdtfg22hq46reabf",
                "cn-qingdao": "m-m5e1bdeo0a0ujw2s0ovv",
                "cn-zhangjiakou": "m-8vbdd1mi6hwlddicqt2x"
            },
            "infrastructure": "alicloud",
            "name": "bosh-alicloud-kvm-hubuntu-trusty-go_agent",
            "os_distro": "ubuntu",
            "os_type": "linux",
            "root_device_name": "/dev/vda1",
            "source_url": "",
            "version": 1003
        }
    ],
    "context": {
        "director_uuid": "073eac6e-7a35-4a49-8c42-68988ea16ca7"
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
