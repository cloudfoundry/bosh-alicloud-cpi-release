/*
 * Copyright (C) 2017-2017 Alibaba Group Holding Limited
 */
package integration

import (
	"fmt"
	"os"
)

var (
	// provider config
	regionId         = envOrDefault("CPI_REGION", "cn-beijing")
	zoneId           = envOrDefault("CPI_ZONE", "cn-beijing-c")
	registry_address = envOrDefault("REGISTRY_ADDRESS", "172.16.0.3")

	// Configurable defaults
	stemcellId		 = envOrDefault("BOSH_STEMCELL_FILE", "m-2zeggz4i4n2z510ajcvw")
	securityGroupId  = envOrDefault("SECURITY_GROUP_ID", "sg-2ze7qg9qdmt1lt9lgvgt")
	vswitchId        = envOrDefault("VSWITCH_ID", "vsw-2ze1oepoom33cdt6nsk88")

	cfgContent = fmt.Sprintf(`{
    "cloud": {
        "plugin": "alicloud",
        "properties": {
            "alicloud": {
                "region_id": "%v",
                "zone_id": "%v",
                "access_key_id": "${ACCESS_KEY_ID}",
                "access_key_secret": "${ACCESS_KEY_CONFIG}"
            },
            "registry": {
                "user": "registry",
                "password": "2a57f7c0-7726-4e76-43aa-00b10b073229",
                "protocol": "http",
                "address": "%v",
                "port": 6901
            },
            "agent": {
                "ntp": ["0.pool.ntp.org", "1.pool.ntp.org"],
                "mbus": "http://mbus:mbus@0.0.0.0:6868",
                "blobstore": {
                    "provider": "dav",
                    "options": {
                        "endpoint": "http://10.0.0.2:25250",
                        "user": "agent",
                        "password": "agent-password"
                    }
                }
            }
        }
    }
}`, regionId, zoneId, registry_address)
)

func envOrDefault(key, defaultVal string) (val string) {
	if val = os.Getenv(key); val == "" {
		val = defaultVal
	}
	return
}
