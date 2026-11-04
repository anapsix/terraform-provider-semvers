// Copyright (c) HashiCorp, Inc.
// Copyright (c) Anastas Dancha
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/anapsix/terraform-provider-semvers/internal/semverutil"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ function.Function = &SemversSortFunction{}
)

type SemversSortFunction struct{}

func NewSemversSortFunction() function.Function {
	return &SemversSortFunction{}
}

func (r *SemversSortFunction) Metadata(_ context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "sort"
}

func (r *SemversSortFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Returns sorted list of semver strings",
		MarkdownDescription: "Returns sorted and deduped list of semver strings. Pass `true` as the optional second argument to skip invalid semver strings in the list. The provider `ignore_invalid_tags` setting does not apply to functions.",
		Parameters: []function.Parameter{
			function.ListParameter{
				ElementType:         types.StringType,
				Name:                "versions",
				MarkdownDescription: "List of semver strings",
			},
		},
		VariadicParameter: ignoreInvalidParameter,
		Return: function.ListReturn{
			ElementType: types.StringType,
		},
	}
}

func (r *SemversSortFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var versions []types.String
	var opts []bool

	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &versions, &opts))
	if resp.Error != nil {
		return
	}

	ignoreInvalidTags, funcErr := ignoreInvalidFromArgs(opts, 1)
	if funcErr != nil {
		resp.Error = funcErr
		return
	}

	semvers, err := semverutil.StringsToStrings(ctx, versionStrings(versions), ignoreInvalidTags)

	if err != nil {
		resp.Error = function.NewArgumentFuncError(0, err.Error())
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, semvers))
}
