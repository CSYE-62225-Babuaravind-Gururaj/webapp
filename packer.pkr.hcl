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
  default = "e2-medium"
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
    type = string
    default =""
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

  provisioner "file" {
    source      = "./webapp.zip"
    destination = "/tmp/webapp.zip"
  }

  provisioner "file" {
    source      = "./application.service"
    destination = "/tmp/application.service"
  }

  provisioner "shell" {
    inline = [
    "sudo yum update -y",
    "sudo groupadd csye6225",
    "sudo useradd -g csye6225 -m csye6225",
    // "sudo id -u csye6225 &>/dev/null || sudo useradd -g csye6225 -m csye6225",
    "sudo yum install unzip -y",
    "sudo mv /tmp/application.service /etc/systemd/system/${var.service_name}",
    "cd /tmp",
    "sudo unzip webapp.zip -d /home/csye6225",
    "cd",

    //Golang
    "sudo yum install wget -y",
    "wget https://golang.org/dl/go1.21.6.linux-amd64.tar.gz",
    "sudo tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz",
    "echo 'export PATH=$PATH:/usr/local/go/bin' >> $HOME/.bash_profile",
    "echo 'export GOPATH=$HOME/go' >> $HOME/.bash_profile",
    "echo 'export PATH=$PATH:$GOPATH/bin' >> $HOME/.bash_profile",
    "source $HOME/.bash_profile",

    //Postgres
    "sudo yum install -y postgresql-server postgresql-contrib",
    "sudo postgresql-setup --initdb",
    "sudo systemctl enable postgresql",
    "sudo systemctl start postgresql",
    "sudo -u postgres psql -c \"CREATE USER csye6225 WITH ENCRYPTED PASSWORD 'root';\"",
    "sudo -u postgres psql -c \"GRANT ALL PRIVILEGES ON DATABASE postgres TO csye6225;\"",

    "sudo systemctl enable postgresql",
    "sudo usermod --shell /usr/sbin/nologin csye6225",
    "sudo chown -R csye6225:csye6225 /home/csye6225",
    "sudo chmod -R 755 /home/csye6225",
    "sudo chown csye6225:csye6225 /etc/systemd/system/${var.service_name}",
    "sudo chmod 644 /etc/systemd/system/${var.service_name}"
    ]
  }

  provisioner "shell" {
    script = "./webapp_start.sh"
  }
}