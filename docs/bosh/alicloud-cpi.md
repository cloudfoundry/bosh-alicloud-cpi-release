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
    ip: 10.0.0.3
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
  ip: 47.47.47.47
```

---
## <a id='resource-pools'></a> Resource Pools / VM Types

Schema for `cloud_properties` section:

* **availability_zone** [String, required]: Availability zone to use for creating instances. Example: `cn-beijing-a`.
* **instance\_type** [String, required]: Type of the [instance](https://www.alibabacloud.com/help/zh/doc-detail/25378.htm/). Example: `ecs.n1.small`.
* **instance\_name** [String, required]: Instance host name.
* **charge\_type** [String, optional]: Charge type of instance: `PrePaid` or `PostPaid`. Default is `PostPaid`
* **charge\_period** [Integer, optional]: Prepaid months (range: 1-9, 12, 24, 36), required if `charge_type` is `PrePaid`. 
* **spot\_strategy** [String, optional]: The spot strategy of a Pay-As-You-Go instance, and it takes effect only when the `charge_type` is `PostPaid`. Value range:
    * **NoSpot** : A normal Pay-As-You-Go instance.
    * **SpotWithPriceLimit** : A price threshold for a spot instance.
    * **SpotAsPriceGo** : A price based on the highest Pay-As-You-Go instance will be automatically generated.
* **spot\_price\_limit** [Float, optional]: The hourly price threshold of a instance, and it takes effect only when the `spot_strategy` is `SpotWithPriceLimit`. Three decimals is allowed at most.
* **auto\_renew** [Boolean, optional]: `True` or `False`, when charge type is `Prepaid`, will auto renew your payment. Default is `False`.
* **auto\_renew\_period** [Integer, optional]: Required if `auto_renew` is `True`, by months (range: 1, 2, 3, 6, 12). 
* **password** [String, optional]: Root password, no-effect with some stemcell, use jumpbox-user instead.
* **key\_pair\_name** [String, optional]: Key pair name, no-effect with some stemcell, use jumpbox-user instead. Example: `bosh`.
* **slbs** [Array, optional]: Array of [Load Balancer](https://www.alibabacloud.com/help/zh/product/27537.htm) that should be attached to created VMs. Example: `["lb-2zegrgbsmjvxx1r1v26pn"]`.
* **slb_weight** [Integer, optional]: SLB weight of VMs. Example `100`. Default is `[100]`.
* **system_disk** [Hash, optional]: system disk of custom size.
    * **size** [Integer, required]: Specifies the disk size in megabytes.
    * **category** [String, optional]: Category of the [disk](https://www.alibabacloud.com/help/doc-detail/25383.htm): `cloud_efficiency`, `cloud_ssd`. Defaults to `cloud_efficiency`.
* **ephemeral_disk** [Hash, optional]: ephemeral disk of custom size.
    * **size** [Integer, required]: Specifies the disk size in megabytes. Default is `51200`
    * **category** [String, optional]: Category of the [disk](https://www.alibabacloud.com/
    * **delete\_with\_instance** [Boolean, optional]: Will deleted with instance. Default is `True`
help/doc-detail/25383.htm): `cloud_efficiency`, `cloud_ssd`. Defaults to `cloud_efficiency`.
    
```yaml
resource_pools:
- name: default
  network: default
  stemcell:
    name: bosh-stemcell-alicloud-kvm-ubuntu-trusty-go_agent
    version: 1009
  cloud_properties:
    availability_zone: cn-beijing-a
    instance_type: ecs.n1.small
    instance_charge_type: PostPaid
    slbs: ["lb-2zegrgbsmjvxx1r1v26pn"]
    system_disk: {"size": "61_440", "category": "cloud_efficiency"}

```

---
## <a id='disk-pools'></a> Disk Pools / Disk Types

Schema for `cloud_properties` section:

* **category** [String, optional]: Category of the [disk](https://www.alibabacloud.com/help/doc-detail/25383.htm): `cloud_efficiency`, `cloud_ssd`. Defaults to `cloud_efficiency`.
* **encrypted** [Boolean, optional]: Enables encryption of your data disk. Default is `False`

Example of 20GB disk:

```yaml
disk_pools:
- name: default
  disk_size: 20_480
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
  cloud_properties:
    availability_zone: cn-beijing-a
- name: z2
  cloud_properties:
    availability_zone: cn-beijing-d
- name: z3
  cloud_properties:
    availability_zone: cn-beijing-e

vm_types:
- name: minimal
  cloud_properties:
    instance_type: ecs.mn4.small
    ephemeral_disk: {size: "51_200"}
- name: small
  cloud_properties:
    instance_type: ecs.sn2.medium
    ephemeral_disk: {size: "51_200"}
- name: default
  cloud_properties:
    instance_type: ecs.sn2.medium
    ephemeral_disk: {size: "51_200"}	
- name: small-highmem
  cloud_properties:
    instance_type: ecs.sn2ne.xlarge
    ephemeral_disk: {size: "51_200"}
- name: compiler
  cloud_properties:
    instance_type: ecs.sn1.large
    ephemeral_disk: {size: "51_200"}

disk_types:
- name: 5GB
  disk_size: 20_480
- name: 10GB
  disk_size: 20_480
- name: 100GB
  disk_size: 102_400

vm_extensions:
- name: 5GB_ephemeral_disk
  cloud_properties:
    ephemeral_disk: {size: "20_480"}
- name: 10GB_ephemeral_disk
  cloud_properties:
    ephemeral_disk: {size: "20_480"}
- name: 50GB_ephemeral_disk
  cloud_properties:
    ephemeral_disk: {size: "50_120"}
- name: 100GB_ephemeral_disk
  cloud_properties:
    ephemeral_disk: {size: "102_400"}
- name: 500GB_ephemeral_disk
  cloud_properties:
    ephemeral_disk: {size: "512_000"}
- name: 1TB_ephemeral_disk
  cloud_properties:
    ephemeral_disk: {size: "1024_000"}
- name: cf-router-network-properties
  cloud_properties: 
    slbs: ["lb-2zegrgbsmjvxx1r1v26pn"]
- name: cf-tcp-router-network-properties
- name: diego-ssh-proxy-network-properties
  
networks:
- name: default
  type: manual
  subnets:
  - range: 192.168.10.0/24
    gateway: 192.168.10.1
    az: z1
    dns: [8.8.8.8]
    cloud_properties:
      vswitch_id: vsw-2zeamad3a8cscoicqb5c5
      security_group_id: sg-2zei0mcphxbdxj49qtmz
  - range: 192.168.16.0/24
    gateway: 192.168.16.1
    az: z2
    dns: [8.8.8.8]
    cloud_properties:
      vswitch_id: vsw-2zerkt1jluc2xdxygeu5t
      security_group_id: sg-2zei0mcphxbdxj49qtmz
  - range: 192.168.11.0/24
    gateway: 192.168.11.1
    az: z3
    dns: [8.8.8.8]
    cloud_properties:
      vswitch_id: vsw-2zedja4ggcyrahgz0s7cc
      security_group_id: sg-2zei0mcphxbdxj49qtmz 
- name: vip
  type: vip

compilation:
  workers: 5
  reuse_compilation_vms: true
  az: z1
  vm_type: compiler
  network: default
```

---
## <a id='errors'></a> Errors

**TODO**