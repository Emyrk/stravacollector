package strava

import (
	"encoding/json"
	"time"
)

type Map struct {
	ID              string `json:"id"`
	Polyline        string `json:"polyline"`
	ResourceState   int    `json:"resource_state"`
	SummaryPolyline string `json:"summary_polyline"`
}

type ActivitySummary struct {
	ID         int64  `json:"id"`
	ExternalID string `json:"external_id"`
	UploadID   int64  `json:"upload_id"`
	Athlete    struct {
		ID            int64 `json:"id"`
		ResourceState int   `json:"resource_state"`
	} `json:"athlete"`
	Name               string    `json:"name"`
	Distance           float64   `json:"distance"`
	MovingTime         float64   `json:"moving_time"`
	ElapsedTime        float64   `json:"elapsed_time"`
	TotalElevationGain float64   `json:"total_elevation_gain"`
	Type               string    `json:"type"`
	SportType          string    `json:"sport_type"`
	StartDate          time.Time `json:"start_date"`
	StartDateLocal     time.Time `json:"start_date_local"`
	Timezone           string    `json:"timezone"`
	UtcOffset          float64   `json:"utc_offset"`
	WorkoutType        int32     `json:"workout_type"`
	LocationCity       string    `json:"location_city"`
	LocationState      string    `json:"location_state"`
	LocationCountry    string    `json:"location_country"`
	AchievementCount   int32     `json:"achievement_count"`
	KudosCount         int32     `json:"kudos_count"`
	CommentCount       int32     `json:"comment_count"`
	AthleteCount       int32     `json:"athlete_count"`
	PhotoCount         int32     `json:"photo_count"`
	Map                struct {
		ID              string `json:"id"`
		SummaryPolyline string `json:"summary_polyline"`
	} `json:"map"`
	Trainer              bool    `json:"trainer"`
	Commute              bool    `json:"commute"`
	Manual               bool    `json:"manual"`
	Private              bool    `json:"private"`
	Flagged              bool    `json:"flagged"`
	GearID               string  `json:"gear_id"`
	FromAcceptedTag      bool    `json:"from_accepted_tag"`
	AverageSpeed         float64 `json:"average_speed"`
	MaxSpeed             float64 `json:"max_speed"`
	AverageCadence       float64 `json:"average_cadence"`
	AverageWatts         float64 `json:"average_watts"`
	WeightedAverageWatts float64 `json:"weighted_average_watts"`
	Kilojoules           float64 `json:"kilojoules"`
	DeviceWatts          bool    `json:"device_watts"`
	HasHeartrate         bool    `json:"has_heartrate"`
	AverageHeartrate     float64 `json:"average_heartrate"`
	MaxHeartrate         float64 `json:"max_heartrate"`
	MaxWatts             float64 `json:"max_watts"`
	PrCount              int32   `json:"pr_count"`
	TotalPhotoCount      int32   `json:"total_photo_count"`
	HasKudoed            bool    `json:"has_kudoed"`
	SufferScore          float64 `json:"suffer_score"`
}

