// Copyright (c) Anastas Dancha
// SPDX-License-Identifier: MPL-2.0

package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

// ProviderConfig holds provider-wide settings shared with data sources and functions.
type ProviderConfig struct {
	IgnoreInvalidTags bool
}

type semversProviderModel struct {
	IgnoreInvalidTags types.Bool `tfsdk:"ignore_invalid_tags"`
}
