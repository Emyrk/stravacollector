package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/Emyrk/strava/database"
	"github.com/go-chi/chi/v5"

	"github.com/Emyrk/strava/api/httpapi"
	"github.com/Emyrk/strava/api/httpmw"
	"github.com/Emyrk/strava/api/modelsdk"
)

func (api *API) athleteHugels(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	row := httpmw.Athlete(r)
	athlete := row.Athlete

	activities, err := api.Opts.DB.AthleteHugelActivites(ctx, athlete.ID)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	acts := modelsdk.AthleteHugelActivities{
		Activities: convertHugelAthleteActivities(activities),
	}

	httpapi.Write(ctx, rw, http.StatusOK, acts)
}

func (api *API) athlete(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	row := httpmw.Athlete(r)
	athlete := row.Athlete

	httpapi.Write(ctx, rw, http.StatusOK, modelsdk.AthleteSummary{
		AthleteID:            modelsdk.StringInt(athlete.ID),
		Summit:               athlete.Summit,
		Username:             athlete.Username,
		Firstname:            athlete.Firstname,
		Lastname:             athlete.Lastname,
		Sex:                  athlete.Sex,
		ProfilePicLink:       athlete.ProfilePicLink,
		ProfilePicLinkMedium: athlete.ProfilePicLinkMedium,
		UpdatedAt:            athlete.UpdatedAt,
		HugelCount:           int(row.HugelCount),
	})
}

func (api *API) syncSummary(rw http.ResponseWriter, r *http.Request) {
	var (
		ctx  = r.Context()
		id   = httpmw.AuthenticatedAthleteID(r)
		page = r.URL.Query().Get("page")
	)

	pageNum, err := strconv.ParseInt(page, 10, 64)
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusBadRequest, modelsdk.Response{
			Message: "Invalid offset",
			Detail:  err.Error(),
		})
		return
	}

	detailedLoad, err := api.Opts.DB.GetAthleteLoadDetailed(ctx, id)
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, modelsdk.Response{
			Message: "Failed to fetch authenticated athlete sync information",
			Detail:  err.Error(),
		})
		return
	}
	load := detailedLoad.AthleteLoad

	limit := 100
	if pageNum <= 0 {
		pageNum = 1
	}
	offset := (pageNum - 1) * int64(limit)
	activities, err := api.Opts.DB.AthleteSyncedActivities(ctx, database.AthleteSyncedActivitiesParams{
		AthleteID: id,
		Offset:    int32(offset),
		Limit:     100,
	})
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, modelsdk.Response{
			Message: "Failed to fetch authenticated athlete's activities",
			Detail:  err.Error(),
		})
		return
	}
	sdk := make([]modelsdk.SyncActivitySummary, 0, len(activities))
	for _, act := range activities {
		sdk = append(sdk, modelsdk.SyncActivitySummary{
			ActivitySummary: convertActivitySummary(act.ActivitySummary),
			Synced:          act.DetailExists,
			SyncedAt:        act.DetailUpdatedAt.Time,
		})
	}

	httpapi.Write(ctx, rw, http.StatusOK, modelsdk.AthleteSyncSummary{
		TotalSummary: int(detailedLoad.SummaryCount),
		TotalDetail:  int(detailedLoad.DetailCount),
		Load: modelsdk.AthleteLoad{
			AthleteID:                  load.AthleteID,
			LastBackloadActivityStart:  load.LastBackloadActivityStart,
			LastLoadAttempt:            load.LastLoadAttempt,
			LastLoadIncomplete:         load.LastLoadIncomplete,
			LastLoadError:              load.LastLoadError,
			ActivitesLoadedLastAttempt: load.ActivitesLoadedLastAttempt,
			EarliestActivity:           load.EarliestActivity,
			EarliestActivityDone:       load.EarliestActivityDone,
		},
		SyncedActivities: sdk,
	})
}

