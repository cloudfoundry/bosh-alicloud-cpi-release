# BOSH Alibaba Cloud CPI

This is the [BOSH](http://bosh.io) CPI for Alibaba Cloud developed by 4 Alibaba engineers with great excitement, energy and coffee.

## Usage

### Installation

As simple as you just need type in :

```

make

```

### Configuration

Create a configuration file, for example:

```
{
    "alicloud": {
        "region_id": "cn-beijing",
        "access_key_id": "${ACCESS_KEY_ID}",
        "access_key_secret": "${ACCESS_KEY_SECRET}",
        "regions": [
            {
                "name": "cn-beijing",
                "image_id": "m-2zeggz4i4n2z510ajcvw"
            },
            {
                "name": "cn-hangzhou",
                "image_id": "m-bp1bidv1aeiaynlyhmu9"
            }
        ]
    },
    "actions": {
        "agent": {
            "mbus": "http://mbus:mbus@0.0.0.0:6868",
            "blobstore": {
                "provider": "dav",
                "options": {
                    "endpoint": "http://10.0.0.2:25250",
                    "user": "agent",
                    "password": "agent-password"
                }
            }
        },
        "registry": {
            "user": "admin",
            "password": "admin",
            "protocol": "http",
            "host": "127.0.0.1",
            "port": "25777"
        }
    }
}
```

*For case of unit test, you can set your AccessKeyId and AccessKeySecret in system env*

### Run

Run this CPI with the previously created configuration file, such as `create_vm`:

```
$ echo "{\"method\": \"create_vm\", \"auguments\": []}" | cpi -configFile="/path/to/configuration_file.json"
```

### Options

There are some alibaba cloud specified options.

#### Network Options

These options are specified under `cloud_properties` at the [networks](http://bosh.io/docs/networks.html) section of a BOSH deployment manifest and are only valid for `manual` or `dynamic` networks:


```
SecurityGroupId [String] Indicates the ID of the security group

VSwitchId [String] ID of a new VSwitch. Only VSwitches in the same zone can be changed.
```

#### Resource pool Options

These options are specified under `cloud_properties` at the [resource_pools](http://bosh.io/docs/deployment-basics.html#resource-pools) section of a BOSH deployment manifest:

```
instance_type [string] ECS instances are categorized into multiple specification types. For values: [Alibaba Cloud Instance Type](https://www.alibabacloud.com/help/doc-detail/25378.htm).

zone [string] Zones are physical areas with independent power grids and networks in one region

system_disk [Struct]
   - Size [String]
   - Type [String]
ephemeral_disk [Struct]
   - Size [String]
   - Type [String]
```

### Persistent Disks

These options are specified under `cloud_properties` at the [disk_pools](http://bosh.io/docs/persistent-disks.html#persistent-disk-pool) sections of a BOSH deployment manifest.

```
persistent_disk [String] The size of persistent disk, usually measureed by GB
```


