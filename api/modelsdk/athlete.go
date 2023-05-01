package modelsdk

import "time"

type AthleteLogin struct {
	AthleteID int64 `db:"athlete_id" json:"athlete_id"`
	Summit    bool  `db:"summit" json:"summit"`
}

// AthleteSummary is the smallest amount of information we need to know about an athlete
// to show them on a page.
type AthleteSummary struct {
	AthleteID            int64     `db:"athlete_id" json:"athlete_id"`
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
	ActivityID int64           `json:"activity_id"`
	AthleteID  int64           `json:"athlete_id"`
	Elapsed    int64           `json:"elapsed"`
	Rank       int64           `json:"rank"`
	Efforts    []SegmentEffort `json:"efforts"`
	Athlete    MinAthlete      `json:"athlete"`
}

type SegmentEffort struct {
	EffortID     int64     `json:"effort_id"`
	StartDate    time.Time `json:"start_date"`
	SegmentID    int64     `json:"segment_id"`
	ElapsedTime  int64     `json:"elapsed_time"`
	MovingTime   int64     `json:"moving_time"`
	DeviceWatts  bool      `json:"device_watts"`
	AverageWatts float64   `json:"average_watts"`
}

type MinAthlete struct {
	AthleteID      int64  `db:"athlete_id" json:"athlete_id"`
	Username       string `json:"username"`
	Firstname      string `json:"firstname"`
	Lastname       string `json:"lastname"`
	Sex            string `json:"sex"`
	ProfilePicLink string `json:"profile_pic_link"`
}
