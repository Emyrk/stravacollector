package stravalimit

import (
	"fmt"
	"io"
	"testing"
	"time"

	"github.com/rs/zerolog"

	"github.com/stretchr/testify/require"
)

func TestInterval(t *testing.T) {
	t.Parallel()

	ti := time.Now()
	start := GetInterval(ti)
	for i := 0; i < 50; i++ {
		ti = ti.Add(time.Minute * 15)

		next := GetInterval(ti)
		require.Truef(t, next == start+1, "next: %d, start: %d", next, start)
		start = next
	}
}

func TestDay(t *testing.T) {
	t.Parallel()

	ti := time.Now()
	start := GetDay(ti)
	for i := 0; i < 400; i++ {
		ti = ti.Add(time.Hour * 24)

		next := GetDay(ti)
		require.Truef(t, next != start, "next: %d, start: %d", next, start)
		start = next
	}
}

func TestTimeToMidnight(t *testing.T) {
	now := time.Now()
	fmt.Println(NextDailyReset(now))
}

//nolint:tparallel,paralleltest
func TestCanCall(t *testing.T) {
	testCases := []struct {
		Name                 string
		IntervalLimit        int64
		DailyLimit           int64
		CurrentIntervalUsage int64
		CurrentDailyUsage    int64
		Calls                int64
		IntervalBuffer       int64
		DailyBuffer          int64
		Expected             bool
	}{
		{
			Name:                 "Empty",
			IntervalLimit:        0,
			DailyLimit:           0,
			CurrentIntervalUsage: 0,
			CurrentDailyUsage:    0,
			Expected:             false,
			Calls:                1,
		},
		{
			Name:                 "OverBufferInterval",
			IntervalLimit:        200,
			DailyLimit:           2000,
			CurrentIntervalUsage: 170,
			CurrentDailyUsage:    100,
			Expected:             false,
			IntervalBuffer:       50,
			DailyBuffer:          200,
			Calls:                1,
		},
		{
			Name:                 "OverBufferDaily",
			IntervalLimit:        200,
			DailyLimit:           2000,
			CurrentIntervalUsage: 0,
			CurrentDailyUsage:    1800,
			Expected:             false,
			IntervalBuffer:       50,
			DailyBuffer:          500,
			Calls:                1,
		},
	}

	//nolint:
	for _, c := range testCases {
		t.Run(c.Name, func(t *testing.T) {
			now := time.Now()
			limiter.IntervalLimit = c.IntervalLimit
			limiter.DailyLimit = c.DailyLimit
			limiter.CurrentIntervalUsage = c.CurrentIntervalUsage
			limiter.CurrentDailyUsage = c.CurrentDailyUsage
			limiter.CurrentInterval = GetInterval(now)
			limiter.CurrentDay = GetDay(now)

			can, _ := CanLogger(c.Calls, c.IntervalBuffer, c.DailyBuffer, zerolog.New(io.Discard))
			require.Equal(t, c.Expected, can)
		})
	}

	limiter = New()
}
