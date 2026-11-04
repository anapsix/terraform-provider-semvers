// Copyright (c) HashiCorp, Inc.
// Copyright (c) Anastas Dancha
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/anapsix/terraform-provider-semvers/internal/semverutil"
	"github.com/hashicorp/terraform-plugin-framework/function"
)

var (
	_ function.Function = SemversEqualsFunction{}
)

func NewSemversEqualsFunction() function.Function {
	return SemversEqualsFunction{}
}

type SemversEqualsFunction struct{}

func (r SemversEqualsFunction) Metadata(_ context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "equals"
}

func (r SemversEqualsFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Checks the equality of two semver strings, returns a boolean",
		MarkdownDescription: `Checks the equality of two semver strings, returns a boolean.
		<br><br>~> **NOTE:** Versions are compared by X.Y.Z. Build metadata is ignored. Prerelease is
    lower than the version without a prerelease. Compare always takes into account
    prereleases. See [Masterminds/semver](https://github.com/Masterminds/semver).`,
		Parameters: []function.Parameter{
			function.StringParameter{
				Name:                "version_a",
				MarkdownDescription: "Semver string version used as base for comparison",
			},
			function.StringParameter{
				Name:                "version_b",
				MarkdownDescription: "Semver string version to compare to `version_a`",
			},
		},
		Return: function.BoolReturn{},
	}
}

func (r SemversEqualsFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var versionA string
	var versionB string

	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &versionA, &versionB))
	if resp.Error != nil {
		return
	}

	a, err := semverutil.ParseVersion(versionA)
	if err != nil {
		resp.Error = function.NewArgumentFuncError(0, err.Error())
		return
	}
	b, err := semverutil.ParseVersion(versionB)
	if err != nil {
		resp.Error = function.NewArgumentFuncError(1, err.Error())
		return
	}

	resp.Error = resp.Result.Set(ctx, a.Equal(b))
}
