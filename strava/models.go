package strava

import "time"

type SegmentEffort struct {
	ID            int    `json:"id"`
	ResourceState int    `json:"resource_state"`
	Name          string `json:"name"`
	Activity      struct {
		ID            int `json:"id"`
		ResourceState int `json:"resource_state"`
	} `json:"activity"`
	Athlete struct {
		ID int64 `json:"id"`
		// 1 == Meta
		// 2 == Summary
		// 3 == Detail
		ResourceState int `json:"resource_state"`
	} `json:"athlete"`
	ElapsedTime    int            `json:"elapsed_time"`
	MovingTime     int            `json:"moving_time"`
	StartDate      time.Time      `json:"start_date"`
	StartDateLocal time.Time      `json:"start_date_local"`
	Distance       float64        `json:"distance"`
	StartIndex     int            `json:"start_index"`
	EndIndex       int            `json:"end_index"`
	DeviceWatts    bool           `json:"device_watts"`
	AverageWatts   float64        `json:"average_watts"`
	Segment        SegmentSummary `json:"segment"`
	KomRank        *int           `json:"kom_rank"`
	PrRank         *int           `json:"pr_rank"`
	Achievements   []interface{}  `json:"achievements"`
}

type SegmentSummary struct {
	ID            int       `json:"id"`
	ResourceState int       `json:"resource_state"`
	Name          string    `json:"name"`
	ActivityType  string    `json:"activity_type"`
	Distance      float64   `json:"distance"`
	AverageGrade  float64   `json:"average_grade"`
	MaximumGrade  float64   `json:"maximum_grade"`
	ElevationHigh float64   `json:"elevation_high"`
	ElevationLow  float64   `json:"elevation_low"`
	StartLatlng   []float64 `json:"start_latlng"`
	EndLatlng     []float64 `json:"end_latlng"`
	ClimbCategory int       `json:"climb_category"`
	City          string    `json:"city"`
	State         string    `json:"state"`
	Country       string    `json:"country"`
	Private       bool      `json:"private"`
	Hazardous     bool      `json:"hazardous"`
	Starred       bool      `json:"starred"`
}

type Athlete struct {
	ID                    int64         `json:"id"`
	Username              string        `json:"username"`
	ResourceState         int           `json:"resource_state"`
	Firstname             string        `json:"firstname"`
	Lastname              string        `json:"lastname"`
	City                  string        `json:"city"`
	State                 string        `json:"state"`
	Country               string        `json:"country"`
	Sex                   string        `json:"sex"`
	Premium               bool          `json:"premium"`
	CreatedAt             time.Time     `json:"created_at"`
	UpdatedAt             time.Time     `json:"updated_at"`
	BadgeTypeID           int           `json:"badge_type_id"`
	ProfileMedium         string        `json:"profile_medium"`
	Profile               string        `json:"profile"`
	Friend                interface{}   `json:"friend"`
	Follower              interface{}   `json:"follower"`
	FollowerCount         int           `json:"follower_count"`
	FriendCount           int           `json:"friend_count"`
	MutualFriendCount     int           `json:"mutual_friend_count"`
	AthleteType           int           `json:"athlete_type"`
	DatePreference        string        `json:"date_preference"`
	MeasurementPreference string        `json:"measurement_preference"`
	Clubs                 []interface{} `json:"clubs"`
	Ftp                   interface{}   `json:"ftp"`
	Weight                int           `json:"weight"`
	Bikes                 []Equipment   `json:"bikes"`
	Shoes                 []Equipment   `json:"shoes"`
}

type Equipment struct {
	ID            string `json:"id"`
	Primary       bool   `json:"primary"`
	Name          string `json:"name"`
	ResourceState int    `json:"resource_state"`
	Distance      int    `json:"distance"`
}
