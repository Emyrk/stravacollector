package modelsdk

import "time"

type VerifyRouteResponse struct {
	MissingSegments []SegmentSummary `json:"missing_segments"`
}

type CompetitiveRoute struct {
	Name        string           `json:"name"`
	DisplayName string           `json:"display_name"`
	Description string           `json:"description"`
	Segments    []SegmentSummary `json:"segments"`
}

type SegmentSummary struct {
	ID   StringInt `json:"id"`
	Name string    `json:"name"`
}

type PersonalSegment struct {
	DetailedSegment DetailedSegment `json:"detailed_segment"`

	Starred      bool                       `json:"starred,omitempty"`
	PersonalBest *PersonalBestSegmentEffort `json:"personal_best,omitempty"`
}

type PersonalBestSegmentEffort struct {
	BestEffortID             StringInt `json:"best_effort_id"`
	BestEffortElapsedTime    float64   `json:"best_effort_elapsed_time"`
	BestEffortMovingTime     float64   `json:"best_effort_moving_time"`
	BestEffortStartDate      time.Time `json:"best_effort_start_date"`
	BestEffortStartDateLocal time.Time `json:"best_effort_start_date_local"`
	BestEffortDeviceWatts    bool      `json:"best_effort_device_watts"`
	BestEffortAverageWatts   float64   `json:"best_effort_average_watts"`
	BestEffortActivitiesID   StringInt `json:"best_effort_activities_id"`
}

type DetailedSegment struct {
	ID            StringInt `json:"id"`
	Name          string    `json:"name"`
	FriendlyName  string    `json:"friendly_name"`
	ActivityType  string    `json:"activity_type"`
	Distance      float64   `json:"distance"`
	AverageGrade  float64   `json:"average_grade"`
	MaximumGrade  float64   `json:"maximum_grade"`
	ElevationHigh float64   `json:"elevation_high"`
	ElevationLow  float64   `json:"elevation_low"`
	StartLatlng   []float64 `json:"start_latlng"`
	EndLatlng     []float64 `json:"end_latlng"`
	// A small image of the elevation profile of this segment.
	ElevationProfile   string    `json:"elevation_profile"`
	ClimbCategory      int32     `json:"climb_category"`
	City               string    `json:"city"`
	State              string    `json:"state"`
	Country            string    `json:"country"`
	Private            bool      `json:"private"`
	Hazardous          bool      `json:"hazardous"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	TotalElevationGain float64   `json:"total_elevation_gain"`
	Map                Map       `json:"map"`
	TotalEffortCount   int32     `json:"total_effort_count"`
	TotalAthleteCount  int32     `json:"total_athlete_count"`
	TotalStarCount     int32     `json:"total_star_count"`
	// The time at which this segment was fetched from the Strava API.
	FetchedAt time.Time `json:"fetched_at"`
}
