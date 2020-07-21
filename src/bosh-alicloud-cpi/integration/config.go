/*
 * Copyright (C) 2017-2019 Alibaba Group Holding Limited
 */
package integration

import (
	"bosh-alicloud-cpi/alicloud"
	"encoding/json"
	"fmt"
	"os"
)

var (
	//
	// OpenApi config
	regionId        = envOrDefault("CPI_REGION", "cn-beijing")
	zoneId          = envOrDefault("CPI_ZONE", "cn-beijing-a")
	endpoint        = envOrDefault("CPI_ENDPOINT", "aliyuncs.com")
	accessKeyId     = os.Getenv("CPI_ACCESS_KEY_ID")
	accessKeySecret = os.Getenv("CPI_ACCESS_KEY_SECRET")

	// Configurable defaults
	// A stemcell that will be created in integration_suite_test.go
	existingStemcell string
	stemcellFile     = envOrDefault("CPI_STEMCELL_FILE", "")
	stemcellVersion  = envOrDefault("CPI_STEMCELL_VERSION", "")
	imageOssBucket   = envOrDefault("CPI_STEMCELL_OSS_BUCKET", "")
	imageOssObject   = envOrDefault("CPI_STEMCELL_OSS_OBJECT", "")

	securityGroupId      = os.Getenv("CPI_SECURITY_GROUP_ID")
	vswitchId            = os.Getenv("CPI_VSWITCH_ID")
	internalIp           = envOrDefault("CPI_INTERNAL_IP", "172.16.0.31")
	internalCidr         = envOrDefault("CPI_INTERNAL_CIDR", "172.16.0.0/24")
	internalNetmask      = envOrDefault("CPI_INTERNAL_NETMASK", "255.240.0.0")
	internalGw           = envOrDefault("CPI_INTERNAL_GW", "172.16.0.1")
	slbId                = envOrDefault("CPI_SLB_ID", "")
	slbWeight            = envOrDefault("CPI_SLB_WEIGHT", "100")
	slbServerGroupId     = envOrDefault("CPI_SLB_SERVER_GROUP_ID", "")
	slbServerGroupPort   = envOrDefault("CPI_SLB_SERVER_GROUP_PORT", "33333")
	slbServerGroupWeight = envOrDefault("CPI_SLB_SERVER_GROUP_WEIGHT", "100")
	externalIp           = envOrDefault("CPI_EXTERNAL_IP", "")

	// spot
	spotStrategy   = envOrDefault("CPI_SPOT_STRATEGY", "SpotWithPriceLimit")
	spotPriceLimit = envOrDefault("CPI_SPOT_PRICE_LIMIT", "0.18")

	// ram
	ramRoleName = envOrDefault("RAM_ROLE_NAME", "DirectorRole")
	//tags
	tags = envOrDefault("Tags", `{ "name": "boshTag", "foo": "bar" }`)

	//
	// registry
	registryUser     = envOrDefault("CPI_REGISTRY_USER", "admin")
	registryPassword = envOrDefault("CPI_REGISTRY_PASSWORD", "admin-password")
	registryHost     = envOrDefault("CPI_REGISTRY_ADDRESS", "172.0.0.1")
	registryPort     = envOrDefault("CPI_REGISTRY_PORT", "25777")

	// Channel that will be used to retrieve IPs to use
	ips chan string
)

func ApplySystemEnv(config *alicloud.Config) error {
	a := &config.OpenApi

	a.Region = regionId
	if a.Region == "" {
		return fmt.Errorf("can't find sysenv: CPI_REGION")
	}

	a.AvailabilityZone = zoneId
	if a.AvailabilityZone == "" {
		return fmt.Errorf("can't find sysenv: CPI_ZONE")
	}

	a.AccessKeyId = accessKeyId
	if a.AccessKeyId == "" {
		return fmt.Errorf("can't find sysenv: CPI_ACCESS_KEY_ID")
	}

	a.AccessKeySecret = accessKeySecret
	if a.AccessKeySecret == "" {
		return fmt.Errorf("can't find sysenv: CPI_ACCESS_KEY_SECRET")
	}

	a.AccessEndpoint = endpoint

	registry := &config.Registry
	registry.User = registryUser
	registry.Password = registryPassword
	registry.Host = registryHost
	registry.Port = json.Number(registryPort)

	return nil
}

func envOrDefault(key, defaultVal string) (val string) {
	if val = os.Getenv(key); val == "" {
		val = defaultVal
	}
	return
}
