package api

import (
	"net/http"

	"github.com/Emyrk/strava/api/httpmw"

	"github.com/Emyrk/strava/api/httpapi"
	"github.com/Emyrk/strava/api/modelsdk"
	"github.com/Emyrk/strava/database"
)

func (api *API) getSegments(rw http.ResponseWriter, r *http.Request) {
	var (
		id, athleteLoggedIn = httpmw.AuthenticatedAthleteIDOptional(r)
		ctx                 = r.Context()
	)

	var requestedSegments []modelsdk.StringInt
	if !httpapi.Read(ctx, rw, r, &requestedSegments) {
		return
	}

	requestedSegmentsInts := make([]int64, len(requestedSegments))
	for i, seg := range requestedSegments {
		requestedSegmentsInts[i] = int64(seg)
	}

	var resp []modelsdk.PersonalSegment
	if athleteLoggedIn {
		segments, err := api.Opts.DB.GetPersonalSegments(ctx, database.GetPersonalSegmentsParams{
			AthleteID:  id,
			SegmentIds: requestedSegmentsInts,
		})
		if err != nil {
			httpapi.Write(ctx, rw, http.StatusInternalServerError, modelsdk.Response{
				Message: "Failed to load segments",
				Detail:  err.Error(),
			})
			return
		}
		resp = convertSegmentRows(segments)
	} else {
		segments, err := api.Opts.DB.GetSegments(ctx, requestedSegmentsInts)
		if err != nil {
			httpapi.Write(ctx, rw, http.StatusInternalServerError, modelsdk.Response{
				Message: "Failed to load segments",
				Detail:  err.Error(),
			})
			return
		}
		resp = convertSegmentRows(segments)
	}

	httpapi.Write(ctx, rw, http.StatusOK, resp)
}

type segmentRow interface {
	database.GetPersonalSegmentsRow | database.GetSegmentsRow
}

func convertSegmentRows[S segmentRow](rows []S) []modelsdk.PersonalSegment {
	segments := make([]modelsdk.PersonalSegment, len(rows))
	for i, row := range rows {
		segments[i] = convertSegmentRow(row)
	}
	return segments
}

func convertSegmentRow[S segmentRow](row S) modelsdk.PersonalSegment {
	var segment database.Segment
	var dbMap database.Map
	var starred bool
	var best *modelsdk.PersonalBestSegmentEffort
	switch row := any(row).(type) {
	case database.GetPersonalSegmentsRow:
		segment = row.Segment
		dbMap = row.Map
		starred = row.Starred
		if row.BestEffortID > 0 {
			best = &modelsdk.PersonalBestSegmentEffort{
				BestEffortID:             modelsdk.StringInt(row.BestEffortID),
				BestEffortElapsedTime:    row.BestEffortElapsedTime,
				BestEffortMovingTime:     row.BestEffortMovingTime,
				BestEffortStartDate:      row.BestEffortStartDate,
				BestEffortStartDateLocal: row.BestEffortStartDateLocal,
				BestEffortDeviceWatts:    row.BestEffortDeviceWatts,
				BestEffortAverageWatts:   row.BestEffortAverageWatts,
				BestEffortActivitiesID:   modelsdk.StringInt(row.BestEffortActivitiesID),
			}
		}
	case database.GetSegmentsRow:
		segment = row.Segment
		dbMap = row.Map
	}

	return modelsdk.PersonalSegment{
		DetailedSegment: modelsdk.DetailedSegment{
			ID:                 modelsdk.StringInt(segment.ID),
			Name:               segment.Name,
			ActivityType:       segment.ActivityType,
			Distance:           segment.Distance,
			AverageGrade:       segment.AverageGrade,
			MaximumGrade:       segment.MaximumGrade,
			ElevationHigh:      segment.ElevationHigh,
			ElevationLow:       segment.ElevationLow,
			StartLatlng:        segment.StartLatlng,
			EndLatlng:          segment.EndLatlng,
			ElevationProfile:   segment.ElevationProfile,
			ClimbCategory:      segment.ClimbCategory,
			City:               segment.City,
			State:              segment.State,
			Country:            segment.Country,
			Private:            segment.Private,
			Hazardous:          segment.Hazardous,
			CreatedAt:          segment.CreatedAt,
			UpdatedAt:          segment.UpdatedAt,
			TotalElevationGain: segment.TotalElevationGain,
			Map:                convertMap(dbMap),
			TotalEffortCount:   segment.TotalEffortCount,
			TotalAthleteCount:  segment.TotalAthleteCount,
			TotalStarCount:     segment.TotalStarCount,
			FetchedAt:          segment.FetchedAt,
		},
		Starred:      starred,
		PersonalBest: best,
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
