data "nspawn-images" "test" {
  mock = "mock"
}

locals {
  machines = data.nspawn-images.test.machines
  images = data.nspawnn-images.test.images
}

source "null" "basic-example" {
  communicator = "none"
}

build {
  sources = [
    "source.null.basic-example"
  ]

  provisioner "shell-local" {
    inline = [
      "echo machines: ${local.machines}",
      "echo images: ${local.images}",
    ]
  }
}
