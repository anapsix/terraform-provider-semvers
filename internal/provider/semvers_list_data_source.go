// Copyright (c) Anastas Dancha
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"sort"

	"github.com/Masterminds/semver/v3"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Data Source definition
type semversListDataSource struct{}

var versionAttributes = map[string]schema.Attribute{
	"version": schema.StringAttribute{
		Computed: true,
	},
	"major": schema.Int64Attribute{
		Computed: true,
	},
	"minor": schema.Int64Attribute{
		Computed: true,
	},
	"patch": schema.Int64Attribute{
		Computed: true,
	},
	"prerelease": schema.StringAttribute{
		Computed: true,
	},
	"metadata": schema.StringAttribute{
		Computed: true,
	},
}

func (d *semversListDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_list"
}

func (d *semversListDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"list": schema.ListAttribute{
				ElementType: types.StringType,
				Required:    true,
			},
			"first": schema.SingleNestedAttribute{
				Computed:   true,
				Attributes: versionAttributes,
			},
			"last": schema.SingleNestedAttribute{
				Computed:   true,
				Attributes: versionAttributes,
			},
			"sorted_versions": schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: versionAttributes,
				},
				Computed: true,
			},
		},
	}
}

func (d *semversListDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var model semversListDataSourceModel

	// Get the data from the request
	diags := req.Config.Get(ctx, &model)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Ensure that the list of semver strings is provided
	if model.List.IsNull() || model.List.IsUnknown() || len(model.List.Elements()) == 0 {
		resp.Diagnostics.AddError(
			"Missing Semver List",
			"The 'list' attribute must contain at least one semver string.",
		)
		return
	}

	// Parse and sort the semver strings
	semverStrings := convertTerraformListToStringSlice(model.List)
	semvers := make([]*semver.Version, len(semverStrings))
	for i, raw := range semverStrings {
		v, err := semver.NewVersion(raw)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error parsing semver",
				"Could not parse version: "+raw,
			)
			return
		}
		semvers[i] = v
	}
	sort.Sort(semver.Collection(semvers))

	// Prepare the sorted list of semver strings
	sortedSemvers := make([]map[string]interface{}, len(semvers))
	for i, v := range semvers {
		sortedSemvers[i] = map[string]interface{}{
			"version":    v.String(),
			"major":      v.Major(),
			"minor":      v.Minor(),
			"patch":      v.Patch(),
			"prerelease": v.Prerelease(),
			"metadata":   v.Metadata(),
		}
	}

	// Set the data source state
	convertedValues := convertToTerraformValueList(ctx, sortedSemvers)

	model.First = convertedValues[0].(types.Object)
	model.Last = convertedValues[len(convertedValues)-1].(types.Object)

	model.SortedVersions = types.ListValueMust(types.ObjectType{
		AttrTypes: map[string]attr.Type{
			"version":    types.StringType,
			"major":      types.Int64Type,
			"minor":      types.Int64Type,
			"patch":      types.Int64Type,
			"prerelease": types.StringType,
			"metadata":   types.StringType,
		},
	}, convertedValues)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read a data source")

	resp.Diagnostics.Append(resp.State.Set(ctx, &model)...)
}

func convertToTerraformValueList(ctx context.Context, values []map[string]interface{}) []attr.Value {
	result := make([]attr.Value, len(values))

	for i, v := range values {
		// Define the attribute types for the object
		attrTypes := map[string]attr.Type{
			"version":    types.StringType,
			"major":      types.Int64Type,
			"minor":      types.Int64Type,
			"patch":      types.Int64Type,
			"prerelease": types.StringType,
			"metadata":   types.StringType,
		}

		// Create a map to hold the attribute values
		attrValues := map[string]attr.Value{}
		for mk, mv := range v {
			switch val := mv.(type) {
			case string:
				attrValues[mk] = types.StringValue(val)
			case uint64:
				attrValues[mk] = types.Int64Value(int64(val))
			case int64:
				attrValues[mk] = types.Int64Value(val)
			case bool:
				attrValues[mk] = types.BoolValue(val)
			}
		}

		// Use types.ObjectValueMust with the correct arguments
		result[i] = types.ObjectValueMust(attrTypes, attrValues)
	}

	return result
}

func convertTerraformListToStringSlice(list types.List) []string {
	result := make([]string, len(list.Elements()))
	for i, v := range list.Elements() {
		result[i] = v.(types.String).ValueString()
	}
	return result
}

// Data source model
type semversListDataSourceModel struct {
	List           types.List   `tfsdk:"list"`
	First          types.Object `tfsdk:"first"`
	Last           types.Object `tfsdk:"last"`
	SortedVersions types.List   `tfsdk:"sorted_versions"`
}
