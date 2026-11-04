package semverutil

import (
	"context"
	"regexp"
	"testing"

	"github.com/Masterminds/semver/v3"
)

func mustConstraint(t *testing.T, s string) *semver.Constraints {
	t.Helper()
	c, err := ParseConstraint(s)
	if err != nil {
		t.Fatalf("ParseConstraint(%q): %v", s, err)
	}
	return c
}

func TestPickFromSemverStrings_ignoreInvalid(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	got, err := PickFromSemverStrings(ctx, []string{"develop-latest", "1.2.3-dev.1", "0.1.0"}, mustConstraint(t, "~> 1.2.3-dev"), true)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"1.2.3-dev.1"}
	if len(got) != len(want) {
		t.Fatalf("got %d results, want %d: %v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("got[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}

func TestPickFromSemverStrings_strictInvalid(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	_, err := PickFromSemverStrings(ctx, []string{"develop-latest", "1.2.3-dev.1"}, mustConstraint(t, "~> 1.2.3-dev"), false)
	if err == nil {
		t.Fatal("expected error for invalid semver tag, got nil")
	}

	if !regexp.MustCompile(`(?i)invalid semantic version`).MatchString(err.Error()) {
		t.Fatalf("unexpected error: %v", err)
	}
}

// Equal versions must keep their input order so dedup keeps the first occurrence.
func TestStringsToSemvers_keepsFirstDuplicate(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	list := []string{
		"v1.0.0", "1.0.0", "1.0", "v1", "1",
		"0.9.0", "2.0.0", "v0.9.0", "0.9", "v2", "2",
		"1.0.0", "v1.0", "0.9.0", "2.0", "v1.0.0",
		"3.0.0", "0.1.0", "v3", "3",
	}

	got, err := StringsToSemvers(ctx, list, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := map[string]string{
		"0.1.0": "0.1.0",
		"0.9.0": "0.9.0",
		"1.0.0": "v1.0.0",
		"2.0.0": "2.0.0",
		"3.0.0": "3.0.0",
	}
	if len(got) != len(want) {
		t.Fatalf("got %d versions, want %d: %v", len(got), len(want), SemversToStrings(got))
	}
	for _, v := range got {
		if v.Original() != want[v.String()] {
			t.Errorf("version %s kept original %q, want %q", v, v.Original(), want[v.String()])
		}
	}
}
