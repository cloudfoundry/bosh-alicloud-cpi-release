/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package integration

import (
	"fmt"
	"os"
	"bosh-alicloud-cpi/alicloud"
	"encoding/json"
)

var (
	//
	// OpenApi config
	regionId        = envOrDefault("CPI_REGION", "cn-beijing")
	zoneId          = envOrDefault("CPI_ZONE", "cn-beijing-a")
	accessKeyId		= os.Getenv("CPI_ACCESS_KEY_ID")
	accessKeySecret = os.Getenv("CPI_ACCESS_KEY_SECRET")

	//
	// for <light-bosh-stemcell-1008-alicloud-kvm-ubuntu-trusty-go_agent> in beijing
	stemcellId		 = envOrDefault("CPI_STEMCELL_ID", "m-2ze2ct08gslmli5e6fw5")

	securityGroupId = os.Getenv("CPI_SECURITY_GROUP_ID")
	vswitchId       = os.Getenv("CPI_VSWITCH_ID")
	internalIp      = envOrDefault("CPI_INTERNAL_IP", "172.16.0.31")
	internalCidr 	= envOrDefault("CPI_INTERNAL_CIDR", "172.16.0.0/24")
	internalNetmask	= envOrDefault("CPI_INTERNAL_NETMASK", "255.240.0.0")
	internalGw  	= envOrDefault("CPI_INTERNAL_GW", "172.16.0.1")
	slbId		  	= envOrDefault("CPI_SLB_ID", "")
	externalIp      = envOrDefault("CPI_EXTERNAL_IP", "")

	//
	// registry
	registryUser 		= envOrDefault("CPI_REGISTRY_USER", "admin")
	registryPassword	= envOrDefault("CPI_REGISTRY_PASSWORD", "admin-password")
	registryHost		= envOrDefault("CPI_REGISTRY_ADDRESS", "172.0.0.1")
	registryPort		= envOrDefault("CPI_REGISTRY_PORT", "25777")
)

func ApplySystemEnv(config *alicloud.Config) (error) {
	a := &config.OpenApi

	a.RegionId = regionId
	if a.RegionId == "" {
		return fmt.Errorf("can't find sysenv: CPI_REGION")
	}

	a.ZoneId = zoneId
	if a.ZoneId == "" {
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
