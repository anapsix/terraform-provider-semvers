// Copyright (c) Anastas Dancha
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/anapsix/terraform-provider-semvers/internal/semverutil"
	"github.com/hashicorp/terraform-plugin-framework/function"
)

var (
	_ function.Function = SemversConstrainedFunction{}
)

func NewSemversConstrainedFunction() function.Function {
	return SemversConstrainedFunction{}
}

type SemversConstrainedFunction struct{}

func (r SemversConstrainedFunction) Metadata(_ context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "constrained"
}

func (r SemversConstrainedFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Checks whether semver is within constraint, returns a boolean",
		MarkdownDescription: "Checks whether semver is within constraint, returns a boolean",
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "version",
				MarkdownDescription: "Semver string version used as base for comparison",
			},
			function.StringParameter{
				Name:                "constraint",
				MarkdownDescription: "Constraint string",
			},
		},
		Return: function.BoolReturn{},
	}
}

func (r SemversConstrainedFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var version string
	var constraint string

	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &version, &constraint))
	if resp.Error != nil {
		return
	}

	result, err := semverutil.Constrained(version, constraint)

	if err != nil {
		resp.Error = function.NewArgumentFuncError(errorPosition(err, 0, 1), err.Error())
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, result))
}
