package api

import (
	"net/http"

	"github.com/Emyrk/strava/api/httpapi"
	"github.com/Emyrk/strava/api/modelsdk"
	"github.com/Emyrk/strava/database"
)

func (api *API) getSegments(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var requestedSegments []modelsdk.StringInt
	if !httpapi.Read(ctx, rw, r, &requestedSegments) {
		return
	}

	requestedSegmentsInts := make([]int64, len(requestedSegments))
	for i, seg := range requestedSegments {
		requestedSegmentsInts[i] = int64(seg)
	}

	segments, err := api.Opts.DB.GetSegments(ctx, requestedSegmentsInts)
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, modelsdk.Response{
			Message: "Failed to load segments",
			Detail:  err.Error(),
		})
		return
	}

	httpapi.Write(ctx, rw, http.StatusOK, convertSegmentRows(segments))
}

func convertSegmentRows(rows []database.GetSegmentsRow) []modelsdk.DetailedSegment {
	segments := make([]modelsdk.DetailedSegment, len(rows))
	for i, row := range rows {
		segments[i] = convertSegmentRow(row)
	}
	return segments
}

func convertSegmentRow(row database.GetSegmentsRow) modelsdk.DetailedSegment {
	return modelsdk.DetailedSegment{
		ID:                 modelsdk.StringInt(row.Segment.ID),
		Name:               row.Segment.Name,
		ActivityType:       row.Segment.ActivityType,
		Distance:           row.Segment.Distance,
		AverageGrade:       row.Segment.AverageGrade,
		MaximumGrade:       row.Segment.MaximumGrade,
		ElevationHigh:      row.Segment.ElevationHigh,
		ElevationLow:       row.Segment.ElevationLow,
		StartLatlng:        row.Segment.StartLatlng,
		EndLatlng:          row.Segment.EndLatlng,
		ElevationProfile:   row.Segment.ElevationProfile,
		ClimbCategory:      row.Segment.ClimbCategory,
		City:               row.Segment.City,
		State:              row.Segment.State,
		Country:            row.Segment.Country,
		Private:            row.Segment.Private,
		Hazardous:          row.Segment.Hazardous,
		CreatedAt:          row.Segment.CreatedAt,
		UpdatedAt:          row.Segment.UpdatedAt,
		TotalElevationGain: row.Segment.TotalElevationGain,
		MapID:              convertMap(row.Map),
		TotalEffortCount:   row.Segment.TotalEffortCount,
		TotalAthleteCount:  row.Segment.TotalAthleteCount,
		TotalStarCount:     row.Segment.TotalStarCount,
		FetchedAt:          row.Segment.FetchedAt,
	}
}

func convertMap(m database.Map) modelsdk.Map {
	return modelsdk.Map{
		ID:              m.ID,
		Polyline:        m.Polyline,
		SummaryPolyline: m.SummaryPolyline,
		UpdatedAt:       m.UpdatedAt,
	}
}
