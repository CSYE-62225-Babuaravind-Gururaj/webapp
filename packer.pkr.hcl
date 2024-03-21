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
  default = "asia-east1-a"
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
  // provisioner "shell" {
  //   inline = [
  //     "sudo yum update -y",
  //     "sudo yum install -y postgresql-server postgresql-contrib",
  //     "sudo postgresql-setup --initdb --unit postgresql",
  //     "sudo systemctl enable postgresql",
  //     "sudo systemctl start postgresql",
  //   ]
  // }

  // PostgreSQL user and database creation and assign perms
  // provisioner "shell" {
  //   script = "./db.sh"
  // }

  provisioner "file" {
    source      = "./webapp"
    destination = "/tmp/webapp"
  }

  provisioner "file" {
    source      = "./webapp.zip"
    destination = "/tmp/webapp.zip"
  }

  provisioner "file" {
    source      = "./webapp.service"
    destination = "/tmp/webapp.service"
  }

  // provisioner "file" {
  //   source      = "./.env"
  //   destination = "/tmp/.env"
  // }

  provisioner "file" {
  source      = "./config.yaml"  
  destination = "/tmp/config.yaml"     
}

  provisioner "shell" {
    inline = [
      // Create group and user
      "getent group csye6225 || sudo groupadd csye6225",
      // Check if user exists, create if it does not
      "id -u csye6225 || sudo useradd -g csye6225 -m csye6225",

      // Move webapp and enable service
      "sudo mv /tmp/webapp /usr/local/bin",
      "sudo mv /tmp/webapp.zip /usr/local/bin",
      // "sudo mv /tmp/.env /usr/local/bin",
      "sudo mv /tmp/webapp.service /etc/systemd/system",

      "sudo sed -i 's/^SELINUX=.*/SELINUX=disabled/' /etc/selinux/config",
      "sudo restorecon -rv /usr/local/bin/webapp",

      // "sudo touch /home/csye6225/webapp/userdata.sh",
      "sudo chown csye6225:csye6225 /usr/local/bin/webapp",
      // "sudo chown csye6225:csye6225 /usr/local/bin/.env",
      "sudo chmod 750 /usr/local/bin/webapp",
      // "sudo chmod 755 /usr/local/bin/.env",

      "sudo mkdir -p /var/log/myapp",
      "sudo touch /var/log/myapp/app.log",
      "sudo chown csye6225:csye6225 /var/log/myapp/app.log",
      "sudo chmod 766 /var/log/myapp/app.log",

      //set nologin to webapp user
      "sudo usermod csye6225 --shell /usr/sbin/nologin",
      "cd"
    ]
  }

  // Cloud Ops Agent installation
  provisioner "shell" {
    inline = [
      "curl -sSO https://dl.google.com/cloudagents/add-google-cloud-ops-agent-repo.sh",
      "sudo bash add-google-cloud-ops-agent-repo.sh --also-install"
    ]
  }

  provisioner "shell" {
    inline = [
      "sudo mv /tmp/config.yaml /etc/google-cloud-ops-agent/config.yaml",
      "sudo systemctl enable google-cloud-ops-agent",
      "sudo systemctl restart google-cloud-ops-agent.service"
    ]
  }

  // Enable and start webapp
  provisioner "shell" {
    script = "./webapp_start.sh"
  }
}