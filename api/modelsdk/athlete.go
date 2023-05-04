package modelsdk

import (
	"time"
)

type AthleteLogin struct {
	AthleteID StringInt `db:"athlete_id" json:"athlete_id"`
	Summit    bool      `db:"summit" json:"summit"`
}

// AthleteSummary is the smallest amount of information we need to know about an athlete
// to show them on a page.
type AthleteSummary struct {
	AthleteID            StringInt `db:"athlete_id" json:"athlete_id"`
	Summit               bool      `json:"summit"`
	Username             string    `json:"username"`
	Firstname            string    `json:"firstname"`
	Lastname             string    `json:"lastname"`
	Sex                  string    `json:"sex"`
	ProfilePicLink       string    `json:"profile_pic_link"`
	ProfilePicLinkMedium string    `json:"profile_pic_link_medium"`
	UpdatedAt            time.Time `json:"updated_at"`
	HugelCount           int       `json:"hugel_count"`
}

type HugelLeaderBoard struct {
	PersonalBest *HugelLeaderBoardActivity  `json:"personal_best,omitempty"`
	Activities   []HugelLeaderBoardActivity `json:"activities"`
}

type HugelLeaderBoardActivity struct {
	RankOneElapsed int64           `json:"rank_one_elapsed"`
	ActivityID     StringInt       `json:"activity_id"`
	AthleteID      StringInt       `json:"athlete_id"`
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

type SuperHugelLeaderBoard struct {
	PersonalBest *SuperHugelLeaderBoardActivity  `json:"personal_best,omitempty"`
	Activities   []SuperHugelLeaderBoardActivity `json:"activities"`
}

type SuperHugelLeaderBoardActivity struct {
	RankOneElapsed int64           `json:"rank_one_elapsed"`
	AthleteID      StringInt       `json:"athlete_id"`
	Elapsed        int64           `json:"elapsed"`
	Rank           int64           `json:"rank"`
	Efforts        []SegmentEffort `json:"efforts"`
	Athlete        MinAthlete      `json:"athlete"`
}

type SegmentEffort struct {
	ActivityID   StringInt `json:"activity_id"`
	EffortID     StringInt `json:"effort_id"`
	StartDate    time.Time `json:"start_date"`
	SegmentID    StringInt `json:"segment_id"`
	ElapsedTime  int64     `json:"elapsed_time"`
	MovingTime   int64     `json:"moving_time"`
	DeviceWatts  bool      `json:"device_watts"`
	AverageWatts float64   `json:"average_watts"`
}

type MinAthlete struct {
	AthleteID      StringInt `db:"athlete_id" json:"athlete_id"`
	Username       string    `json:"username"`
	Firstname      string    `json:"firstname"`
	Lastname       string    `json:"lastname"`
	Sex            string    `json:"sex"`
	ProfilePicLink string    `json:"profile_pic_link"`
	HugelCount     int       `json:"hugel_count"`
}
