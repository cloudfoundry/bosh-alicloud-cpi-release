/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package action

import (
	"testing"
)

func TestDeviceTrim(t *testing.T) {
	device := "/dev/xvdc"
	if device[5] == 'x' {
		device = "/dev/" + string(device[6:])
	}

	t.Log(device)
}