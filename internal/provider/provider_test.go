// Copyright (c) HashiCorp, Inc.
// Copyright (c) Anastas Dancha
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"semvers": func() (tfprotov6.ProviderServer, error) {
		return providerserver.NewProtocol6WithError(New("test")())()
	},
}

func testAccProtoV6ProviderFactoriesIgnoreInvalidTags(ignoreInvalidTags bool) map[string]func() (tfprotov6.ProviderServer, error) {
	return map[string]func() (tfprotov6.ProviderServer, error){
		"semvers": func() (tfprotov6.ProviderServer, error) {
			p := New("test")().(*semversProvider)
			p.config.IgnoreInvalidTags = ignoreInvalidTags
			return providerserver.NewProtocol6WithError(p)()
		},
	}
}

func testAccPreCheck(t *testing.T) {
}
