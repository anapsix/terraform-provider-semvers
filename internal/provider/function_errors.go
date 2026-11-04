// Copyright (c) Anastas Dancha
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"errors"

	"github.com/anapsix/terraform-provider-semvers/internal/semverutil"
)

// errorPosition returns the position of the function argument responsible for err:
// constraintPos when err reports an invalid constraint, versionPos otherwise.
func errorPosition(err error, versionPos, constraintPos int64) int64 {
	var cErr *semverutil.InvalidConstraintError
	if errors.As(err, &cErr) {
		return constraintPos
	}
	return versionPos
}

// versionPairPosition returns the position (0 or 1) of the argument holding the
// invalid version reported by err, for functions taking two version arguments.
func versionPairPosition(err error, versionA string) int64 {
	var vErr *semverutil.InvalidVersionError
	if errors.As(err, &vErr) && vErr.Version != versionA {
		return 1
	}
	return 0
}
