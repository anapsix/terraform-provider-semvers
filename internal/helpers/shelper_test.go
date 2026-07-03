package shelper

import (
	"context"
	"regexp"
	"testing"
)

func TestPickFromSemverStrings_ignoreInvalid(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	got, err := PickFromSemverStrings(ctx, []string{"develop-latest", "1.2.3-dev.1", "0.1.0"}, "~> 1.2.3-dev", true)
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

	_, err := PickFromSemverStrings(ctx, []string{"develop-latest", "1.2.3-dev.1"}, "~> 1.2.3-dev", false)
	if err == nil {
		t.Fatal("expected error for invalid semver tag, got nil")
	}

	if !regexp.MustCompile(`(?i)invalid semantic version`).MatchString(err.Error()) {
		t.Fatalf("unexpected error: %v", err)
	}
}
