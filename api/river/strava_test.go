package river

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestHugelEventDailyReserve(t *testing.T) {
	for _, tc := range []struct {
		instant            string
		potential, boosted bool
	}{
		{"2026-11-14T05:59:59Z", true, false},
		{"2026-11-14T06:00:00Z", true, true},
		{"2026-11-14T16:44:00Z", true, true},
		{"2026-11-15T05:59:59Z", true, true},
		{"2026-11-15T06:00:00Z", true, false},
		{"2027-11-14T12:00:00Z", true, false},
		{"2026-11-14T16:44:00Z", false, false},
	} {
		t.Run(tc.instant+"/"+map[bool]string{true: "potential", false: "ordinary"}[tc.potential], func(t *testing.T) {
			now, err := time.Parse(time.RFC3339, tc.instant)
			require.NoError(t, err)
			cap := hugelEventDailyReserve(now, tc.potential)
			if tc.boosted {
				require.Equal(t, []int64{100}, cap)
			} else {
				require.Nil(t, cap)
			}
		})
	}
}

func TestHugelAllowancePreservesIntervalAndOrdinaryPolicy(t *testing.T) {
	for _, tc := range []struct {
		hour            int
		interval, daily int64
	}{
		{16, 100, 555}, {22, 75, 350}, {23, 65, 150},
	} {
		now := time.Date(2026, 11, 14, tc.hour, 30, 0, 0, time.UTC)
		i, d := jobStravaBuffers(now, 5, 50)
		require.Equal(t, tc.interval, i)
		require.Equal(t, tc.daily, d)
		boostedI, boostedD := jobStravaBuffers(now, 5, 50, hugelEventDailyReserve(now, true)...)
		require.Equal(t, i, boostedI)
		require.Equal(t, int64(100), boostedD)
		ordinaryI, ordinaryD := jobStravaBuffers(now, 0, 0, hugelEventDailyReserve(now, false)...)
		require.Equal(t, i+5, ordinaryI)
		require.Equal(t, d+50, ordinaryD)
	}
	// Preserve the existing lower manual reserve near reset rather than increasing it.
	now := time.Date(2026, 11, 14, 23, 30, 0, 0, time.UTC)
	_, d := jobStravaBuffers(now, 10, 115, 100)
	require.Equal(t, int64(85), d)
	// The same instant in another location must not alter reset proximity.
	central := time.FixedZone("CST", -6*60*60)
	i, d := jobStravaBuffers(now.In(central), 5, 50, 100)
	require.Equal(t, int64(65), i)
	require.Equal(t, int64(100), d)
}
