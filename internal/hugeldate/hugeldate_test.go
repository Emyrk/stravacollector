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
