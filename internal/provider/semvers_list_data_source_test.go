// Copyright (c) HashiCorp, Inc.
// Copyright (c) Anastas Dancha
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestAccDataSource(t *testing.T) {
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
			{
				Config: `data "semvers_list" "dedup" {
            list = ["1", "1.0", "1.0.0", "v1", "v1.0", "v1.0.0"]
        }`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.semvers_list.dedup", tfjsonpath.New("sorted_versions"), knownvalue.ListSizeExact(1)),
				},
			},
		},
	})
}

func TestAccDataSource_IgnoreInvalidTags(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `data "semvers_list" "strict" {
            list = ["develop-latest", "0.1.0"]
        }`,
				ExpectError: regexp.MustCompile(`Error parsing semver list`),
			},
			{
				Config: `provider "semvers" {
            ignore_invalid_tags = true
        }

        data "semvers_list" "lenient" {
            list = ["develop-latest", "0.2.0", "0.1.0"]
        }`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.semvers_list.lenient", tfjsonpath.New("sorted_versions"), knownvalue.ListSizeExact(2)),
					statecheck.ExpectKnownValue("data.semvers_list.lenient", tfjsonpath.New("last").AtMapKey("version"), knownvalue.StringExact("0.2.0")),
				},
			},
			{
				Config: `provider "semvers" {
            ignore_invalid_tags = true
        }

        data "semvers_list" "none_valid" {
            list = ["develop-latest", "not-a-version"]
        }`,
				ExpectError: regexp.MustCompile(`No valid semver strings`),
			},
		},
	})
}

// ignore_invalid_tags can be unknown during plan (e.g. derived from a resource).
// The data source must not silently fall back to strict parsing in that case.
func TestAccDataSource_IgnoreInvalidTagsUnknown(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `resource "terraform_data" "flag" {
            input = true
        }

        provider "semvers" {
            ignore_invalid_tags = terraform_data.flag.output
        }

        data "semvers_list" "valid_only" {
            list = ["0.2.0", "0.1.0"]
        }`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("data.semvers_list.valid_only", tfjsonpath.New("last").AtMapKey("version"), knownvalue.StringExact("0.2.0")),
				},
			},
			{
				Config: `resource "terraform_data" "flag" {
            input = true
        }

        # flag2 is new in this step, so its output is unknown during plan;
        # flag already exists from step 1 and would be known.
        resource "terraform_data" "flag2" {
            input = true
        }

        provider "semvers" {
            ignore_invalid_tags = terraform_data.flag2.output
        }

        data "semvers_list" "with_invalid" {
            list = ["develop-latest", "0.2.0", "0.1.0"]
        }`,
				ExpectError: regexp.MustCompile(`ignore_invalid_tags\s+is\s+not\s+known`),
			},
		},
	})
}
