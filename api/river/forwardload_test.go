package river

import (
	"github.com/Emyrk/strava/strava"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestForwardLoadMaxAttempts(t *testing.T) {
	require.Equal(t, 2, (ForwardLoadArgs{}).InsertOpts().MaxAttempts)
}

func TestOnHugelDate2026(t *testing.T) {
	for _, tc := range []struct {
		date string
		want bool
	}{
		{"2025-11-08T12:00:00Z", false},
		{"2026-11-13T05:59:59Z", false},
		{"2026-11-13T06:00:00Z", true},
		{"2026-11-14T12:00:00Z", true},
		{"2026-11-16T05:59:59Z", true},
		{"2026-11-16T06:00:00Z", false},
	} {
		t.Run(tc.date, func(t *testing.T) {
			start, err := time.Parse(time.RFC3339, tc.date)
			require.NoError(t, err)
			require.Equal(t, tc.want, onHugelDate(strava.ActivitySummary{StartDate: start}))
		})
	}
}
