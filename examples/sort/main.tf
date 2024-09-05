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
    "1",
    "1.0",
    "1.0.0",
    "v1",
    "v1.0",
    "v1.0.0",
    "v0.90.1",
    "9",
    "9.0.0",
    "v0.9.0",
    "v0.9",
    "v0.80.0",
    "v0.91.0",
    "v0.9.10",
    "v1.0.1",
    "5.0.1-rc1+dead",
    "v2+test",
    "v2.0.0-rc1",
    "v2.0.1-rc1",
    "v0.1.0",
  ]
}

locals {
  sorted_by_function = provider::semvers::sort(local.versions)
}

output "semvers_list_sorted_versions" {
  value = local.sorted_by_function
}
