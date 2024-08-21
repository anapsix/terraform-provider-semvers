// Copyright (c) HashiCorp, Inc.
// Copyright (c) Anastas Dancha
// SPDX-License-Identifier: MPL-2.0

package provider

import (
  "testing"

  "github.com/hashicorp/terraform-plugin-testing/knownvalue"
  "github.com/hashicorp/terraform-plugin-testing/statecheck"
  "github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

  "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)



func TestAccDataSource (t *testing.T) {
  resource.Test(t, resource.TestCase{
    PreCheck:                 func() { testAccPreCheck(t) },
    ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
    Steps: []resource.TestStep{
      {
        Config: `data "semvers_list" "example" {
            list = ["0.1.1-rc1+a231f59", "0.1.1", "0.1.10", "0.1.2-rc1"]
        }`,
        ConfigStateChecks: []statecheck.StateCheck{
          statecheck.ExpectKnownValue("data.semvers_list.example", tfjsonpath.New("sorted_versions").AtSliceIndex(0).AtMapKey("major"), knownvalue.Int64Exact(0)),
          statecheck.ExpectKnownValue("data.semvers_list.example", tfjsonpath.New("sorted_versions").AtSliceIndex(0).AtMapKey("minor"), knownvalue.Int64Exact(1)),
          statecheck.ExpectKnownValue("data.semvers_list.example", tfjsonpath.New("sorted_versions").AtSliceIndex(0).AtMapKey("patch"), knownvalue.Int64Exact(1)),
          statecheck.ExpectKnownValue("data.semvers_list.example", tfjsonpath.New("sorted_versions").AtSliceIndex(0).AtMapKey("prerelease"), knownvalue.StringExact("rc1")),
          statecheck.ExpectKnownValue("data.semvers_list.example", tfjsonpath.New("sorted_versions").AtSliceIndex(0).AtMapKey("metadata"), knownvalue.StringExact("a231f59")),
          statecheck.ExpectKnownValue("data.semvers_list.example", tfjsonpath.New("sorted_versions").AtSliceIndex(3).AtMapKey("major"), knownvalue.Int64Exact(0)),
          statecheck.ExpectKnownValue("data.semvers_list.example", tfjsonpath.New("sorted_versions").AtSliceIndex(3).AtMapKey("minor"), knownvalue.Int64Exact(1)),
          statecheck.ExpectKnownValue("data.semvers_list.example", tfjsonpath.New("sorted_versions").AtSliceIndex(3).AtMapKey("patch"), knownvalue.Int64Exact(10)),
          statecheck.ExpectKnownValue("data.semvers_list.example", tfjsonpath.New("sorted_versions").AtSliceIndex(3).AtMapKey("prerelease"), knownvalue.StringExact("")),
          statecheck.ExpectKnownValue("data.semvers_list.example", tfjsonpath.New("sorted_versions").AtSliceIndex(3).AtMapKey("metadata"), knownvalue.StringExact("")),
        },
      },
    },
  })
}
