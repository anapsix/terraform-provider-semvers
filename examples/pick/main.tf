terraform {
  required_providers {
    semvers = {
      source = "anapsix/semvers"
    }
  }
}

provider "semvers" {}

locals {
  versions = [
    "v0.1.0",
    "v0.1.10",
    "v0.2.0",
    "v0.2.1+hotfix",
    "v0.2.2-rc1",
    "v3.0.0",
  ]
}

output "semvers_list_picked" {
  value = [
    {
      constraint = "~> 0.1"
      list       = provider::semvers::pick(local.versions, "~> 0.1")
    },
    {
      constraint = "~> 0.2"
      list       = provider::semvers::pick(local.versions, "~> 0.2")
    },
    {
      constraint = "<= 2.0"
      list       = provider::semvers::pick(local.versions, "<= 2.0")
    },
    {
      constraint = "= 0.2.0"
      list       = provider::semvers::pick(local.versions, "= 0.2.0")
    },
  ]
}
