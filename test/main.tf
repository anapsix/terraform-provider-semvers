# add the following to ~/.terraformrc
# provider_installation {
#   dev_overrides {
#     "anapsix/semvers" = "~/go/bin"
#   }
#   direct {}
# }

terraform {
  required_providers {
    semvers = {
      source = "anapsix/semvers"
    }
  }
}

provider "semvers" {}

# Providing a custom list as input (required)
data "semvers_list" "example" {
  provider = semvers
  list = ["0.90.1", "0.9.0", "0.80.0", "0.91.0", "0.9.10", "1.0.1", "2.0.0", "2.0.0-rc1", "2.0.1-rc1", "0.1", "5.0.1-rc1+dead"]
}

locals {
  version_count = length(data.semvers_list.example.sorted_versions)

  local_first = data.semvers_list.example.sorted_versions[0]
  local_last  = data.semvers_list.example.sorted_versions[local.version_count - 1]

  local_last_no_prerelease = [for v in reverse(data.semvers_list.example.sorted_versions): v if v.prerelease == ""][0]
}

output "sorted_versions" {
  value = data.semvers_list.example.sorted_versions
}

output "first" {
  value = data.semvers_list.example.first
}

output "first_local" {
  value = local.local_first
}

output "last" {
  value = data.semvers_list.example.last
}

output "last_local" {
  value = local.local_last
}

output "last_no_prerelease_local" {
  value = local.local_last
}


