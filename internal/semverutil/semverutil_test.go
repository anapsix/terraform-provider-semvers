// Copyright (c) Anastas Dancha
// SPDX-License-Identifier: MPL-2.0

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

func TestInvalidConstraintError_message(t *testing.T) {
	t.Parallel()

	tests := []struct {
		constraint string
		want       string
	}{
		// the underlying error already quotes the whole constraint
		{constraint: "~~> 0.2", want: `improper constraint: "~~> 0.2"`},
		// the underlying error quotes only the offending part
		{constraint: "1.0 || ~~1", want: `invalid constraint "1.0 || ~~1": improper constraint: " ~~1"`},
	}
	for _, tt := range tests {
		_, err := ParseConstraint(tt.constraint)
		if err == nil {
			t.Fatalf("ParseConstraint(%q): expected error", tt.constraint)
		}
		if err.Error() != tt.want {
			t.Errorf("ParseConstraint(%q) error = %q, want %q", tt.constraint, err.Error(), tt.want)
		}
	}
}

// Versions differing only in build metadata compare equal but are not duplicates.
func TestStringsToStrings_keepsBuildMetadataVariants(t *testing.T) {
	t.Parallel()

	got, err := StringsToStrings(context.Background(), []string{"1.0.0+b", "v1.0.0+a", "1.0.0+b"}, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	want := []string{"1.0.0+b", "1.0.0+a"}
	if len(got) != len(want) || got[0] != want[0] || got[1] != want[1] {
		t.Fatalf("got %v, want %v", got, want)
	}
}
