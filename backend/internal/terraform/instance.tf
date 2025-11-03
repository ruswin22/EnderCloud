resource "aws_instance" "mc" {
  ami                    = data.aws_ami.ubuntu.id
  instance_type          = var.instance_type
  key_name               = aws_key_pair.endercloud.key_name
  vpc_security_group_ids = [aws_security_group.mc_sg.id]

  tags = {
    Name = var.minecraft_tag
  }

  # optional: add a small user_data to update packages quickly (Ansible will do rest)
  user_data = <<-EOF
              #!/bin/bash
              apt-get update -y
              EOF
}

output "public_ip" {
  value = aws_instance.mc.public_ip
}

output "private_key_path" {
  value = local_file.private_key.filename
}
