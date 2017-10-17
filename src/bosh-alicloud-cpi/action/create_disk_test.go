package action

import "testing"

var createDiskTestArgs = []byte (`{
    "method": "create_disk",
    "arguments": [
        30,
        {},
        "i-2zefl7hfr7yb97ni5skw"
    ],
    "context": {
        "director_uuid": "911133bb-7d44-4811-bf8a-b215608bf084"
    }
}`)


func TestCreateDisk(t *testing.T) {
	CallTestCase(TestConfig, createDiskTestArgs, t)
}
