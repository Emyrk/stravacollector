package superlative

import (
	"time"

	"github.com/Emyrk/strava/api/modelsdk"
)

type List struct {
	// EarlyBird is the earliest activity start
	EarlyBird Entry[time.Time] `json:"early_bird"`
	// NightOwl is the latest end time
	NightOwl Entry[time.Time] `json:"night_owl"`
	// MostStoppage value is in seconds
	MostStoppage Entry[int64] `json:"most_stoppage"`
	// LeastStoppage value is in seconds
	LeastStoppage Entry[int64]   `json:"least_stoppage"`
	MostWatts     Entry[float64] `json:"most_watts"`
	MostCadence   Entry[float64] `json:"most_cadence"`
	LeastCadence  Entry[float64] `json:"least_cadence"`

	//Shortest ride = most efficient
	//Longest ride =
	//Most elevation = mountain climber
	//Max amount of PRs = most improved
	//Suffer score = most pain
	//Best of $Segment
}

type Entry[T any] struct {
	Activity modelsdk.StringInt `json:"activity_id"`
	Value    T                  `json:"value"`
}