type DetailedActivity struct {
	ID            int64  `json:"id"`
	ResourceState int    `json:"resource_state"`
	ExternalID    string `json:"external_id"`
	UploadID      int64  `json:"upload_id"`
	Athlete       struct {
		ID            int64 `json:"id"`
		ResourceState int   `json:"resource_state"`
	} `json:"athlete"`
	Name                 string                  `json:"name"`
	Distance             float64                 `json:"distance"`
	MovingTime           float64                 `json:"moving_time"`
	ElapsedTime          float64                 `json:"elapsed_time"`
	TotalElevationGain   float64                 `json:"total_elevation_gain"`
	Type                 string                  `json:"type"`
	SportType            string                  `json:"sport_type"`
	StartDate            time.Time               `json:"start_date"`
	StartDateLocal       time.Time               `json:"start_date_local"`
	Timezone             string                  `json:"timezone"`
	UtcOffset            float64                 `json:"utc_offset"`
	StartLatlng          []float64               `json:"start_latlng"`
	EndLatlng            []float64               `json:"end_latlng"`
	AchievementCount     int32                   `json:"achievement_count"`
	KudosCount           int32                   `json:"kudos_count"`
	CommentCount         int32                   `json:"comment_count"`
	AthleteCount         int32                   `json:"athlete_count"`
	PhotoCount           int32                   `json:"photo_count"`
	Map                  Map                     `json:"map"`
	Trainer              bool                    `json:"trainer"`
	Commute              bool                    `json:"commute"`
	Manual               bool                    `json:"manual"`
	Private              bool                    `json:"private"`
	Flagged              bool                    `json:"flagged"`
	GearID               string                  `json:"gear_id"`
	FromAcceptedTag      bool                    `json:"from_accepted_tag"`
	AverageSpeed         float64                 `json:"average_speed"`
	MaxSpeed             float64                 `json:"max_speed"`
	AverageCadence       float64                 `json:"average_cadence"`
	AverageTemp          float64                 `json:"average_temp"`
	AverageWatts         float64                 `json:"average_watts"`
	WeightedAverageWatts float64                 `json:"weighted_average_watts"`
	Kilojoules           float64                 `json:"kilojoules"`
	DeviceWatts          bool                    `json:"device_watts"`
	HasHeartrate         bool                    `json:"has_heartrate"`
	AverageHeartrate     float64                 `json:"average_heartrate"`
	MaxHeartrate         float64                 `json:"max_heartrate"`
	MaxWatts             float64                 `json:"max_watts"`
	ElevHigh             float64                 `json:"elev_high"`
	ElevLow              float64                 `json:"elev_low"`
	PrCount              int32                   `json:"pr_count"`
	TotalPhotoCount      int32                   `json:"total_photo_count"`
	HasKudoed            bool                    `json:"has_kudoed"`
	WorkoutType          int32                   `json:"workout_type"`
	SufferScore          float64                 `json:"suffer_score"`
	Description          string                  `json:"description"`
	Calories             float64                 `json:"calories"`
	SegmentEfforts       []DetailedSegmentEffort `json:"segment_efforts"`
	SplitsMetric         []struct {
		Distance            float64 `json:"distance"`
		ElapsedTime         int     `json:"elapsed_time"`
		ElevationDifference float64 `json:"elevation_difference"`
		MovingTime          int     `json:"moving_time"`
		Split               int     `json:"split"`
		AverageSpeed        float64 `json:"average_speed"`
		PaceZone            int     `json:"pace_zone"`
	} `json:"splits_metric"`
	Laps []struct {
		ID            int64  `json:"id"`
		ResourceState int    `json:"resource_state"`
		Name          string `json:"name"`
		Activity      struct {
			ID            int `json:"id"`
			ResourceState int `json:"resource_state"`
		} `json:"activity"`
		Athlete struct {
			ID            int `json:"id"`
			ResourceState int `json:"resource_state"`
		} `json:"athlete"`
		ElapsedTime        float64   `json:"elapsed_time"`
		MovingTime         float64   `json:"moving_time"`
		StartDate          time.Time `json:"start_date"`
		StartDateLocal     time.Time `json:"start_date_local"`
		Distance           float64   `json:"distance"`
		StartIndex         float64   `json:"start_index"`
		EndIndex           int       `json:"end_index"`
		TotalElevationGain float64   `json:"total_elevation_gain"`
		AverageSpeed       float64   `json:"average_speed"`
		MaxSpeed           float64   `json:"max_speed"`
		AverageCadence     float64   `json:"average_cadence"`
		DeviceWatts        bool      `json:"device_watts"`
		AverageWatts       float64   `json:"average_watts"`
		LapIndex           int       `json:"lap_index"`
		Split              int       `json:"split"`
	} `json:"laps"`
	Gear struct {
		ID            string  `json:"id"`
		Primary       bool    `json:"primary"`
		Name          string  `json:"name"`
		ResourceState int     `json:"resource_state"`
		Distance      float64 `json:"distance"`
	} `json:"gear"`
	PartnerBrandTag interface{} `json:"partner_brand_tag"`
	Photos          struct {
		Primary struct {
			ID       interface{} `json:"id"`
			UniqueID string      `json:"unique_id"`
			Urls     struct {
				Num100 string `json:"100"`
				Num600 string `json:"600"`
			} `json:"urls"`
			Source int `json:"source"`
		} `json:"primary"`
		UsePrimaryPhoto bool `json:"use_primary_photo"`
		Count           int  `json:"count"`
	} `json:"photos"`
	HighlightedKudosers []struct {
		DestinationURL string `json:"destination_url"`
		DisplayName    string `json:"display_name"`
		AvatarURL      string `json:"avatar_url"`
		ShowName       bool   `json:"show_name"`
	} `json:"highlighted_kudosers"`
	HideFromHome             bool   `json:"hide_from_home"`
	DeviceName               string `json:"device_name"`
	EmbedToken               string `json:"embed_token"`
	SegmentLeaderboardOptOut bool   `json:"segment_leaderboard_opt_out"`
	LeaderboardOptOut        bool   `json:"leaderboard_opt_out"`
}

