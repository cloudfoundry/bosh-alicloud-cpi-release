/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package mock

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestBooks(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Mock Suite")
}