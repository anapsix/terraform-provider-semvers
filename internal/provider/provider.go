// Copyright (c) HashiCorp, Inc.
// Copyright (c) Anastas Dancha
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ provider.Provider              = &semversProvider{}
	_ provider.ProviderWithFunctions = &semversProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &semversProvider{
			version: version,
		}
	}
}

type semversProvider struct {
	version string
}

func (p *semversProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "semvers"
	resp.Version = p.version
}

// Schema defines the provider-level schema for configuration data.
func (p *semversProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Implements `semvers_list` data source and `sort` provider function, using `github.com/Masterminds/semver/v3`. Usage of `provider::semvers::sort([])` requires Terraform version 1.8 and above.",
		Attributes:  nil,
	}
}

func (p *semversProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		func() datasource.DataSource {
			return &semversListDataSource{}
		},
	}
}

func (p *semversProvider) Resources(ctx context.Context) []func() resource.Resource {
	// Return nil since this provider doesn't have resources.
	return nil
}

func (p *semversProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// No configuration necessary for this provider
}

func (p *semversProvider) Functions(_ context.Context) []func() function.Function {
	return []func() function.Function{
		NewSemversSortFunction,
		NewSemversPickFunction,
	}
}