type DetailedSegmentEffort struct {
	ID            int64  `json:"id"`
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
	ElapsedTime    float64        `json:"elapsed_time"`
	MovingTime     float64        `json:"moving_time"`
	StartDate      time.Time      `json:"start_date"`
	StartDateLocal time.Time      `json:"start_date_local"`
	Distance       float64        `json:"distance"`
	StartIndex     int32          `json:"start_index"`
	EndIndex       int32          `json:"end_index"`
	DeviceWatts    bool           `json:"device_watts"`
	AverageWatts   float64        `json:"average_watts"`
	Segment        SegmentSummary `json:"segment"`
	KomRank        int32          `json:"kom_rank"`
	PrRank         int32          `json:"pr_rank"`
	Achievements   []interface{}  `json:"achievements"`
}

type SegmentSummary struct {
	ID            int64     `json:"id"`
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
	ID                    int64           `json:"id"`
	Username              string          `json:"username"`
	ResourceState         int             `json:"resource_state"`
	Firstname             string          `json:"firstname"`
	Lastname              string          `json:"lastname"`
	City                  string          `json:"city"`
	State                 string          `json:"state"`
	Country               string          `json:"country"`
	Sex                   string          `json:"sex"`
	Premium               bool            `json:"premium"`
	Summit                bool            `json:"summit"`
	CreatedAt             time.Time       `json:"created_at"`
	UpdatedAt             time.Time       `json:"updated_at"`
	BadgeTypeID           int             `json:"badge_type_id"`
	ProfileMedium         string          `json:"profile_medium"`
	Profile               string          `json:"profile"`
	Friend                interface{}     `json:"friend"`
	Follower              interface{}     `json:"follower"`
	FollowerCount         int             `json:"follower_count"`
	FriendCount           int             `json:"friend_count"`
	MutualFriendCount     int             `json:"mutual_friend_count"`
	AthleteType           int             `json:"athlete_type"`
	DatePreference        string          `json:"date_preference"`
	MeasurementPreference string          `json:"measurement_preference"`
	Clubs                 json.RawMessage `json:"clubs"`
	Ftp                   float64         `json:"ftp"`
	Weight                float64         `json:"weight"`
	Bikes                 []Equipment     `json:"bikes"`
	Shoes                 []Equipment     `json:"shoes"`
}

type AthleteSummary struct {
	ID            int64       `json:"id"`
	Username      string      `json:"username"`
	ResourceState int         `json:"resource_state"`
	Firstname     string      `json:"firstname"`
	Lastname      string      `json:"lastname"`
	Bio           string      `json:"bio"`
	City          string      `json:"city"`
	State         string      `json:"state"`
	Country       string      `json:"country"`
	Sex           string      `json:"sex"`
	Premium       bool        `json:"premium"`
	Summit        bool        `json:"summit"`
	CreatedAt     time.Time   `json:"created_at"`
	UpdatedAt     time.Time   `json:"updated_at"`
	BadgeTypeID   int         `json:"badge_type_id"`
	ProfileMedium string      `json:"profile_medium"`
	Profile       string      `json:"profile"`
	Friend        interface{} `json:"friend"`
	Follower      interface{} `json:"follower"`
}

type Equipment struct {
	ID            string `json:"id"`
	Primary       bool   `json:"primary"`
	Name          string `json:"name"`
	ResourceState int    `json:"resource_state"`
	Distance      int    `json:"distance"`
}

