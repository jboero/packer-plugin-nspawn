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
