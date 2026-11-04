// Copyright (c) Anastas Dancha
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/function"
)

// ignoreInvalidParameter is the optional trailing argument accepted by list functions.
// Provider configuration is not available to provider-defined functions, so leniency
// must be requested per call.
var ignoreInvalidParameter = function.BoolParameter{
	Name:                "ignore_invalid",
	MarkdownDescription: "Optional. When `true`, invalid semver strings in the list are skipped instead of returning an error. Skipped strings are logged at INFO level when `TF_LOG` is enabled. Defaults to `false`.",
}

// ignoreInvalidFromArgs resolves the variadic `ignore_invalid` argument, which must be
// passed at most once. argPos is the position of the variadic argument.
func ignoreInvalidFromArgs(opts []bool, argPos int) (bool, *function.FuncError) {
	switch len(opts) {
	case 0:
		return false, nil
	case 1:
		return opts[0], nil
	default:
		return false, function.NewArgumentFuncError(int64(argPos), fmt.Sprintf("ignore_invalid accepts at most one value, got %d", len(opts)))
	}
}