type SegmentDetailed struct {
	ID                  int64     `json:"id"`
	ResourceState       int       `json:"resource_state"`
	Name                string    `json:"name"`
	ActivityType        string    `json:"activity_type"`
	Distance            float64   `json:"distance"`
	AverageGrade        float64   `json:"average_grade"`
	MaximumGrade        float64   `json:"maximum_grade"`
	ElevationHigh       float64   `json:"elevation_high"`
	ElevationLow        float64   `json:"elevation_low"`
	StartLatlng         []float64 `json:"start_latlng"`
	EndLatlng           []float64 `json:"end_latlng"`
	ElevationProfile    string    `json:"elevation_profile"`
	ClimbCategory       int32     `json:"climb_category"`
	City                string    `json:"city"`
	State               string    `json:"state"`
	Country             string    `json:"country"`
	Private             bool      `json:"private"`
	Hazardous           bool      `json:"hazardous"`
	Starred             bool      `json:"starred"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
	TotalElevationGain  float64   `json:"total_elevation_gain"`
	Map                 Map       `json:"map"`
	EffortCount         int32     `json:"effort_count"`
	AthleteCount        int32     `json:"athlete_count"`
	StarCount           int32     `json:"star_count"`
	AthleteSegmentStats struct {
		PrElapsedTime interface{} `json:"pr_elapsed_time"`
		PrDate        interface{} `json:"pr_date"`
		PrActivityID  interface{} `json:"pr_activity_id"`
		EffortCount   int         `json:"effort_count"`
	} `json:"athlete_segment_stats"`
	Xoms struct {
		Kom         string `json:"kom"`
		Qom         string `json:"qom"`
		Overall     string `json:"overall"`
		Destination struct {
			Href string `json:"href"`
			Type string `json:"type"`
			Name string `json:"name"`
		} `json:"destination"`
	} `json:"xoms"`
	LocalLegend struct {
		AthleteID         int    `json:"athlete_id"`
		Title             string `json:"title"`
		Profile           string `json:"profile"`
		EffortDescription string `json:"effort_description"`
		EffortCount       string `json:"effort_count"`
		EffortCounts      struct {
			Overall string `json:"overall"`
			Female  string `json:"female"`
		} `json:"effort_counts"`
		Destination string `json:"destination"`
	} `json:"local_legend"`
}

type Route struct {
	Athlete       AthleteSummary `json:"athlete"`
	Description   string         `json:"description"`
	Distance      float64        `json:"distance"`
	ElevationGain float64        `json:"elevation_gain"`
	ID            int64          `json:"id"`
	IDStr         string         `json:"id_str"`
	Map           Map            `json:"map"`
	MapUrls       struct {
		URL       string `json:"url"`
		RetinaURL string `json:"retina_url"`
	} `json:"map_urls"`
	Name          string    `json:"name"`
	Private       bool      `json:"private"`
	ResourceState int       `json:"resource_state"`
	Starred       bool      `json:"starred"`
	SubType       int       `json:"sub_type"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	// When it was created
	Timestamp           int              `json:"timestamp"`
	Type                int              `json:"type"`
	EstimatedMovingTime int              `json:"estimated_moving_time"`
	Segments            []SegmentSummary `json:"segments"`
}

type Club struct {
	ID            int    `json:"id"`
	ResourceState int    `json:"resource_state"`
	Name          string `json:"name"`
	ProfileMedium string `json:"profile_medium"`
	Profile       string `json:"profile"`
	// These can be null
	CoverPhoto      *string `json:"cover_photo"`
	CoverPhotoSmall *string `json:"cover_photo_small"`
	// ActivityTypes Examples
	//		"Handcycle",
	//		"EBikeRide",
	//		"VirtualRide",
	//		"Velomobile",
	//		"Ride"
	ActivityTypes     []string `json:"activity_types"`
	ActivityTypesIcon string   `json:"activity_types_icon"`
	// Dimension example
	//		"distance",
	//		"num_activities",
	//		"best_activities_distance",
	//		"velocity",
	//		"elev_gain",
	//		"moving_time"
	Dimensions         []string `json:"dimensions"`
	SportType          string   `json:"sport_type"`
	LocalizedSportType string   `json:"localized_sport_type"`
	City               string   `json:"city"`
	State              string   `json:"state"`
	Country            string   `json:"country"`
	Private            bool     `json:"private"`
	MemberCount        int      `json:"member_count"`
	Featured           bool     `json:"featured"`
	Verified           bool     `json:"verified"`
	URL                string   `json:"url"`
	Membership         string   `json:"membership"`
	Admin              bool     `json:"admin"`
	Owner              bool     `json:"owner"`
}
