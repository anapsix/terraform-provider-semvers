# add the following to ~/.terraformrc
# provider_installation {
#   dev_overrides {
#     "anapsix-dev/semvers" = "~/go/bin"
#   }
#   direct {}
# }

terraform {
  required_providers {
    semvers = {
      source = "anapsix-dev/semvers"
      # source = "anapsix/semvers"
      # version = "0.2.0"
    }
  }
}

provider "semvers" {}

data "semvers_list" "example" {
  list = local.versions
}

locals {
  version_prefix = "v"
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

  version_count = length(data.semvers_list.example.sorted_versions)

  local_first = data.semvers_list.example.sorted_versions[0]
  local_last  = data.semvers_list.example.sorted_versions[local.version_count - 1]

  local_last_no_prerelease = [for v in reverse(data.semvers_list.example.sorted_versions) : v if v.prerelease == ""][0]

  sorted_by_function = provider::semvers::sort(local.versions)

}

output "semvers_sorted_by_function" {
  value = local.sorted_by_function
}

output "semvers_list_sorted_versions" {
  value = data.semvers_list.example.sorted_versions
}

output "semvers_list_sorted_versions_dups" {
  value = [for v in data.semvers_list.example.sorted_versions : v if v["version"] == "1.0.0"]
}

output "semvers_list_first" {
  value = data.semvers_list.example.first
}

output "first_local" {
  value = local.local_first
}

output "semvers_list_last" {
  value = data.semvers_list.example.last
}

output "last_local" {
  value = local.local_last
}

output "last_local_noprerelease" {
  value = local.local_last_no_prerelease
}
