// Copyright (c) HashiCorp, Inc.
// Copyright (c) Anastas Dancha
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"

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
			config:  &ProviderConfig{},
		}
	}
}

type semversProvider struct {
	version string
	config  *ProviderConfig
}

func (p *semversProvider) ignoreInvalidTags() bool {
	if p.config == nil {
		return false
	}
	return p.config.IgnoreInvalidTags
}

func (p *semversProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "semvers"
	resp.Version = p.version
}

func (p *semversProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Implements `semvers_list` data source, `compare`, `constrained`, `equals`, `sort` and `pick` provider function. All functionality is based on `github.com/Masterminds/semver/v3`. Usage of provider functions requires Terraform version 1.8 and above.",
		Attributes: map[string]schema.Attribute{
			"ignore_invalid_tags": schema.BoolAttribute{
				Optional:            true,
				Description:         "When true, invalid semver strings in list inputs are skipped instead of returning an error. Applies to `semvers_list`, `sort`, and `pick`.",
				MarkdownDescription: "When true, invalid semver strings in list inputs are skipped instead of returning an error. Applies to `semvers_list`, `sort`, and `pick`. Skipped tags are logged at INFO level when `TF_LOG` is enabled. Defaults to `false`.",
			},
		},
	}
}

func (p *semversProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		func() datasource.DataSource {
			return &semversListDataSource{}
		},
	}
}

func (p *semversProvider) Resources(_ context.Context) []func() resource.Resource {
	return nil
}

func (p *semversProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var config semversProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	p.config = &ProviderConfig{
		IgnoreInvalidTags: false,
	}
	if !config.IgnoreInvalidTags.IsNull() && !config.IgnoreInvalidTags.IsUnknown() {
		p.config.IgnoreInvalidTags = config.IgnoreInvalidTags.ValueBool()
	}

	resp.DataSourceData = p.config
}

func (p *semversProvider) Functions(_ context.Context) []func() function.Function {
	return []func() function.Function{
		func() function.Function {
			return &SemversSortFunction{provider: p}
		},
		func() function.Function {
			return &SemversPickFunction{provider: p}
		},
		NewSemversCompareFunction,
		NewSemversEqualsFunction,
		NewSemversConstrainedFunction,
	}
}

func providerConfigFromDataSourceConfigure(req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) *ProviderConfig {
	if req.ProviderData == nil {
		return &ProviderConfig{}
	}

	providerConfig, ok := req.ProviderData.(*ProviderConfig)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Provider Data Type",
			fmt.Sprintf("Expected *provider.ProviderConfig, got %T", req.ProviderData),
		)
		return &ProviderConfig{}
	}

	return providerConfig
}
