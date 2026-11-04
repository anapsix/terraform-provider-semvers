package semverutil

import (
	"context"
	"fmt"
	"sort"

	"github.com/Masterminds/semver/v3"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// InvalidVersionError reports a string that cannot be parsed as a semver.
type InvalidVersionError struct {
	Version string
	Err     error
}

func (e *InvalidVersionError) Error() string {
	return fmt.Sprintf("invalid version %q: %s", e.Version, e.Err)
}

func (e *InvalidVersionError) Unwrap() error {
	return e.Err
}

// InvalidConstraintError reports a string that cannot be parsed as a semver constraint.
type InvalidConstraintError struct {
	Constraint string
	Err        error
}

func (e *InvalidConstraintError) Error() string {
	return fmt.Sprintf("invalid constraint %q: %s", e.Constraint, e.Err)
}

func (e *InvalidConstraintError) Unwrap() error {
	return e.Err
}

// ParseVersion parses s as a semver, returning an *InvalidVersionError on failure.
func ParseVersion(s string) (*semver.Version, error) {
	v, err := semver.NewVersion(s)
	if err != nil {
		return nil, &InvalidVersionError{Version: s, Err: err}
	}
	return v, nil
}

// ParseConstraint parses s as a semver constraint, returning an
// *InvalidConstraintError on failure.
func ParseConstraint(s string) (*semver.Constraints, error) {
	c, err := semver.NewConstraint(s)
	if err != nil {
		return nil, &InvalidConstraintError{Constraint: s, Err: err}
	}
	return c, nil
}

// RemoveDups removes duplicate versions from a list of semver.Version pointers
func RemoveDups(list []*semver.Version) []*semver.Version {
	seen := make(map[string]struct{})
	var result []*semver.Version
	for _, v := range list {
		versionStr := v.String()
		if _, exists := seen[versionStr]; !exists {
			seen[versionStr] = struct{}{} // Using struct{}{} to save space
			result = append(result, v)
		}
	}
	return result
}

func logSkippedInvalidTags(ctx context.Context, skipped []string) {
	if len(skipped) == 0 {
		return
	}

	tflog.Info(ctx, "Skipping invalid semver tags", map[string]interface{}{
		"count": len(skipped),
		"tags":  skipped,
	})
}

// StringsToSemvers converts a list of version strings to semver.Version pointers,
// sorts them, and removes duplicates. Sorting is stable, so of several equal
// versions the one appearing first in list is kept.
func StringsToSemvers(ctx context.Context, list []string, ignoreInvalid bool) ([]*semver.Version, error) {
	var semvers []*semver.Version
	var skipped []string

	for _, raw := range list {
		v, err := ParseVersion(raw)
		if err != nil {
			if ignoreInvalid {
				skipped = append(skipped, raw)
				continue
			}
			return nil, err
		}
		semvers = append(semvers, v)
	}

	logSkippedInvalidTags(ctx, skipped)
	sort.Stable(semver.Collection(semvers))
	return RemoveDups(semvers), nil
}

// SemversToStrings converts a list of semver.Version pointers to version strings
func SemversToStrings(semversList []*semver.Version) []string {
	semverStrings := make([]string, len(semversList))
	for i, v := range semversList {
		semverStrings[i] = v.String()
	}
	return semverStrings
}

// StringsToStrings converts a list of version strings to a sorted and deduplicated
// list of version strings
func StringsToStrings(ctx context.Context, list []string, ignoreInvalid bool) ([]string, error) {
	semvers, err := StringsToSemvers(ctx, list, ignoreInvalid)
	if err != nil {
		return nil, err
	}
	return SemversToStrings(semvers), nil
}

func PickFromSemverStrings(ctx context.Context, list []string, c *semver.Constraints, ignoreInvalid bool) ([]string, error) {
	semvers, err := StringsToSemvers(ctx, list, ignoreInvalid)
	if err != nil {
		return nil, err
	}

	filtered := []string{}
	for _, v := range semvers {
		if c.Check(v) {
			filtered = append(filtered, v.String())
		}
	}

	return filtered, nil
}
