/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package alicloud

import (
	"testing"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAlicloud(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Alicloud Suite")
}
