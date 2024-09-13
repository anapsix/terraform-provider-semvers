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
	_ function.Function = SemversPickFunction{}
)

func NewSemversPickFunction() function.Function {
	return SemversPickFunction{}
}

type SemversPickFunction struct{}

func (r SemversPickFunction) Metadata(_ context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "pick"
}

func (r SemversPickFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Returns semver from list of semvers according to constraint",
		MarkdownDescription: "Returns semver from list of semvers according to constraint",
		Parameters: []function.Parameter{
			function.ListParameter{
				AllowNullValue:      false,
				AllowUnknownValues:  false,
				ElementType:         types.StringType,
				Name:                "versions",
				MarkdownDescription: "List of semver strings",
			},
			function.StringParameter{
				AllowNullValue:      false,
				AllowUnknownValues:  false,
				Name:                "constraint",
				MarkdownDescription: "Semver constraint",
			},
		},
		Return: function.ListReturn{
			ElementType: types.StringType,
		},
	}
}

func (r SemversPickFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var versions []string
	var constraint string

	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &versions, &constraint))
	if resp.Error != nil {
		return
	}

	filtered_semvers, err := shelper.PickFromSemverStrings(versions, constraint)

	if err != nil {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError("Error performing operation: "+err.Error()))
		return
	}

	if len(filtered_semvers) == 0 {
		empty_list := make([]string, 0)
		resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, empty_list))
	} else {
		resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, filtered_semvers))
	}
}
