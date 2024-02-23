packer {
  required_version = ">= 1.7.0"
  required_plugins {
    googlecompute = {
      version = ">= 1.0"
      source  = "github.com/hashicorp/googlecompute"
    }
  }
}

locals {
  timestamp = regex_replace(timestamp(), "[- TZ:]", "")
}

variable "gcp_project_id" {
  type    = string
  default = "csye-6225-terraform-packer"
}

variable "source_image_family" {
  type    = string
  default = "centos-stream-8"
}

variable "machine_type" {
  type    = string
  default = "custom-1-2048"
}

variable "application_name" {
  type    = string
  default = "webapp"
}

variable "service_name" {
  type    = string
  default = "webapp.service"
}

variable "zone" {
  type    = string
  default = "us-central1-a"
}

variable "ssh_username" {
  type    = string
  default = "centos"
}

variable golang_version {
  type    = string
  default = ""
}

source "googlecompute" "webapp-source" {
  image_name          = "webapp-${local.timestamp}"
  project_id          = var.gcp_project_id
  machine_type        = var.machine_type
  source_image_family = var.source_image_family
  ssh_username        = var.ssh_username
  zone                = var.zone
}

build {
  sources = [
    "source.googlecompute.webapp-source"
  ]

  // PostgreSQL Installation
  provisioner "shell" {
    inline = [
      "sudo yum update -y",
      "sudo yum install -y postgresql-server postgresql-contrib",
      "sudo postgresql-setup --initdb --unit postgresql",
      "sudo systemctl enable postgresql",
      "sudo systemctl start postgresql",
    ]
  }

  // PostgreSQL user and database creation and assign perms
  provisioner "shell" {
    script = "./db.sh"
  }

  provisioner "file" {
    source      = "./webapp"
    destination = "/tmp/webapp"
  }

  provisioner "file" {
    source      = "./webapp.service"
    destination = "/tmp/webapp.service"
  }

  // provisioner "file" {
  //   source      = "./.env"
  //   destination = "/tmp/.env"
  // }

  provisioner "shell" {
    inline = [
      // Create group and user
      "sudo groupadd csye6225",
      "sudo useradd -g csye6225 -m csye6225",

      // Move webapp and enable service
      "sudo mv /tmp/webapp /usr/local/bin",
      // "sudo mv /tmp/.env /usr/local/bin",
      "sudo mv /tmp/webapp.service /etc/systemd/system",

      "sudo sed -i 's/^SELINUX=.*/SELINUX=disabled/' /etc/selinux/config",
      "sudo restorecon -rv /usr/local/bin/webapp",

      "sudo chown csye6225:csye6225 /usr/local/bin/webapp",
      // "sudo chown csye6225:csye6225 /usr/local/bin/.env",
      "sudo chmod 755 /usr/local/bin/webapp",
      // "sudo chmod 755 /usr/local/bin/.env",

      //set nologin to webapp user
      "sudo usermod csye6225 --shell /usr/sbin/nologin",
    ]
  }

  // Enable and start webapp
  provisioner "shell" {
    script = "./webapp_start.sh"
  }
}