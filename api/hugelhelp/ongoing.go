package hugelhelp

import (
	"time"

	"github.com/Emyrk/strava/internal/hugeldate"
)

func HugelOngoing(now time.Time) bool {
	// Use the current leaderboard's event dates, including Friday and Sunday,
	// with an exclusive end. Comparisons work for both UTC and Central timestamps.
	return !now.Before(hugeldate.Year2026.Start) && now.Before(hugeldate.Year2026.End)
}
