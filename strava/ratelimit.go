package strava

import (
	"fmt"
)

type RateLimit struct {
	DailyLimit int64
	DailyUsage int64

	// Intervals are every 15min
	IntervalLimit int64
	IntervalUsage int64
}

func (r RateLimit) String() string {
	return fmt.Sprintf("Daily: %d/%d | Interval: %d/%d", r.DailyUsage, r.DailyLimit, r.IntervalUsage, r.IntervalLimit)
}
