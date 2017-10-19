/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */

package utils

import (
	"fmt"
	"math/rand"
)

func GenerateHexString(length int) (string) {
	s := ""
	for i := 0; i < length; i++ {
		s = s + fmt.Sprintf("%.2x", rand.Intn(256))
	}
	return s
}

func GenerateDiskId() (string) {
	return ""
}

func GenerateInstanceId() (string) {
	return ""
}

func GenerateStemcellId() (string) {
	return ""
}