variable "region" {
}

variable "env_name" {
}

variable "public_key" {
}

variable "concourse_worker_ip" {
  default = ""
}

# Name of the pre-created RAM role attached to the Director VM. It is created
# once, out of band, because the pipeline's provisioning role is intentionally
# not granted ram:CreatePolicy / ram:AttachPolicyToRole (that would let the
# pipeline grant itself arbitrary permissions). See ci/README.md.
variable "director_role_name" {
  default = "BoshDirectorRole"
}

terraform {
  backend "oss" {
  }
  required_providers {
    alicloud = {
      source  = "aliyun/alicloud"
      version = "1.238.0"
    }
  }
}

# Credentials come from the environment (ALICLOUD_ACCESS_KEY, ALICLOUD_SECRET_KEY
# and ALICLOUD_SECURITY_TOKEN), which the pipeline fills from the role it
# assumed. Declaring them as variables would put them into the plan and into any
# saved plan file.
provider "alicloud" {
  region = var.region
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
  cidr_ip           = alicloud_vpc.default.cidr_block
}

# Allow Concourse worker to reach BOSH Director ports (mbus, API, SSH, etc.)
resource "alicloud_security_group_rule" "concourse-in" {
  count             = var.concourse_worker_ip != "" ? 1 : 0
  type              = "ingress"
  ip_protocol       = "all"
  nic_type          = "intranet"
  policy            = "accept"
  port_range        = "-1/-1"
  priority          = 1
  security_group_id = alicloud_security_group.default.id
  cidr_ip           = "${var.concourse_worker_ip}/32"
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

resource "random_integer" "default" {
  max = 99999
  min = 10000
}

resource "alicloud_oss_bucket" "blobstore" {
  bucket = "cpi-pipeline-blobstore-${var.env_name}-${random_integer.default.result}"
  acl    = "private"
}

resource "alicloud_key_pair" "director" {
  key_pair_name = var.env_name
  public_key    = var.public_key
}

locals {
  # No fallback to env_name: the role is pre-created and shared, not per-env.
  director_role_name = var.director_role_name
}

# The Director VM runs the CPI under this pre-created RAM role and discovers its
# credentials from instance metadata, so no access key is ever written to the
# Director. The role and its permission policy are created once, out of band
# (see ci/README.md); the pipeline only references it by name and attaches it to
# the VM at create time.
#
# Terraform deliberately does not look up, create, or manage the role: the
# provisioning role is intentionally granted neither ram:ListRoles nor
# ram:CreatePolicy / ram:AttachPolicyToRole, so it cannot grant itself
# permissions. If the role is missing, `bosh create-env` fails with a clear
# "RAM role not found" error when it attaches the role to the Director VM.

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
  value = local.director_role_name
}

