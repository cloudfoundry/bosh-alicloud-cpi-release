package alicloud_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestAlicloud(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Alicloud Suite")
}
