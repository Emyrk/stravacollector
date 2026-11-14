package hugeldate

import (
	"testing"
	"time"
)

func TestYear2026(t *testing.T) {
	expectedStart := time.Date(2026, time.November, 13, 0, 0, 0, 0, CentralTimeZone)
	if !Year2026.Start.Equal(expectedStart) {
		t.Fatalf("unexpected 2026 start: got %s, want %s", Year2026.Start, expectedStart)
	}

	expectedEnd := time.Date(2026, time.November, 16, 0, 0, 0, 0, CentralTimeZone)
	if !Year2026.End.Equal(expectedEnd) {
		t.Fatalf("unexpected 2026 end: got %s, want %s", Year2026.End, expectedEnd)
	}
}

// Use UTC instants rather than the location under test: a UTC fallback must fail.
func TestYear2026CentralBoundaries(t *testing.T) {
	if got, want := Year2026.Start.UTC(), time.Date(2026, 11, 13, 6, 0, 0, 0, time.UTC); !got.Equal(want) {
		t.Fatalf("start = %s, want %s", got, want)
	}
	if got, want := Year2026.End.UTC(), time.Date(2026, 11, 16, 6, 0, 0, 0, time.UTC); !got.Equal(want) {
		t.Fatalf("end = %s, want %s", got, want)
	}
}
