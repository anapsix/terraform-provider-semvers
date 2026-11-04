// Copyright (c) Anastas Dancha
// SPDX-License-Identifier: MPL-2.0

package provider

import "github.com/hashicorp/terraform-plugin-framework/types"

// ProviderConfig holds provider-wide settings shared with data sources.
type ProviderConfig struct {
	IgnoreInvalidTags bool
	// IgnoreInvalidTagsUnknown is true when ignore_invalid_tags is not known yet,
	// e.g. during plan when it is derived from a resource attribute.
	IgnoreInvalidTagsUnknown bool
}

type semversProviderModel struct {
	IgnoreInvalidTags types.Bool `tfsdk:"ignore_invalid_tags"`
}
