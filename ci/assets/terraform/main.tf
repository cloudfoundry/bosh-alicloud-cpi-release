provider "alicloud" {
  access_key = "${var.alicloud_access_key}"
  secret_key = "${var.alicloud_secret_key}"
  region = "${var.alicloud_region}"
}

data "alicloud_instance_types" "1c2g" {
  cpu_core_count = 2
  memory_size = 4
}

//// Images data source for image_id
//data "alicloud_images" "default" {
//  most_recent = true
//  owners = "system"
//  name_regex = "${var.image_name_regex}"
//}
//
data "alicloud_zones" "default" {
  "available_instance_type"= "${data.alicloud_instance_types.1c2g.instance_types.1.id}"
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

//resource "alicloud_nat_gateway" "default" {
//  vpc_id = "${alicloud_vpc.default.id}"
//  spec = "Small"
//  name = "bosh_init_nat${var.prefix}"
//  bandwidth_packages = [{
//    ip_count = 2
//    bandwidth = 100
//    zone = "${data.alicloud_zones.default.zones.0.id}"
//  }]
//  depends_on = [
//    "alicloud_vswitch.default"]
//}
////resource "alicloud_snat_entry" "default"{
////  snat_table_id = "${alicloud_nat_gateway.default.snat_table_ids}"
////  source_vswitch_id = "${alicloud_vswitch.default.id}"
////  snat_ip = "${element(split(",", alicloud_nat_gateway.default.bandwidth_packages.0.public_ip_addresses),0)}"
////}
//
//resource "alicloud_forward_entry" "default"{
//  forward_table_id = "${alicloud_nat_gateway.default.forward_table_ids}"
//  external_ip = "${element(split(",", alicloud_nat_gateway.default.bandwidth_packages.0.public_ip_addresses),1)}"
//  external_port = "80"
//  ip_protocol = "tcp"
//  internal_ip = "${var.router_private_ip}"
////  internal_ip = "${alicloud_instance.default.private_ip}"
//  internal_port = "80"
//}
//
//resource "alicloud_forward_entry" "443"{
//  forward_table_id = "${alicloud_nat_gateway.default.forward_table_ids}"
//  external_ip = "${element(split(",", alicloud_nat_gateway.default.bandwidth_packages.0.public_ip_addresses),1)}"
//  external_port = "443"
//  ip_protocol = "tcp"
//  internal_ip = "${var.uaa_private_ip}"
//  internal_port = "8443"
//}
//
////resource "alicloud_forward_entry" "ssh"{
////  forward_table_id = "${alicloud_nat_gateway.default.forward_table_ids}"
////  external_ip = "${element(split(",", alicloud_nat_gateway.default.bandwidth_packages.0.public_ip_addresses),1)}"
////  external_port = "22"
////  ip_protocol = "tcp"
////  internal_ip = "${alicloud_instance.default.private_ip}"
////  internal_port = "22"
////}
//
//resource "alicloud_forward_entry" "bosh_target"{
//  forward_table_id = "${alicloud_nat_gateway.default.forward_table_ids}"
//  external_ip = "${element(split(",", alicloud_nat_gateway.default.bandwidth_packages.0.public_ip_addresses),1)}"
//  external_port = "25555"
//  ip_protocol = "tcp"
//  internal_ip = "${var.bosh_ip}"
//  internal_port = "25555"
//}
//
//resource "alicloud_forward_entry" "bosh_target_6868"{
//  forward_table_id = "${alicloud_nat_gateway.default.forward_table_ids}"
//  external_ip = "${element(split(",", alicloud_nat_gateway.default.bandwidth_packages.0.public_ip_addresses),1)}"
//  external_port = "6868"
//  ip_protocol = "tcp"
//  internal_ip = "${var.bosh_ip}"
//  internal_port = "6868"
//}
//
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
//resource "alicloud_security_group_rule" "http-in" {
//  type = "ingress"
//  ip_protocol = "tcp"
//  nic_type = "intranet"
//  policy = "accept"
//  port_range = "80/80"
//  priority = 1
//  security_group_id = "${alicloud_security_group.default.id}"
//  cidr_ip = "0.0.0.0/0"
//}
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
  bandwidth = "10"
  internet_charge_type = "PayByBandwidth"
}

//resource "alicloud_key_pair" "key_pair" {
//  key_name = "key_pair_for_cloudfoundary"
//  key_file = "../deploy/roles/bosh-deploy/files/bosh-init.pem"
//}
//
//resource "alicloud_instance" "default" {
////  availability_zone = "${data.alicloud_zones.default.zones.0.id}"
//  security_groups = [
//    "${alicloud_security_group.default.id}"]
//
//  vswitch_id = "${alicloud_vswitch.default.id}"
//
//  password = "${var.ecs_password}"
//
//  # series III
//  instance_charge_type = "PostPaid"
//  instance_type = "${data.alicloud_instance_types.1c2g.instance_types.1.id}"
//  internet_max_bandwidth_out = 0
//
//  system_disk_category = "cloud_efficiency"
//  system_disk_size = 100
//  image_id = "${data.alicloud_images.default.images.0.id}"
//  instance_name = "bosh_init${var.prefix}"
//  //  allocate_public_ip = true
//
//
//  //  key_name = "${alicloud_key_pair.key_pair.id}"
//
//  provisioner "local-exec" {
//    command = <<EOF
//        echo [CloudFoundaryServer] > ../hosts
//        echo ${element(split(",", alicloud_nat_gateway.default.bandwidth_packages.0.public_ip_addresses),1)} ansible_user=root ansible_ssh_pass=${var.ecs_password} >> ../hosts
//        echo \n
//        echo internal_cidr: ${var.vpc_cidr} >> ../group_vars/all
//        echo internal_gw: 172.16.0.1 >> ../group_vars/all
//        echo internal_ip: ${var.bosh_ip} >> ../group_vars/all
//        echo security_group_id: ${alicloud_security_group.default.id} >> ../group_vars/all
//        echo subnet_id: ${alicloud_vswitch.default.id} >> ../group_vars/all
//        echo xip_ip_domain: ${element(split(",", alicloud_nat_gateway.default.bandwidth_packages.0.public_ip_addresses),1)}.xip.io >> ../group_vars/all
//
//  EOF
//  }
//}
