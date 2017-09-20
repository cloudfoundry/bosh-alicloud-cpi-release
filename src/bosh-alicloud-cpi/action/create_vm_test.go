package action

import "testing"
import "bosh-alicloud-cpi/action"

var createVmArgs = []byte(`
{
    "method": "create_vm",
    "arguments": [
        "be387a69-c5d5-4b94-86c2-978581354b50",
        "m-2zehhdtfg22hq46reabf",
        {
            "ephemeral_disk": {
                "size": 50,
                "type": "cloud_efficiency"
            },
            "image_id": "m-2ze200tcuotb5uk2kol4",
            "instance_name": "test-cc",
            "instance_type": "ecs.n4.small",
            "system_disk": {
                "size": 50,
                "type": "cloud_efficiency"
            }
        },
        {
            "private": {
                "ip": "172.16.0.63",
                "netmask": "255.240.0.0",
                "cloud_properties": {
                    "SecurityGroupId": "sg-2zec8ubi1q5aeo5mqcbb",
                    "VSwitchId": "vsw-2zevwt3w7h5u761o405rd"
                },
                "default": [
                    "dns",
                    "gateway"
                ],
                "dns": [
                    "8.8.8.8"
                ],
                "gateway": "172.16.0.1"
            }
        },
        [],
        {}
    ],
    "context": {
        "director_uuid": "911133bb-7d44-4811-bf8a-b215608bf084"
    }
}
`)

func TestCreateVm(t *testing.T) {
	action.CallTestCase(action.TestConfig, createVmArgs, t)
}