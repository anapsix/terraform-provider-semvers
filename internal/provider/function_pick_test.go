// Copyright (c) HashiCorp, Inc.
// Copyright (c) Anastas Dancha
// SPDX-License-Identifier: MPL-2.0

package provider

import (
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
				Config: `
        output "semvers_filtered" {
          value = provider::semvers::pick(
						["0.1.1-rc1+a231f59", "0.1.1", "0.1.10", "0.1.2-rc1", "0.2.1"],
						"~> 0.2"
					)
        }
        `,
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
				Config: `
        output "semvers_filtered" {
          value = provider::semvers::pick(
						["0.1.0", "0.1.1-rc1+a231f59", "0.1.1", "0.1.10", "0.1.2-rc1", "0.2.1"],
						">= 0.1.1"
					)
        }
        `,
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
				Config: `
        output "semvers_filtered" {
          value = provider::semvers::pick(
						["0.1.0", "0.1.1-rc1+a231f59", "0.1.1", "0.1.10", "0.1.2-rc1", "0.2.1"],
						">= 3.0"
					)
        }
        `,
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

// func TPick_Null(t *testing.T) {
//   resource.UnitTest(t, resource.TestCase{
//     TerraformVersionChecks: []tfversion.TerraformVersionCheck{
//       tfversion.SkipBelow(tfversion.Version1_8_0),
//     },
//     ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
//     Steps: []resource.TestStep{
//       {
//         Config: `
//         output "test" {
//           value = provider::semvers::sort_semvers(null)
//         }
//         `,
//         // The parameter does not enable AllowNullValue
//         ExpectError: regexp.MustCompile(`argument must not be null`),
//       },
//     },
//   })
// }

// func TPick_Unknown(t *testing.T) {
//   resource.UnitTest(t, resource.TestCase{
//     TerraformVersionChecks: []tfversion.TerraformVersionCheck{
//       tfversion.SkipBelow(tfversion.Version1_8_0),
//     },
//     ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
//     Steps: []resource.TestStep{
//       {
//         Config: `
//         resource "terraform_data" "test" {
//           input = "testvalue"
//         }

//         output "test" {
//           value = provider::semvers::sort_semvers(terraform_data.test.output)
//         }
//         `,
//         Check: resource.ComposeAggregateTestCheckFunc(
//           resource.TestCheckOutput("test", "testvalue"),
//         ),
//       },
//     },
//   })
// }
