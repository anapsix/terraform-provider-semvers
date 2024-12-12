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

func TestSemversConstrainedFunction_Known(t *testing.T) {
	t.Parallel()

	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `output "results_one" {
          value = provider::semvers::constrained("0.1.1", ">= 0.1")
        }
        output "results_two" {
          value = provider::semvers::constrained("0.1", ">= 0.1.1")
        }
        output "results_three" {
          value = provider::semvers::constrained("0.1.2", "~> 0.1")
        }
        output "results_four" {
          value = provider::semvers::constrained("0.2", "~> 0.1")
        }
        output "results_five" {
          value = provider::semvers::constrained("0.2", "0.1.1 - 0.2.0")
        }`,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue(
						"results_one",
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownOutputValue(
						"results_two",
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownOutputValue(
						"results_three",
						knownvalue.Bool(true),
					),
					statecheck.ExpectKnownOutputValue(
						"results_four",
						knownvalue.Bool(false),
					),
					statecheck.ExpectKnownOutputValue(
						"results_five",
						knownvalue.Bool(true),
					),
				},
			},
		},
	})
}

func TestSemversConstrainedFunction_Invalid(t *testing.T) {
	t.Parallel()

	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `output "results_one" {
          value = provider::semvers::constrained("0.1.1", "blah")
        }`,
				ExpectError: regexp.MustCompile(`improper constraint`),
			},
		},
	})
}
