// Copyright (c) HashiCorp, Inc.
// Copyright (c) Anastas Dancha
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/anapsix/terraform-provider-semvers/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ function.Function = SemversSortFunction{}
)

func NewSemversSortFunction() function.Function {
	return SemversSortFunction{}
}

type SemversSortFunction struct{}

func (r SemversSortFunction) Metadata(_ context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "sort"
}

func (r SemversSortFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Returns sorted list of semver strings",
		MarkdownDescription: "Returns sorted and deduped list of semver strings",
		Parameters: []function.Parameter{
			function.ListParameter{
				ElementType:         types.StringType,
				Name:                "versions",
				MarkdownDescription: "List of semver strings",
			},
		},
		Return: function.ListReturn{
			ElementType: types.StringType,
		},
	}
}

func (r SemversSortFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var versions []string

	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &versions))
	semvers, err := shelper.StringsToStrings(versions)

	if err != nil {
		tflog.Error(ctx, "Error in shelper.StringsToStrings()")
	}

	if resp.Error != nil {
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, semvers))
}
