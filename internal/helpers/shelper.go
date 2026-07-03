package shelper

import (
	"context"
	"sort"

	"github.com/Masterminds/semver/v3"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Compare compares this version to another one. It returns -1, 0, or 1 if
// the version smaller, equal, or larger than the other version. 99 is returned
// on error.
//
// Versions are compared by X.Y.Z. Build metadata is ignored. Prerelease is
// lower than the version without a prerelease. Compare always takes into account
// prereleases. If you want to work with ranges using typical range syntaxes that
// skip prereleases if the range is not looking for them use constraints.
func Compare(a string, b string) (int, error) {
	version_a, err := semver.NewVersion(a)
	if err != nil {
		return 99, err
	}
	version_b, err := semver.NewVersion(b)
	if err != nil {
		return 99, err
	}
	result := version_a.Compare(version_b)
	return result, nil
}

func Equals(a string, b string) (bool, error) {
	compare, err := Compare(a, b)
	if err != nil || compare == 99 {
		return false, err
	}
	if compare == 0 {
		return true, nil
	}
	return false, nil
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
// sorts them, and removes duplicates.
func StringsToSemvers(ctx context.Context, list []string, ignoreInvalid bool) ([]*semver.Version, error) {
	var semvers []*semver.Version
	var skipped []string

	for _, raw := range list {
		v, err := semver.NewVersion(raw)
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
	sort.Sort(semver.Collection(semvers))
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

func PickFromSemverStrings(ctx context.Context, list []string, constraint string, ignoreInvalid bool) ([]string, error) {
	semver_compare, err := semver.NewConstraint(constraint)
	if err != nil {
		return nil, err
	}

	semvers, err := StringsToSemvers(ctx, list, ignoreInvalid)
	if err != nil {
		return nil, err
	}

	var semvers_filtered []string
	for _, v := range semvers {
		if semver_compare.Check(v) {
			semvers_filtered = append(semvers_filtered, v.String())
		}
	}

	if len(semvers_filtered) == 0 {
		return []string{}, nil
	}

	return semvers_filtered, nil
}

func Constrained(version string, constraint string) (bool, error) {
	c, err := semver.NewConstraint(constraint)
	if err != nil {
		return false, err
	}

	v, err := semver.NewVersion(version)
	if err != nil {
		return false, err
	}

	return c.Check(v), nil
}
