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

func TestSemversCompareFunction_Known(t *testing.T) {
	t.Parallel()

	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `output "results_one" {
          value = provider::semvers::compare("0.1.1", "0.1.1")
        }
        output "results_two" {
          value = provider::semvers::compare("0.1.0", "0.1.0+test")
        }`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue(
						"results_one",
						knownvalue.Int32Exact(0),
					),
					statecheck.ExpectKnownOutputValue(
						"results_two",
						knownvalue.Int32Exact(0),
					),
				},
			},
			{
				Config: `output "results" {
          value = provider::semvers::compare("0.1.0", "0.1.1")
        }`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue(
						"results",
						knownvalue.Int32Exact(-1),
					),
				},
			},
			{
				Config: `output "results" {
          value = provider::semvers::compare("0.1.1", "0.1.0")
        }`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue(
						"results",
						knownvalue.Int32Exact(1),
					),
				},
			},
		},
	})
}

func TestSemversCompareFunction_Invalid(t *testing.T) {
	t.Parallel()

	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `output "results" {
          value = provider::semvers::compare("blah", "0.1.0")
        }`,
				ExpectError: regexp.MustCompile(`Invalid Semantic Version`),
			},
			{
				Config: `output "results" {
          value = provider::semvers::compare("0.1.0", "blah")
        }`,
				ExpectError: regexp.MustCompile(`Invalid Semantic Version`),
			},
		},
	})
}
