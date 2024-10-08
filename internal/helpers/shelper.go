package shelper

import (
	"sort"

	"github.com/Masterminds/semver/v3"
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

// StringsToSemvers converts a list of version strings to semver.Version pointers,
// sorts them, and removes duplicates.
func StringsToSemvers(list []string) ([]*semver.Version, error) {
	var semvers []*semver.Version
	for _, raw := range list {
		v, err := semver.NewVersion(raw)
		if err != nil {
			return nil, err
		}
		semvers = append(semvers, v)
	}
	// Sort the semver versions
	sort.Sort(semver.Collection(semvers))
	// Remove duplicates
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
func StringsToStrings(list []string) ([]string, error) {
	semvers, err := StringsToSemvers(list)
	if err != nil {
		return nil, err
	}
	return SemversToStrings(semvers), nil
}

func PickFromSemverStrings(list []string, contraint string) ([]string, error) {
	var semvers_filtered []string
	semvers_list, err := StringsToSemvers(list)
	if err != nil {
		return nil, err
	}
	semver_compare, err := semver.NewConstraint(contraint)
	if err != nil {
		return nil, err
	}

	for _, v := range semvers_list {
		match := semver_compare.Check(v)
		// match, msgs := semver_compare.Validate(v)
		// for _, m := range msgs {
		// 	fmt.Println(m)
		// }
		if match {
			semvers_filtered = append(semvers_filtered, v.String())
		}
	}

	if len(semvers_filtered) == 0 {
		var empty_results []string
		return empty_results, nil
	}

	return semvers_filtered, nil
}
