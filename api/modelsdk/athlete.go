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
