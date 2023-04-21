package strava

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
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

func ParseRateLimit(resp http.Response) RateLimit {
	lInt, lDay := splitInts(resp.Header.Get("X-Ratelimit-Limit"))
	uInt, uDay := splitInts(resp.Header.Get("X-Ratelimit-Usage"))
	return RateLimit{
		DailyLimit:    lDay,
		DailyUsage:    uDay,
		IntervalLimit: lInt,
		IntervalUsage: uInt,
	}
}

func splitInts(s string) (int64, int64) {
	split := strings.Split(s, ",")
	if len(split) != 2 {
		return -1, -1
	}
	a, err := strconv.ParseInt(strings.TrimSpace(split[0]), 10, 64)
	if err != nil {
		a = -1
	}
	b, err := strconv.ParseInt(strings.TrimSpace(split[1]), 10, 64)
	if err != nil {
		b = -1
	}
	return a, b
}
