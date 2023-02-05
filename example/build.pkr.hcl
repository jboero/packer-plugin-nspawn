/*
packer {
  required_plugins {
    nspawn = {
      version = ">=v0.1.0"
      source  = "github.com/hashicorp/nspawn"
    }
  }
}
*/

source "nspawn-machine" "basic-example" {
  image = "rawhide"
}

build {
  sources = [
    "source.nspawn-machine.basic-example"
  ]

  provisioner "shell-local" {
    inline = [
      "echo Hello from build container!",
    ]
  }
}
