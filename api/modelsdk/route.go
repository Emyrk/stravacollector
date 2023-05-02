package modelsdk

type CompetitiveRoute struct {
	Name        string           `json:"name"`
	DisplayName string           `json:"display_name"`
	Description string           `json:"description"`
	Segments    []SegmentSummary `json:"segments"`
}

type SegmentSummary struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}
