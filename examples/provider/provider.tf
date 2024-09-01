terraform {
  required_providers {
    semvers = {
      source = "anapsix/semvers"
    }
  }
}

provider "semvers" {}

data "semvers_list" "example" {
  list = [
    "0.1",
    "v0.1.1+hotfix",
    "v0.2.0",
    "v0.10.0",
    "v5.0.1-rc1+43fbedb",
  ]
}
