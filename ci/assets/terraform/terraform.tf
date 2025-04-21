variable "access_key" {
}

variable "secret_key" {
}

variable "security_token" {
}

variable "region" {
}

variable "env_name" {
}

variable "public_key" {
}

terraform {
  backend "oss" {
  }
  required_providers {
    alicloud = {
      source  = "aliyun/alicloud"
      version = "1.193.0"
    }
  }
}

provider "alicloud" {
  access_key     = var.access_key
  secret_key     = var.secret_key
  security_token = var.security_token
  region         = var.region
}

data "alicloud_zones" "default" {
}

# Create a VPC to launch our instances into
resource "alicloud_vpc" "default" {
  vpc_name   = var.env_name
  cidr_block = "172.16.0.0/16"
}

# Create an nat gateway to give our vswitch access to the outside world
resource "alicloud_nat_gateway" "default" {
  vpc_id     = alicloud_vpc.default.id
  name       = var.env_name
  vswitch_id = alicloud_vswitch.default.id
  nat_type   = "Enhanced"
}

resource "alicloud_eip" "default" {
  internet_charge_type = "PayByTraffic"
  name                 = var.env_name
}

resource "alicloud_eip_association" "default" {
  instance_id   = alicloud_nat_gateway.default.id
  allocation_id = alicloud_eip.default.id
}

resource "alicloud_snat_entry" "a" {
  snat_table_id     = alicloud_nat_gateway.default.snat_table_ids
  source_vswitch_id = alicloud_vswitch.default.id
  snat_ip           = alicloud_eip.default.ip_address
  depends_on        = [alicloud_eip_association.default]
}

resource "alicloud_snat_entry" "b" {
  snat_table_id     = alicloud_nat_gateway.default.snat_table_ids
  source_vswitch_id = alicloud_vswitch.backup.id
  snat_ip           = alicloud_eip.default.ip_address
  depends_on        = [alicloud_eip_association.default]
}

resource "alicloud_snat_entry" "c" {
  snat_table_id     = alicloud_nat_gateway.default.snat_table_ids
  source_vswitch_id = alicloud_vswitch.manual.id
  snat_ip           = alicloud_eip.default.ip_address
  depends_on        = [alicloud_eip_association.default]
}

resource "alicloud_vswitch" "default" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 0)
  zone_id      = data.alicloud_zones.default.zones[0].id
  vswitch_name = var.env_name
}

resource "alicloud_vswitch" "backup" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 2)
  zone_id      = data.alicloud_zones.default.zones[1].id
  vswitch_name = var.env_name
}

resource "alicloud_vswitch" "manual" {
  vpc_id       = alicloud_vpc.default.id
  cidr_block   = cidrsubnet(alicloud_vpc.default.cidr_block, 8, 4)
  zone_id      = data.alicloud_zones.default.zones[0].id
  vswitch_name = var.env_name
}

resource "alicloud_security_group" "default" {
  name        = var.env_name
  description = "Allow all inbound and outgoing traffic"
  vpc_id      = alicloud_vpc.default.id
}

resource "alicloud_security_group_rule" "all-in" {
  type              = "ingress"
  ip_protocol       = "all"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "-1/-1"
  priority          = 1
  security_group_id = alicloud_security_group.default.id
  cidr_ip           = "0.0.0.0/0"
}

resource "alicloud_security_group_rule" "all-out" {
  type              = "egress"
  ip_protocol       = "all"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "-1/-1"
  priority          = 1
  security_group_id = alicloud_security_group.default.id
  cidr_ip           = "0.0.0.0/0"
}

resource "alicloud_eip" "director" {
  internet_charge_type = "PayByTraffic"
  name                 = var.env_name
}

resource "alicloud_eip" "deployment" {
  internet_charge_type = "PayByTraffic"
  name                 = var.env_name
}

