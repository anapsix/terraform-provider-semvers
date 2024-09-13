// Copyright (c) HashiCorp, Inc.
// Copyright (c) Anastas Dancha
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestSemversPickFunction_Known(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `output "semvers_filtered" {
          value = provider::semvers::pick(
           ["0.1.1-rc1+a231f59", "0.1.1", "0.1.10", "0.1.2-rc1", "0.2.1"],
           "~> 0.2"
          )
        }`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue(
						"semvers_filtered",
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("0.2.1"),
						}),
					),
				},
			},
			{
				Config: `output "semvers_filtered" {
          value = provider::semvers::pick(
            ["0.1.0", "0.1.1-rc1+a231f59", "0.1.1", "0.1.10", "0.1.2-rc1", "0.2.1"],
            ">= 0.1.1"
          )
        }`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue(
						"semvers_filtered",
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("0.1.1"),
							knownvalue.StringExact("0.1.10"),
							knownvalue.StringExact("0.2.1"),
						}),
					),
				},
			},
			{
				Config: `output "semvers_filtered" {
          value = provider::semvers::pick(
            ["0.1.0", "0.1.1-rc1+a231f59", "0.1.1", "0.1.10", "0.1.2-rc1", "0.2.1"],
            ">= 3.0"
          )
        }`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue(
						"semvers_filtered",
						knownvalue.ListSizeExact(0),
					),
				},
			},
			{
				Config: `output "semvers_filtered" {
          value = provider::semvers::pick([], ">= 3.0")
        }`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue(
						"semvers_filtered",
						knownvalue.ListSizeExact(0),
					),
				},
			},
		},
	})
}

func TestSemversPickFunction_Invalid(t *testing.T) {
	t.Parallel()

	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `output "result" {
          value = provider::semvers::pick(null, ">= 3.0")
        }`,
				ExpectError: regexp.MustCompile(`Invalid function argument`),
			},
			{
				Config: `output "results" {
          value = provider::semvers::pick(["blah", "0.1.1"], "~> 0.2")
        }`,
				ExpectError: regexp.MustCompile(`Invalid Semantic Version`),
			},
			{
				Config: `output "results" {
          value = provider::semvers::pick(["0.1.0", "0.1.1"], "~~> 0.2")
        }`,
				ExpectError: regexp.MustCompile(`improper constraint`),
			},
			{
				Config: `output "results" {
          value = provider::semvers::pick(true, "~> 0.2")
        }`,
				ExpectError: regexp.MustCompile(`Invalid function argument`),
			},
			{
				Config: `output "results" {
          value = provider::semvers::pick(["0.1.0", "0.1.1"], true)
        }`,
				ExpectError: regexp.MustCompile(`improper constraint`),
			},
			{
				Config: `output "results" {
          value = provider::semvers::pick(["0.1.0", "0.1.1"], {one: 1})
        }`,
				ExpectError: regexp.MustCompile(`Invalid function argument`),
			},
			{
				Config: `output "results" {
          value = provider::semvers::pick(null, "~> 0.2")
        }`,
				ExpectError: regexp.MustCompile(`Invalid function argument`),
			},
			{
				Config: `output "results" {
          value = provider::semvers::pick(["0.1.0", "0.1.1"], null)
        }`,
				ExpectError: regexp.MustCompile(`Invalid function argument`),
			},
		},
	})
}
