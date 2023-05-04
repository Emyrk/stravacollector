package modelsdk

import (
	"strconv"
	"time"
)

// Int64String is used because javascript can't handle 64 bit integers
func Int64String(v int64) string {
	return strconv.FormatInt(v, 10)
}

type AthleteLogin struct {
	AthleteID string `db:"athlete_id" json:"athlete_id"`
	Summit    bool   `db:"summit" json:"summit"`
}

// AthleteSummary is the smallest amount of information we need to know about an athlete
// to show them on a page.
type AthleteSummary struct {
	AthleteID            string    `db:"athlete_id" json:"athlete_id"`
	Summit               bool      `json:"summit"`
	Username             string    `json:"username"`
	Firstname            string    `json:"firstname"`
	Lastname             string    `json:"lastname"`
	Sex                  string    `json:"sex"`
	ProfilePicLink       string    `json:"profile_pic_link"`
	ProfilePicLinkMedium string    `json:"profile_pic_link_medium"`
	UpdatedAt            time.Time `db:"updated_at" json:"updated_at"`
}

type HugelLeaderBoard struct {
	PersonalBest *HugelLeaderBoardActivity  `json:"personal_best,omitempty"`
	Activities   []HugelLeaderBoardActivity `json:"activities"`
}

type HugelLeaderBoardActivity struct {
	RankOneElapsed int64           `json:"rank_one_elapsed"`
	ActivityID     string          `json:"activity_id"`
	AthleteID      string          `json:"athlete_id"`
	Elapsed        int64           `json:"elapsed"`
	Rank           int64           `json:"rank"`
	Efforts        []SegmentEffort `json:"efforts"`
	Athlete        MinAthlete      `json:"athlete"`

	// Activity info
	ActivityName               string    `json:"activity_name"`
	ActivityDistance           float64   `json:"activity_distance"`
	ActivityMovingTime         int64     `json:"activity_moving_time"`
	ActivityElapsedTime        int64     `json:"activity_elapsed_time"`
	ActivityStartDate          time.Time `json:"activity_start_date"`
	ActivityTotalElevationGain float64   `json:"activity_total_elevation_gain"`
}

type SegmentEffort struct {
	EffortID     string    `json:"effort_id"`
	StartDate    time.Time `json:"start_date"`
	SegmentID    string    `json:"segment_id"`
	ElapsedTime  int64     `json:"elapsed_time"`
	MovingTime   int64     `json:"moving_time"`
	DeviceWatts  bool      `json:"device_watts"`
	AverageWatts float64   `json:"average_watts"`
}

type MinAthlete struct {
	AthleteID      string `db:"athlete_id" json:"athlete_id"`
	Username       string `json:"username"`
	Firstname      string `json:"firstname"`
	Lastname       string `json:"lastname"`
	Sex            string `json:"sex"`
	ProfilePicLink string `json:"profile_pic_link"`
}
