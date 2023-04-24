package stravalimit

import (
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
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

	sync.RWMutex
}

func New() *Limiter {
	now := time.Now()
	return &Limiter{
		CurrentInterval: GetInterval(now),
		CurrentDay:      GetDay(now),
		IntervalLimit:   100,
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
	l.RLock()
	defer l.RUnlock()
	l.updateInterval()

	return l.IntervalLimit - l.CurrentIntervalUsage, l.DailyLimit - l.CurrentDailyUsage
}

func (l *Limiter) updateInterval() {
	l.Lock()
	defer l.Unlock()

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

	l.UpdateUsage(uInt, uDay, lInt, lDay)
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

func Can(buffer int64) bool {
	i, d := Remaining()
	if i < buffer || d < buffer {
		return false
	}
	return true
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
