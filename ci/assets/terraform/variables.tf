variable "alicloud_access_key" {
}

variable "alicloud_secret_key" {
}

variable "alicloud_region" {
}

variable "vpc_cidr" {
  default = "172.16.0.0/24"
}
variable "vswitch_cidr" {
  default = "172.16.0.0/24"
}
variable "router_private_ip" {
  default = "172.16.0.27"
}
variable "uaa_private_ip" {
  default = "172.16.0.25"
}
variable "bosh_ip" {
  default = "172.16.0.2"
}
variable "rule_policy" {
  default = "accept"
}
variable "instance_type" {
  default = "ecs.n4.small"
}
# Image variables
variable "image_name_regex" {
  description = "The ECS image's name regex used to fetch specified image."
  default = "^ubuntu_16.*_64"
}
variable "disk_category"{
  default = "cloud_efficiency"
}
variable "ecs_password"{
  default = "Cloud12345"
}
variable "prefix"{
  default = "_911"
}