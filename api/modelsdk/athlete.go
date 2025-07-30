package modelsdk

import (
	"time"

	"github.com/Emyrk/strava/api/superlative"
)

type AthleteHugelActivity struct {
	Summary          ActivitySummary `json:"summary"`
	Efforts          []SegmentEffort `json:"efforts"`
	TotalTimeSeconds int64           `json:"total_time_seconds"`
}

type SyncActivitySummary struct {
	Activity ActivitySummary `json:"activity_summary"`
	Synced   bool            `json:"synced"`
	SyncedAt time.Time       `json:"synced_at"`
}

type AthleteSyncSummary struct {
	Load             AthleteLoad           `json:"athlete_load"`
	TotalActivities  int                   `json:"total_activities"`
	SyncedActivities []SyncActivitySummary `json:"synced_activities"`
	Athlete          AthleteSummary        `json:"athlete_summary"`

	TotalSummary int `json:"total_summary"`
	TotalDetail  int `json:"total_detail"`
}

// Tracks loading athlete activities. Must be an authenticated athlete.
type AthleteLoad struct {
	AthleteID int64 `json:"athlete_id"`
	// Timestamp start of the last activity loaded. Future ones are not loaded.
	ActivityTimeAfter time.Time `json:"activity_time_after"`
	// Timestamp of the last time the athlete was attempted to be loaded.
	LastLoadAttempt time.Time `json:"last_load_attempt"`
	// True if the last load was completed no more work is needed to catch up.
	LastLoadComplete  bool      `json:"last_load_complete"`
	NextLoadNotBefore time.Time `json:"next_load_not_before"`
}

type AthleteLogin struct {
	AthleteID StringInt `db:"athlete_id" json:"athlete_id"`
	Summit    bool      `db:"summit" json:"summit"`
}

type AthleteHugelActivities struct {
	Activities []AthleteHugelActivity `json:"activities"`
}

type ActivitySummary struct {
	ActivityID     StringInt `db:"activity_id" json:"activity_id"`
	AthleteID      StringInt `db:"athlete_id" json:"athlete_id"`
	UploadID       StringInt `db:"upload_id" json:"upload_id"`
	ExternalID     string    `db:"external_id" json:"external_id"`
	Name           string    `db:"name" json:"name"`
	Distance       float64   `db:"distance" json:"distance"`
	MovingTime     float64   `db:"moving_time" json:"moving_time"`
	ElapsedTime    float64   `db:"elapsed_time" json:"elapsed_time"`
	TotalEleGain   float64   `db:"total_elevation_gain" json:"total_elevation_gain"`
	ActivityType   string    `db:"activity_type" json:"activity_type"`
	SportType      string    `db:"sport_type" json:"sport_type"`
	StartDate      time.Time `db:"start_date" json:"start_date"`
	StartDateLocal time.Time `db:"start_date_local" json:"start_date_local"`
	Timezone       string    `db:"timezone" json:"timezone"`
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
	Superlatives superlative.List           `json:"superlatives"`
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
	ActivitySufferScore        int       `json:"activity_suffer_score"`
	ActivityAchievementCount   int       `json:"activity_achievement_count"`
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