func (api *API) whoAmI(rw http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
		id  = httpmw.AuthenticatedAthleteID(r)
	)

	full, err := api.Opts.DB.GetAthleteLoginFull(ctx, id)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	athlete := full.Athlete

	if errors.Is(err, sql.ErrNoRows) {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, modelsdk.Response{
			Message: "Failed to fetch authenticated athlete",
			Detail:  "Please try to log out and login again.",
		})
		return
	}

	httpapi.Write(ctx, rw, http.StatusOK, modelsdk.AthleteSummary{
		AthleteID:            modelsdk.StringInt(id),
		Summit:               full.AthleteLogin.Summit,
		Username:             athlete.Username,
		Firstname:            athlete.Firstname,
		Lastname:             athlete.Lastname,
		Sex:                  athlete.Sex,
		ProfilePicLink:       athlete.ProfilePicLink,
		ProfilePicLinkMedium: athlete.ProfilePicLinkMedium,
		UpdatedAt:            athlete.UpdatedAt,
		HugelCount:           int(full.HugelCount),
	})
}

func (api *API) manualFetchActivity(rw http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
		id  = httpmw.AuthenticatedAthleteID(r)
	)

	// Only steven can do this
	if id != 2661162 {
		httpapi.Write(ctx, rw, http.StatusUnauthorized, modelsdk.Response{
			Message: "Not authorized",
		})
		return
	}

	actID, err := strconv.ParseInt(chi.URLParam(r, "activity_id"), 10, 64)
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusBadRequest, modelsdk.Response{
			Message: "Invalid activity ID",
			Detail:  err.Error(),
		})
		return
	}

	err = api.Manager.EnqueueFetchActivity(ctx, database.ActivityDetailSourceManual, id, actID)
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusBadRequest, modelsdk.Response{
			Message: "Enqueue fetch",
			Detail:  err.Error(),
		})
		return
	}

	httpapi.Write(ctx, rw, http.StatusOK, modelsdk.Response{
		Message: "Enqueued",
	})
}

func convertActivitySummaries(activities []database.ActivitySummary) []modelsdk.ActivitySummary {
	sdk := make([]modelsdk.ActivitySummary, 0, len(activities))
	for _, act := range activities {
		sdk = append(sdk, convertActivitySummary(act))
	}
	return sdk
}

func convertActivitySummary(activity database.ActivitySummary) modelsdk.ActivitySummary {
	return modelsdk.ActivitySummary{
		ActivityID:     modelsdk.StringInt(activity.ID),
		AthleteID:      modelsdk.StringInt(activity.AthleteID),
		UploadID:       modelsdk.StringInt(activity.UploadID),
		ExternalID:     activity.ExternalID,
		Name:           activity.Name,
		Distance:       activity.Distance,
		MovingTime:     activity.MovingTime,
		ElapsedTime:    activity.ElapsedTime,
		TotalEleGain:   activity.TotalElevationGain,
		ActivityType:   activity.ActivityType,
		SportType:      activity.SportType,
		StartDate:      activity.StartDate,
		StartDateLocal: activity.StartDateLocal,
		Timezone:       activity.Timezone,
	}
}

func convertHugelAthleteActivities(activities []database.AthleteHugelActivitesRow) []modelsdk.HugelLeaderBoardActivity {
	sdk := make([]modelsdk.HugelLeaderBoardActivity, 0, len(activities))
	for _, act := range activities {
		sdk = append(sdk, convertHugelAthleteActivity(act))
	}
	return sdk
}

func convertHugelAthleteActivity(activity database.AthleteHugelActivitesRow) modelsdk.HugelLeaderBoardActivity {
	var efforts []modelsdk.SegmentEffort
	_ = json.Unmarshal(activity.Efforts, &efforts)
	return modelsdk.HugelLeaderBoardActivity{
		RankOneElapsed: activity.BestTime,
		ActivityID:     modelsdk.StringInt(activity.ActivityID),
		AthleteID:      modelsdk.StringInt(activity.AthleteID),
		Elapsed:        activity.TotalTimeSeconds,
		Efforts:        efforts,
		Athlete: modelsdk.MinAthlete{
			AthleteID:      modelsdk.StringInt(activity.AthleteID),
			Username:       activity.Username,
			Firstname:      activity.Firstname,
			Lastname:       activity.Lastname,
			Sex:            activity.Sex,
			ProfilePicLink: activity.ProfilePicLink,
			HugelCount:     int(activity.HugelCount),
		},
		ActivityName:               activity.Name,
		ActivityDistance:           activity.Distance,
		ActivityMovingTime:         int64(activity.MovingTime),
		ActivityElapsedTime:        int64(activity.ElapsedTime),
		ActivityStartDate:          activity.StartDate,
		ActivityTotalElevationGain: activity.TotalElevationGain,
	}
}
