// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2

package database

import (
	"database/sql"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type Activity struct {
	ID        int64 `db:"id" json:"id"`
	AthleteID int64 `db:"athlete_id" json:"athlete_id"`
	UploadID  int64 `db:"upload_id" json:"upload_id"`
	// External ID refers to external source of the activity.
	ExternalID               string    `db:"external_id" json:"external_id"`
	Name                     string    `db:"name" json:"name"`
	MovingTime               float64   `db:"moving_time" json:"moving_time"`
	ElapsedTime              float64   `db:"elapsed_time" json:"elapsed_time"`
	TotalElevationGain       float64   `db:"total_elevation_gain" json:"total_elevation_gain"`
	ActivityType             string    `db:"activity_type" json:"activity_type"`
	SpotType                 string    `db:"spot_type" json:"spot_type"`
	StartDate                time.Time `db:"start_date" json:"start_date"`
	StartDateLocal           time.Time `db:"start_date_local" json:"start_date_local"`
	Timezone                 string    `db:"timezone" json:"timezone"`
	UtcOffset                int32     `db:"utc_offset" json:"utc_offset"`
	StartLatlng              []float64 `db:"start_latlng" json:"start_latlng"`
	EndLatlng                []float64 `db:"end_latlng" json:"end_latlng"`
	AchievementCount         int32     `db:"achievement_count" json:"achievement_count"`
	KudosCount               int32     `db:"kudos_count" json:"kudos_count"`
	CommentCount             int32     `db:"comment_count" json:"comment_count"`
	AthleteCount             int32     `db:"athlete_count" json:"athlete_count"`
	PhotoCount               int32     `db:"photo_count" json:"photo_count"`
	MapID                    string    `db:"map_id" json:"map_id"`
	MapPolyline              string    `db:"map_polyline" json:"map_polyline"`
	MapSummaryPolyline       string    `db:"map_summary_polyline" json:"map_summary_polyline"`
	Trainer                  bool      `db:"trainer" json:"trainer"`
	Commute                  bool      `db:"commute" json:"commute"`
	Manual                   bool      `db:"manual" json:"manual"`
	Private                  bool      `db:"private" json:"private"`
	Flagged                  bool      `db:"flagged" json:"flagged"`
	GearID                   string    `db:"gear_id" json:"gear_id"`
	FromAcceptedTag          bool      `db:"from_accepted_tag" json:"from_accepted_tag"`
	AverageSpeed             float64   `db:"average_speed" json:"average_speed"`
	MaxSpeed                 float64   `db:"max_speed" json:"max_speed"`
	AverageCadence           float64   `db:"average_cadence" json:"average_cadence"`
	AverageTemp              float64   `db:"average_temp" json:"average_temp"`
	AverageWatts             float64   `db:"average_watts" json:"average_watts"`
	WeightedAverageWatts     float64   `db:"weighted_average_watts" json:"weighted_average_watts"`
	Kilojoules               float64   `db:"kilojoules" json:"kilojoules"`
	DeviceWatts              bool      `db:"device_watts" json:"device_watts"`
	HasHeartrate             bool      `db:"has_heartrate" json:"has_heartrate"`
	MaxWatts                 float64   `db:"max_watts" json:"max_watts"`
	ElevHigh                 float64   `db:"elev_high" json:"elev_high"`
	ElevLow                  float64   `db:"elev_low" json:"elev_low"`
	PrCount                  int32     `db:"pr_count" json:"pr_count"`
	TotalPhotoCount          int32     `db:"total_photo_count" json:"total_photo_count"`
	WorkoutType              int32     `db:"workout_type" json:"workout_type"`
	SufferScore              int32     `db:"suffer_score" json:"suffer_score"`
	Calories                 float64   `db:"calories" json:"calories"`
	NumEfforts               int32     `db:"num_efforts" json:"num_efforts"`
	EmbedToken               string    `db:"embed_token" json:"embed_token"`
	SegmentLeaderboardOptOut bool      `db:"segment_leaderboard_opt_out" json:"segment_leaderboard_opt_out"`
	LeaderboardOptOut        bool      `db:"leaderboard_opt_out" json:"leaderboard_opt_out"`
	// Owner of the activity has premium account at the time of the fetch.
	PremiumFetch bool `db:"premium_fetch" json:"premium_fetch"`
}

type Athlete struct {
	ID          int64  `db:"id" json:"id"`
	Summit      bool   `db:"summit" json:"summit"`
	Username    string `db:"username" json:"username"`
	Firstname   string `db:"firstname" json:"firstname"`
	Lastname    string `db:"lastname" json:"lastname"`
	Sex         string `db:"sex" json:"sex"`
	City        string `db:"city" json:"city"`
	State       string `db:"state" json:"state"`
	Country     string `db:"country" json:"country"`
	FollowCount int32  `db:"follow_count" json:"follow_count"`
	FriendCount int32  `db:"friend_count" json:"friend_count"`
	// feet or meters
	MeasurementPreference string          `db:"measurement_preference" json:"measurement_preference"`
	Ftp                   float64         `db:"ftp" json:"ftp"`
	Weight                float64         `db:"weight" json:"weight"`
	Clubs                 json.RawMessage `db:"clubs" json:"clubs"`
	CreatedAt             time.Time       `db:"created_at" json:"created_at"`
	UpdatedAt             time.Time       `db:"updated_at" json:"updated_at"`
	FetchedAt             time.Time       `db:"fetched_at" json:"fetched_at"`
}

type AthleteLogin struct {
	AthleteID int64 `db:"athlete_id" json:"athlete_id"`
	Summit    bool  `db:"summit" json:"summit"`
	// Oauth app client ID
	ProviderID        string    `db:"provider_id" json:"provider_id"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time `db:"updated_at" json:"updated_at"`
	OauthAccessToken  string    `db:"oauth_access_token" json:"oauth_access_token"`
	OauthRefreshToken string    `db:"oauth_refresh_token" json:"oauth_refresh_token"`
	OauthExpiry       time.Time `db:"oauth_expiry" json:"oauth_expiry"`
	OauthTokenType    string    `db:"oauth_token_type" json:"oauth_token_type"`
	ID                uuid.UUID `db:"id" json:"id"`
}

type GueJob struct {
	JobID      string         `db:"job_id" json:"job_id"`
	Priority   int16          `db:"priority" json:"priority"`
	RunAt      time.Time      `db:"run_at" json:"run_at"`
	JobType    string         `db:"job_type" json:"job_type"`
	Args       []byte         `db:"args" json:"args"`
	ErrorCount int32          `db:"error_count" json:"error_count"`
	LastError  sql.NullString `db:"last_error" json:"last_error"`
	Queue      string         `db:"queue" json:"queue"`
	CreatedAt  time.Time      `db:"created_at" json:"created_at"`
	UpdatedAt  time.Time      `db:"updated_at" json:"updated_at"`
}

type Segment struct {
}

type Segment struct {
	ID   int32  `db:"id" json:"id"`
	Name string `db:"name" json:"name"`
}

type SegmentEffort struct {
	ID             int64     `db:"id" json:"id"`
	AthleteID      int64     `db:"athlete_id" json:"athlete_id"`
	SegmentID      int64     `db:"segment_id" json:"segment_id"`
	Name           string    `db:"name" json:"name"`
	ElapsedTime    float64   `db:"elapsed_time" json:"elapsed_time"`
	MovingTime     float64   `db:"moving_time" json:"moving_time"`
	StartDate      time.Time `db:"start_date" json:"start_date"`
	StartDateLocal time.Time `db:"start_date_local" json:"start_date_local"`
	// Distance is in meters
	Distance     float64       `db:"distance" json:"distance"`
	StartIndex   int32         `db:"start_index" json:"start_index"`
	EndIndex     int32         `db:"end_index" json:"end_index"`
	DeviceWatts  bool          `db:"device_watts" json:"device_watts"`
	AverageWatts float64       `db:"average_watts" json:"average_watts"`
	KomRank      sql.NullInt32 `db:"kom_rank" json:"kom_rank"`
	PrRank       sql.NullInt32 `db:"pr_rank" json:"pr_rank"`
}

type WebhookDump struct {
	ID         uuid.UUID `db:"id" json:"id"`
	RecordedAt time.Time `db:"recorded_at" json:"recorded_at"`
	Raw        string    `db:"raw" json:"raw"`
}