# Create a new classic load balancer
resource "alicloud_slb" "default" {
  name                 = var.env_name
  internet_charge_type = "PayByTraffic"
  address_type         = "internet"
  specification        = "slb.s1.small"
}


resource "alicloud_slb_server_group" "default" {
  load_balancer_id = alicloud_slb.default.id
  name             = var.env_name
}
resource "alicloud_slb_listener" "http" {
  load_balancer_id = alicloud_slb.default.id
  backend_port     = 80
  frontend_port    = 80
  protocol         = "http"
  bandwidth        = 10
  health_check     = "off"
}

# Create a new application load balancer
resource "alicloud_slb" "app" {
  name                 = var.env_name
  vswitch_id           = alicloud_vswitch.default.id
  internet_charge_type = "PayByTraffic"
  specification        = "slb.s1.small"
}

resource "alicloud_slb_listener" "app-http" {
  load_balancer_id          = alicloud_slb.app.id
  backend_port              = 80
  frontend_port             = 80
  protocol                  = "http"
  bandwidth                 = 10
  health_check              = "on"
  health_check_timeout      = 4
  health_check_interval     = 5
  health_check_http_code    = "http_2xx"
  health_check_connect_port = 20
}

resource "alicloud_oss_bucket" "blobstore" {
  bucket = "cpi-pipeline-blobstore-${var.env_name}"
  acl    = "private"
}

resource "alicloud_key_pair" "director" {
  key_pair_name = var.env_name
  public_key    = var.public_key
}

resource "alicloud_ram_role" "role" {
  name        = var.env_name
  description = "a role for bosh integration test"
  force       = true
  document    = <<EOF
  {
    "Statement": [
      {
        "Action": "sts:AssumeRole",
        "Effect": "Allow",
        "Principal": {
          "Service": [
            "ecs.aliyuncs.com"
          ]
        }
      }
    ],
    "Version": "1"
  }
  
EOF

}

output "vpc_id" {
  value = alicloud_vpc.default.id
}

output "region" {
  value = var.region
}

# Used by bats
output "key_pair_name" {
  value = alicloud_key_pair.director.key_pair_name
}

output "security_group_id" {
  value = alicloud_security_group.default.id
}

output "external_ip" {
  value = alicloud_eip.director.ip_address
}

output "zone" {
  value = alicloud_vswitch.default.zone_id
}

output "vswitch_id" {
  value = alicloud_vswitch.default.id
}

output "manual_vswitch_id" {
  value = alicloud_vswitch.manual.id
}

output "internal_cidr" {
  value = alicloud_vpc.default.cidr_block
}

output "internal_gw" {
  value = cidrhost(alicloud_vpc.default.cidr_block, 1)
}

output "dns_recursor_ip" {
  value = "8.8.8.8"
}

output "internal_ip" {
  value = cidrhost(alicloud_vpc.default.cidr_block, 6)
}

output "reserved_range" {
  value = "${cidrhost(alicloud_vpc.default.cidr_block, 2)}-${cidrhost(alicloud_vpc.default.cidr_block, 9)}"
}

output "static_range" {
  value = "${cidrhost(alicloud_vpc.default.cidr_block, 10)}-${cidrhost(alicloud_vpc.default.cidr_block, 30)}"
}

output "bats_eip" {
  value = alicloud_eip.deployment.ip_address
}

output "network_static_ip_1" {
  value = cidrhost(alicloud_vpc.default.cidr_block, 29)
}

output "network_static_ip_2" {
  value = cidrhost(alicloud_vpc.default.cidr_block, 30)
}

output "slb" {
  value = alicloud_slb.default.id
}

output "server_group_slb" {
  value = alicloud_slb_server_group.default.id
}

output "blobstore_bucket" {
  value = alicloud_oss_bucket.blobstore.id
}

output "integration_bucket" {
  value = alicloud_oss_bucket.blobstore.id
}

output "ram_role" {
  value = alicloud_ram_role.role.name
}

