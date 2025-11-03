# Generate an SSH key locally and upload the public key as an AWS key pair
resource "tls_private_key" "endercloud" {
  algorithm = "RSA"
  rsa_bits  = 4096
}

resource "aws_key_pair" "endercloud" {
  key_name   = var.ssh_key_name
  public_key = tls_private_key.endercloud.public_key_openssh
}

# Store private key in a secure local folder (outside repo)
locals {
  private_key_path = "${pathexpand("~/.ssh/EnderCloud/endercloud_key.pem")}"
}

resource "local_file" "private_key" {
  content         = tls_private_key.endercloud.private_key_pem
  filename        = local.private_key_path
  file_permission = "0600"
}
