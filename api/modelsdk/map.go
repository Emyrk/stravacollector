package modelsdk

import "time"

type Map struct {
	ID              string    `json:"id"`
	Polyline        string    `json:"polyline"`
	SummaryPolyline string    `json:"summary_polyline"`
	UpdatedAt       time.Time `json:"updated_at"`
}
