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
	_ function.Function = &SemversPickFunction{}
)

type SemversPickFunction struct {
	provider *semversProvider
}

func NewSemversPickFunction() function.Function {
	return &SemversPickFunction{}
}

func (r *SemversPickFunction) Metadata(_ context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "pick"
}

func (r *SemversPickFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Returns semver from list of semvers according to constraint",
		MarkdownDescription: "Returns semver from list of semvers according to constraint. When the provider `ignore_invalid_tags` setting is true, invalid semver strings in the list are skipped and logged at INFO level when TF_LOG is enabled.",
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

func (r *SemversPickFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var versions []string
	var constraint string

	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &versions, &constraint))
	if resp.Error != nil {
		return
	}

	ignoreInvalidTags := false
	if r.provider != nil {
		ignoreInvalidTags = r.provider.ignoreInvalidTags()
	}

	filtered_semvers, err := shelper.PickFromSemverStrings(ctx, versions, constraint, ignoreInvalidTags)

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
