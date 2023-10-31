package stravalimit

import (
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"

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

	// Gauges
	PromCurrentIntervalUsage prometheus.Gauge
	PromCurrentDailyUsage    prometheus.Gauge
	PromIntervalLimit        prometheus.Gauge
	PromDailyLimit           prometheus.Gauge
	PromCurrentDay           prometheus.Gauge
	PromCurrentInterval      prometheus.Gauge

	Registry *prometheus.Registry
	sync.Mutex
}

func New() *Limiter {
	now := time.Now()
	return &Limiter{
		CurrentInterval: GetInterval(now),
		CurrentDay:      GetDay(now),
		IntervalLimit:   200,
		DailyLimit:      1000,
		PromCurrentIntervalUsage: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "strava",
			Subsystem: "api_limiter",
			Name:      "interval_usage",
			Help:      "How many calls in this interval",
		}),
		PromCurrentDailyUsage: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "strava",
			Subsystem: "api_limiter",
			Name:      "daily_usage",
			Help:      "How many calls in this day",
		}),
		PromIntervalLimit: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "strava",
			Subsystem: "api_limiter",
			Name:      "interval_limit",
			Help:      "Interval limit",
		}),
		PromDailyLimit: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "strava",
			Subsystem: "api_limiter",
			Name:      "daily_limit",
			Help:      "Daily limit",
		}),
		PromCurrentDay: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "strava",
			Subsystem: "api_limiter",
			Name:      "current_day",
			Help:      "Current Day",
		}),
		PromCurrentInterval: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: "strava",
			Subsystem: "api_limiter",
			Name:      "current_interval",
			Help:      "Current Interval",
		}),
	}
}

func (l *Limiter) RegisterMetrics(reg *prometheus.Registry) {
	reg.MustRegister(
		l.PromCurrentDailyUsage, l.PromCurrentIntervalUsage,
		l.PromDailyLimit, l.PromIntervalLimit,
		l.PromCurrentInterval, l.PromCurrentDay,
	)
}

func SetRegistry(registry *prometheus.Registry) {
	limiter.Registry = registry
	limiter.RegisterMetrics(registry)
}

func GetInterval(t time.Time) int64 {
	return t.Unix() / (60 * 15)
}

func GetDay(t time.Time) int64 {
	return int64(t.UTC().YearDay())
}

func (l *Limiter) UpdateUsage(intervalUsage, intervalLimit, dailyUsage, dailyLimit int64) {
	l.Lock()
	defer l.Unlock()
	l.updateInterval()

	l.CurrentIntervalUsage = intervalUsage
	l.CurrentDailyUsage = dailyUsage
	l.IntervalLimit = intervalLimit
	l.DailyLimit = dailyLimit

	l.PromCurrentDailyUsage.Set(float64(l.CurrentDailyUsage))
	l.PromCurrentIntervalUsage.Set(float64(l.CurrentIntervalUsage))
	l.PromDailyLimit.Set(float64(l.DailyLimit))
	l.PromIntervalLimit.Set(float64(l.IntervalLimit))
	l.PromCurrentDay.Set(float64(l.CurrentDay))
	l.PromCurrentInterval.Set(float64(l.CurrentInterval))
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
		l.PromCurrentIntervalUsage.Set(0)
	}

	if l.CurrentDay != day {
		l.CurrentDay = day
		l.CurrentDailyUsage = 0
		l.PromCurrentDailyUsage.Set(0)
	}
}

func (l *Limiter) Update(headers http.Header) {
	if headers == nil {
		return
	}
	lInt, lDay := splitInts(headers.Get("X-Readratelimit-Limit"))
	uInt, uDay := splitInts(headers.Get("X-Readratelimit-Usage"))
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

// CanLogger
// Buffer is how many calls to reserve
func CanLogger(calls, intervalBuffer int64, dailyBuffer int64, logger zerolog.Logger) (bool, zerolog.Logger) {
	i, d := Remaining()

	if i < intervalBuffer+calls || d < dailyBuffer+calls {
		return false, logger.With().
			Int64("interval_remaining", i).
			Int64("daily_remaining", d).
			Int64("calls", calls).
			Int64("interval_buffer", intervalBuffer).
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
