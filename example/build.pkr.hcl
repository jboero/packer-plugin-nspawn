packer {
  required_plugins {
    nspawn = {
      version = ">=v0.1.0"
      source  = "github.com/hashicorp/nspawn"
    }
  }
}

source "nspawn" "foo-example" {
  mock = local.foo
}

source "nspawn" "bar-example" {
  mock = local.bar
}

build {
  sources = [
    "source.nspawn.foo-example",
  ]

  source "source.nspawn.bar-example" {
    name = "bar"
  }
}
