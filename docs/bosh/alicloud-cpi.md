---
title: Alicloud CPI
---

This topic describes cloud properties for different resources created by the AWS CPI.

## <a id='azs'></a> AZs

Schema for `cloud_properties` section:

* **availability_zone** [String, required]: Availability zone to use for creating instances. Example: `cn-beijing-a`.

Example:

```yaml
azs:
- name: z1
  cloud_properties:
    availability_zone: cn-beijing-a
```

---
## <a id='networks'></a> Networks

Schema for `cloud_properties` section used by dynamic network or manual network subnet:

* **vswitch_id** [String, required]: VSwitch ID in which the instance will be created. Example: `vsw-2zemyfytfclbcmgfkzokx`.
* **security_group_id** [String, required]: [Security Group](https://www.alibabacloud.com/help/zh/doc-detail/25468.htm), by ID, to apply to all VMs placed on this network. Example: `sg-2zei0mcphxbdxj49qtmz`

Example of manual network:

```yaml
networks:
- name: default
  type: manual
  subnets:
  - range: 10.10.0.0/24
    gateway: 10.10.0.1
    cloud_properties:
      vswitch_id: vsw-2zemyfytfclbcmgfkzokx
      security_group_id: sg-2zei0mcphxbdxj49qtmz
```

Example of dynamic network:

```yaml
networks:
- name: default
  type: dynamic
  cloud_properties:
    vswitch_id: vsw-2zemyfytfclbcmgfkzokx
    security_group_id: sg-2zei0mcphxbdxj49qtmz
```

Example of vip network:

```yaml
networks:
- name: default
  type: vip
```

---
## <a id='resource-pools'></a> Resource Pools / VM Types

Schema for `cloud_properties` section:

* **availability_zone** [String, required]: Availability zone to use for creating instances. Example: `cn-beijing-a`.
* **instance_type** [String, required]: Type of the [instance](https://www.alibabacloud.com/help/zh/doc-detail/25378.htm/). Example: `ecs.n1.small`.
* **security_group_id** [String, required]: See description under [networks](#networks). 
* **key_name** [String, optional]: Key pair name. Example: `bosh`.
* **slbs** [Array, optional]: Array of [Load Balancer](https://www.alibabacloud.com/help/zh/product/27537.htm) that should be attached to created VMs. 
  * **slbs** [Array, required]: SLB ID Example: `["lb-2zegrgbsmjvxx1r1v26pn"]`. 
  * **slb_weight** [Integer, optional]: SLB weight of VMs. Example `100`. Default is `[100]`.
* **system_disk** [Hash, optional]: system disk of custom size.
    * **size** [Integer, required]: Specifies the disk size in megabytes.
    * **category** [String, optional]: Category of the [disk](https://www.alibabacloud.com/help/doc-detail/25383.htm): `cloud_efficiency`, `cloud_ssd`. Defaults to `cloud_efficiency`.

* **system_disk** [Hash, optional]: EBS backed root disk of custom size.
    * **size** [Integer, required]: Specifies the disk size in megabytes. Default is `51200`
    * **category** [String, optional]: Category of the [disk](https://www.alibabacloud.com/help/doc-detail/25383.htm): `cloud_efficiency`, `cloud_ssd`. Defaults to `cloud_efficiency`.
    
```yaml
resource_pools:
- name: default
  network: default
  stemcell:
    name: light-bosh-stemcell-alicloud-kvm-ubuntu-trusty-go_agent
    version: 1008
  cloud_properties:
    instance_type: ecs.n1.small
    availability_zone: cn-beijing-a
```

---
## <a id='disk-pools'></a> Disk Pools / Disk Types

Schema for `cloud_properties` section:

* **category** [String, optional]: Category of the [disk](https://www.alibabacloud.com/help/doc-detail/25383.htm): `cloud_efficiency`, `cloud_ssd`. Defaults to `cloud_efficiency`.

Example of 10GB disk:

```yaml
disk_pools:
- name: default
  disk_size: 10_240
  cloud_properties:
    category: cloud_efficiency
```

---
## <a id='global'></a> Global Configuration

The CPI can only talk to a single Alibaba Cloud region. 

Schema:

* **region** [String, required]: Name of the `Alibaba Cloud` region, [Available regions](https://www.alibabacloud.com/help/doc-detail/40654.htm). Example: `cn-beijing`.
* **access\_key\_id** [String, optional]: Accesss Key ID. Example: `AKI...`.
* **access\_key\_secret** [String, optional]: Accesss Key Secret. Example: `0kwh...`.

Example with hard-coded credentials:

```yaml
properties:
  alicloud:
    region: cn-beijing
    access_key_id: ACCESS-KEY-ID
    secret_access_key: SECRET-ACCESS-KEY
```

---
## <a id='cloud-config'></a> Example Cloud Config

```yaml
azs:
- name: z1
  cloud_properties: {availability_zone: cn-beijing-a}
- name: z2
  cloud_properties: {availability_zone: cn-beijing-b}

vm_types:
- name: default
  cloud_properties:
    instance_type: ecs.m4.small
    ephemeral_disk: {size: 30720, category: cloud_efficiency}
- name: large
  cloud_properties:
    instance_type: ecs.m4.large
    ephemeral_disk: {size: 30720, category: cloud_efficiency}

disk_types:
- name: default
  disk_size: 10_240
  cloud_properties: {category: cloud_efficiency}
- name: large
  disk_size: 51_200
  cloud_properties: {category: cloud_efficiency}

networks:
- name: default
  type: manual
  subnets:
  - range: 10.10.0.0/24
    gateway: 10.10.0.1
    az: z1
    static: [10.10.0.62]
    dns: [10.10.0.2]
    cloud_properties: {vswitch_id: "vsw-9q8asd9q1243sad234234", security_group_id: "sg-a98234oiwoierupoi"}
  - range: 10.10.64.0/24
    gateway: 10.10.64.1
    az: z2
    static: [10.10.64.121, 10.10.64.122]
    dns: [10.10.0.2]
    cloud_properties: {vswitch_id: "vsw-99823hpoiaoipoiu34234", security_group_id: "sg-9uoiuere993499f"}
- name: vip
  type: vip

compilation:
  workers: 5
  reuse_compilation_vms: true
  az: z1
  vm_type: large
  network: default
```

---
## <a id='errors'></a> Errors

**TODO**