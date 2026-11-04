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
	_ function.Function = &SemversPickFunction{}
)

type SemversPickFunction struct{}

func NewSemversPickFunction() function.Function {
	return &SemversPickFunction{}
}

func (r *SemversPickFunction) Metadata(_ context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "pick"
}

func (r *SemversPickFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Returns semver from list of semvers according to constraint",
		MarkdownDescription: "Returns semver from list of semvers according to constraint. Pass `true` as the optional third argument to skip invalid semver strings in the list. The provider `ignore_invalid_tags` setting does not apply to functions.",
		Parameters: []function.Parameter{
			function.ListParameter{
				ElementType:         types.StringType,
				Name:                "versions",
				MarkdownDescription: "List of semver strings",
			},
			function.StringParameter{
				Name:                "constraint",
				MarkdownDescription: "Semver constraint",
			},
		},
		VariadicParameter: ignoreInvalidParameter,
		Return: function.ListReturn{
			ElementType: types.StringType,
		},
	}
}

func (r *SemversPickFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var versions []types.String
	var constraint string
	var opts []bool

	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &versions, &constraint, &opts))
	if resp.Error != nil {
		return
	}

	ignoreInvalidTags, funcErr := ignoreInvalidFromArgs(opts, 2)
	if funcErr != nil {
		resp.Error = funcErr
		return
	}

	filtered, err := semverutil.PickFromSemverStrings(ctx, versionStrings(versions), constraint, ignoreInvalidTags)

	if err != nil {
		resp.Error = function.NewArgumentFuncError(errorPosition(err, 0, 1), err.Error())
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Result.Set(ctx, filtered))
}
