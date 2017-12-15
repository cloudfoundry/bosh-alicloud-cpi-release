/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package mock

import (
	"math/rand"
)

var randomChars = "0123456789abcdefghijklmnopqrstuvwxyz"

func RandomString(length int) (string) {
	s := ""
	for i := 0; i < length; i++ {
		s = s + string(randomChars[rand.Intn(len(randomChars))])
	}
	return s
}

func NewDiskId() (string) {
	return "d-" + RandomString(20)
}

func NewInstanceId() (string) {
	return "i-" + RandomString(20)
}

func NewStemcellId() (string) {
	return "m-" + RandomString(20)
}

func NewOssBucketName() (string) {
	return "oss-" + RandomString(10)
}

func NewOssObejctName() (string) {
	return "obj-" + RandomString(10)
}

func NewVswitchId() (string) {
	return "vsw-" + RandomString(20)
}

func NewSnapshotId() (string) {
	return "s-" + RandomString(20)
}