package hugelhelp

import (
	"testing"
	"time"
)

func TestHugelOngoing2026(t *testing.T) {
	for _, tc := range []struct {
		name, timestamp string
		want            bool
	}{
		{"old event", "2025-11-08T12:00:00-06:00", false},
		{"old November days", "2026-11-08T12:00:00-06:00", false},
		{"before", "2026-11-13T05:59:59Z", false},
		{"start", "2026-11-13T06:00:00Z", true},
		{"event day", "2026-11-14T12:00:00-06:00", true},
		{"last instant", "2026-11-16T05:59:59Z", true},
		{"end", "2026-11-16T06:00:00Z", false},
		{"next year", "2027-11-14T12:00:00-06:00", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			now, err := time.Parse(time.RFC3339, tc.timestamp)
			if err != nil {
				t.Fatal(err)
			}
			if got := HugelOngoing(now); got != tc.want {
				t.Fatalf("HugelOngoing(%s) = %v, want %v", now, got, tc.want)
			}
		})
	}
}
