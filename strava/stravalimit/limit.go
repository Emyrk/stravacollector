package stravalimit

import (
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

var limiter *Limiter = New()

type Limiter struct {
	CurrentInterval int64
	CurrentDay      int64

	// Last known usage
	CurrentIntervalUsage int64
	CurrentDailyUsage    int64

	// Last known limits
	IntervalLimit int64
	DailyLimit    int64

	sync.Mutex
}

func New() *Limiter {
	now := time.Now()
	return &Limiter{
		CurrentInterval: GetInterval(now),
		CurrentDay:      GetDay(now),
		IntervalLimit:   200,
		DailyLimit:      1000,
	}
}

func GetInterval(t time.Time) int64 {
	return t.Unix() / (60 * 15)
}

func GetDay(t time.Time) int64 {
	return int64(t.YearDay())
}

func (l *Limiter) UpdateUsage(intervalUsage, intervalLimit, dailyUsage, dailyLimit int64) {
	l.Lock()
	defer l.Unlock()
	l.updateInterval()

	l.CurrentIntervalUsage = intervalUsage
	l.CurrentDailyUsage = dailyUsage
	l.IntervalLimit = intervalLimit
	l.DailyLimit = dailyLimit
}

func (l *Limiter) Remaining() (int64, int64) {
	l.Lock()
	defer l.Unlock()
	l.updateInterval()

	return l.IntervalLimit - l.CurrentIntervalUsage, l.DailyLimit - l.CurrentDailyUsage
}

func (l *Limiter) updateInterval() {
	now := time.Now()
	interval := GetInterval(now)
	day := GetDay(now)

	if l.CurrentInterval != interval {
		l.CurrentInterval = interval
		l.CurrentIntervalUsage = 0
	}

	if l.CurrentDay != day {
		l.CurrentDay = day
		l.CurrentDailyUsage = 0
	}
}

func (l *Limiter) Update(headers http.Header) {
	if headers == nil {
		return
	}
	lInt, lDay := splitInts(headers.Get("X-RateLimit-Limit"))
	uInt, uDay := splitInts(headers.Get("X-RateLimit-Usage"))
	if lInt == -1 || lDay == -1 || uInt == -1 || uDay == -1 {
		return
	}

	l.UpdateUsage(uInt, lInt, uDay, lDay)
}

func Update(headers http.Header) {
	limiter.Update(headers)
}

func UpdateUsage(intervalUsage, intervalLimit, dailyUsage, dailyLimit int64) {
	limiter.UpdateUsage(intervalUsage, intervalLimit, dailyUsage, dailyLimit)
}

func Remaining() (int64, int64) {
	return limiter.Remaining()
}

func CanLogger(calls, buffer int64, dailyBuffer int64, logger zerolog.Logger) (bool, zerolog.Logger) {
	i, d := Remaining()

	if i < buffer+calls || d < dailyBuffer+calls {
		return false, logger.With().
			Int64("interval_remaining", i).
			Int64("daily_remaining", d).
			Int64("calls", calls).
			Int64("interval_buffer", buffer).
			Int64("daily_buffer", dailyBuffer).

			// Remove
			Int64("interval_limit", limiter.IntervalLimit).
			Int64("daily_limit", limiter.DailyLimit).
			Int64("interval_usage", limiter.CurrentIntervalUsage).
			Int64("daily_usage", limiter.CurrentDailyUsage).
			Logger()
	}
	return true, logger
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

func NextDailyReset(now time.Time) time.Duration {
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	midnight = midnight.AddDate(0, 0, 1)
	return midnight.Sub(now)
}
