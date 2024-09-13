// Copyright (c) HashiCorp, Inc.
// Copyright (c) Anastas Dancha
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	shelper "github.com/anapsix/terraform-provider-semvers/internal/helpers"
	"github.com/hashicorp/terraform-plugin-framework/function"
)

var (
	_ function.Function = SemversCompareFunction{}
)

func NewSemversCompareFunction() function.Function {
	return SemversCompareFunction{}
}

type SemversCompareFunction struct{}

func (r SemversCompareFunction) Metadata(_ context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "compare"
}

func (r SemversCompareFunction) Definition(_ context.Context, _ function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Returns comparison results as integer",
		MarkdownDescription: `Returns comparison results as integer
    <ul>
    <li><code>-1</code> if smaller</li>
    <li><code>0</code> if equals</li>
    <li><code>-1</code> if larger</li>
    <li><code>99</code> if error</li>
    </ul>
    <br>Versions are compared by X.Y.Z. Build metadata is ignored. Prerelease is
    lower than the version without a prerelease. Compare always takes into account
    prereleases.`,
		Parameters: []function.Parameter{
			function.StringParameter{
				AllowNullValue:      false,
				AllowUnknownValues:  false,
				Name:                "version_a",
				MarkdownDescription: "Semver string version used as base for comparison",
			},
			function.StringParameter{
				AllowNullValue:      false,
				AllowUnknownValues:  false,
				Name:                "version_b",
				MarkdownDescription: "Semver string version to compare to `version_a`",
			},
		},
		Return: function.Int32Return{},
	}
}

func (r SemversCompareFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var version_a string
	var version_b string
	var compare_results int

	resp.Error = function.ConcatFuncErrors(req.Arguments.Get(ctx, &version_a, &version_b))
	if resp.Error != nil {
		return
	}

	compare_results, err := shelper.Compare(version_a, version_b)

	if err != nil || compare_results == 99 {
		resp.Error = function.ConcatFuncErrors(resp.Error, function.NewFuncError("Error performing operation: "+err.Error()))
		return
	}

	resp.Error = function.ConcatFuncErrors(resp.Error, resp.Result.Set(ctx, compare_results))
}
