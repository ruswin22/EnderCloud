variable "aws_region" {
  type    = string
  default = "ap-south-1"
}

variable "instance_type" {
  type    = string
  default = "t2.micro"
}

variable "minecraft_tag" {
  type    = string
  default = "EnderCloud-MC"
}

variable "ssh_key_name" {
  type    = string
  default = "endercloud_key"
}
