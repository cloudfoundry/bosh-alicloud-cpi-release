provider "alicloud" {
  access_key = "${var.alicloud_access_key}"
  secret_key = "${var.alicloud_secret_key}"
  region = "${var.alicloud_region}"
}

data "alicloud_instance_types" "1c2g" {
  cpu_core_count = 2
  memory_size = 4
}

data "alicloud_zones" "default" {
  "available_instance_type"= "${data.alicloud_instance_types.1c2g.instance_types.0.id}"
  "available_disk_category"= "${var.disk_category}"
}

resource "alicloud_vpc" "default" {
  name = "concourse_${var.prefix}"
  cidr_block = "${var.vpc_cidr}"
}

resource "alicloud_vswitch" "default" {
  vpc_id = "${alicloud_vpc.default.id}"
  cidr_block = "${var.vswitch_cidr}"
  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
}

resource "alicloud_security_group" "default" {
  name = "bosh_init_sg${var.prefix}"
  description = "tf_sg"
  vpc_id = "${alicloud_vpc.default.id}"
}

resource "alicloud_security_group_rule" "all-in" {
  type = "ingress"
  ip_protocol = "all"
  nic_type = "intranet"
  policy = "accept"
  port_range = "-1/-1"
  priority = 1
  security_group_id = "${alicloud_security_group.default.id}"
  cidr_ip = "0.0.0.0/0"
}

resource "alicloud_security_group_rule" "http-out" {
  type = "egress"
  ip_protocol = "all"
  nic_type = "intranet"
  policy = "accept"
  port_range = "-1/-1"
  priority = 1
  security_group_id = "${alicloud_security_group.default.id}"
  cidr_ip = "0.0.0.0/0"
}


resource "alicloud_security_group_rule" "ssh" {
  type = "ingress"
  ip_protocol = "tcp"
  nic_type = "intranet"
  policy = "accept"
  port_range = "22/22"
  priority = 1
  security_group_id = "${alicloud_security_group.default.id}"
  cidr_ip = "0.0.0.0/0"
}

resource "alicloud_security_group_rule" "boshtarget" {
  type = "ingress"
  ip_protocol = "tcp"
  nic_type = "intranet"
  policy = "accept"
  port_range = "25555/25555"
  priority = 1
  security_group_id = "${alicloud_security_group.default.id}"
  cidr_ip = "0.0.0.0/0"
}

resource "alicloud_security_group_rule" "boshagent" {
  type = "ingress"
  ip_protocol = "tcp"
  nic_type = "intranet"
  policy = "accept"
  port_range = "6868/6868"
  priority = 1
  security_group_id = "${alicloud_security_group.default.id}"
  cidr_ip = "0.0.0.0/0"
}

resource "alicloud_eip" "default" {
  count = 2
  bandwidth = "10"
  internet_charge_type = "PayByBandwidth"
}

resource "alicloud_slb" "http" {
  name = "for_concourse${var.prefix}"
  vswitch_id = "${alicloud_vswitch.default.id}"
  internet_charge_type = "paybytraffic"
  listener = [
    {
      "instance_port" = "80"
      "lb_port" = "80"
      "lb_protocol" = "http"
      "bandwidth" = "5"
    },
    {
      "instance_port" = "8443"
      "lb_port" = "443"
      "lb_protocol" = "tcp"
      "bandwidth" = "5"
    }]
}

resource "alicloud_slb" "tcp" {
  name = "for_concourse${var.prefix}"
  vswitch_id = "${alicloud_vswitch.default.id}"
  internet_charge_type = "paybytraffic"
  listener = [
    {
      "instance_port" = "80"
      "lb_port" = "80"
      "lb_protocol" = "tcp"
      "bandwidth" = "5"
    },
    {
      "instance_port" = "443"
      "lb_port" = "443"
      "lb_protocol" = "tcp"
      "bandwidth" = "5"
    }]
}

