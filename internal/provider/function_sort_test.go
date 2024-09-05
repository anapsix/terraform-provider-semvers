// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestSemversSortFunction_Known(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: `
        locals {
          list = ["0.1.1-rc1+a231f59", "0.1.1", "0.1.10", "0.1.2-rc1"]
        }
        output "semvers_sorted" {
          value = provider::semvers::sort(local.list)
        }
        `,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue(
						"semvers_sorted",
						knownvalue.ListExact([]knownvalue.Check{
							knownvalue.StringExact("0.1.1-rc1+a231f59"),
							knownvalue.StringExact("0.1.1"),
							knownvalue.StringExact("0.1.2-rc1"),
							knownvalue.StringExact("0.1.10"),
						}),
					),
				},
			},
			{
				Config: `
        locals {
          list = ["2","2.0", "2.0.0", "v2", "v2.0", "v2.0.0"]
        }
        output "semvers_deduped" {
          value = provider::semvers::sort(local.list)
        }
        `,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue(
						"semvers_deduped",
						knownvalue.ListSizeExact(1),
					),
				},
			},
			{
				Config: `
        locals {
          list = ["2+abc", "2.0.0", "v2", "v2.0+abc"]
        }
        output "semvers_deduped" {
          value = provider::semvers::sort(local.list)
        }
        `,
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownOutputValue(
						"semvers_deduped",
						knownvalue.ListSizeExact(2),
					),
				},
			},
		},
	})
}

// func TestSemversSortFunction_Null(t *testing.T) {
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

// func TestSemversSortFunction_Unknown(t *testing.T) {
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
