package hugelhelp

import (
	"time"

	"github.com/Emyrk/strava/internal/hugeldate"
)

func HugelOngoing(now time.Time) bool {
	now = now.In(hugeldate.CentralTimeZone)
	if now.Month() == time.November && (now.Day() >= 7 && now.Day() <= 12) {
		return true
	}

	return false
}
