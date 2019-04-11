/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package integration

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("integration:alicloud", func() {
	It("get regions and zones", func() {
		client, _ := caller.Config.NewEcsClient("")
		resp, err := client.DescribeRegions(ecs.CreateDescribeRegionsRequest())
		Expect(err).NotTo(HaveOccurred())

		for _, r := range resp.Regions.Region {
			// region := common.Region(r.RegionId)
			// zones, err := client.DescribeZones(region)
			// Expect(err).NotTo(HaveOccurred())
			GinkgoT().Logf("Alibaba Cloud region [%s]%s zones %v", r.RegionId, r.LocalName)
		}
	})

	It("test for client", func() {
		//client := caller.Config.NewEcsClient()
		//args := ecs.ImportImageArgs{}
		//client.CreateImage()
		//caller.Config.OpenApi
		//client := oss.NewOSSClient(internal)
		//bucket := client.Bucket("Hello")
		//bucket.PutFile()

	})
})
