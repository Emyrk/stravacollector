package river

import (
	"fmt"
	"time"

	"github.com/Emyrk/strava/strava/stravalimit"
	"github.com/rs/zerolog"
)

// jobStravaCheck checks our rate limits and decides if the # of calls can happen.
// This prevents us from hitting the Strava API rate limits.
func (m *Manager) jobStravaCheck(logger zerolog.Logger, calls int64, extraInterval, extraDaily int64) error {
	// We get 300 every interval (15min) and 3k a day
	iBuf, dBuf := int64(105), int64(605) // Buffer

	// Adjust if the next daily reset is close
	switch {
	case stravalimit.NextDailyReset(time.Now()) < time.Hour*1:
		iBuf, dBuf = int64(70), int64(200)
	case stravalimit.NextDailyReset(time.Now()) < time.Hour*3:
		iBuf, dBuf = int64(80), int64(400)
	}

	// Adjust
	iBuf -= extraInterval
	dBuf -= extraDaily
	if iBuf <= 0 {
		iBuf = 0
	}
	if dBuf <= 0 {
		dBuf = 0
	}

	ok, limitLogger := stravalimit.CanLogger(calls, iBuf, dBuf, logger)
	if !ok {
		m.rateLimitLogger.Do(func() {
			limitLogger.Error().
				Msg("hitting strava rate limit, job going to fail and try again later")
		})

		return fmt.Errorf("hitting strava rate limit, job going to fail and try again later")
	}
	return nil
}
