// Copyright (c) HashiCorp, Inc.
// Copyright (c) Anastas Dancha
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	shelper "github.com/anapsix/terraform-provider-semvers/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

	if resp.Error != nil {
		return
	}

	semvers, err := shelper.StringsToStrings(versions)

	if err != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError("Error performing operation: "+err.Error()))
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, semvers))
}
