package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Emyrk/strava/api/httpapi"
	"github.com/Emyrk/strava/api/httpmw"
	"github.com/Emyrk/strava/api/modelsdk"
	river2 "github.com/Emyrk/strava/api/river"
	"github.com/Emyrk/strava/database"
	"github.com/go-chi/chi/v5"
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
		UpdatedAt:            athlete.UpdatedAt.Time,
		HugelCount:           int(row.HugelCount),
	})
}

func (api *API) syncSummary(rw http.ResponseWriter, r *http.Request) {
	var (
		ctx         = r.Context()
		authAth, ok = httpmw.AuthenticatedAthleteIDOptional(r)
		ath         = httpmw.Athlete(r)
		page        = r.URL.Query().Get("page")
		limitStr    = r.URL.Query().Get("limit")
		err         error
	)
	if !ok {
		httpapi.Write(ctx, rw, http.StatusUnauthorized, modelsdk.Response{
			Message: "Synced data requires authentication. No authentication provided",
		})
		return
	}

	// 2661162 is Steven
	if authAth != 2661162 && authAth != ath.Athlete.ID {
		httpapi.Write(ctx, rw, http.StatusUnauthorized, modelsdk.Response{
			Message: "You can only fetch your own sync summary, not another athlete's.",
		})
		return
	}

	limit := int64(100)
	if limitStr != "" {
		limit, err = strconv.ParseInt(limitStr, 10, 64)
		if err != nil {
			httpapi.Write(ctx, rw, http.StatusBadRequest, modelsdk.Response{
				Message: "Invalid limit",
				Detail:  err.Error(),
			})
			return
		}
	}

	pageNum := int64(1)
	if page != "" {
		pageNum, err = strconv.ParseInt(page, 10, 64)
		if err != nil {
			httpapi.Write(ctx, rw, http.StatusBadRequest, modelsdk.Response{
				Message: "Invalid offset",
				Detail:  err.Error(),
			})
			return
		}
	}

	detailedLoad, err := api.Opts.DB.GetAthleteLoadDetailed(ctx, ath.Athlete.ID)
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, modelsdk.Response{
			Message: "Failed to fetch authenticated athlete sync information",
			Detail:  err.Error(),
		})
		return
	}
	load := detailedLoad.AthleteLoad

	if pageNum <= 0 {
		pageNum = 1
	}
	offset := (pageNum - 1) * limit
	total := 0
	activities, err := api.Opts.DB.AthleteSyncedActivities(ctx, database.AthleteSyncedActivitiesParams{
		AthleteID: ath.Athlete.ID,
		Offset:    int32(offset),
		Limit:     int32(limit),
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
			Activity: convertActivitySummary(act.ActivitySummary),
			Synced:   act.DetailExists,
			SyncedAt: act.DetailUpdatedAt.Time,
		})
		total = int(act.Total)
	}

	athlete := detailedLoad.Athlete
	httpapi.Write(ctx, rw, http.StatusOK, modelsdk.AthleteSyncSummary{
		TotalSummary: int(detailedLoad.SummaryCount),
		TotalDetail:  int(detailedLoad.DetailCount),
		Athlete: modelsdk.AthleteSummary{
			AthleteID:            modelsdk.StringInt(athlete.ID),
			Summit:               athlete.Summit,
			Username:             athlete.Username,
			Firstname:            athlete.Firstname,
			Lastname:             athlete.Lastname,
			Sex:                  athlete.Sex,
			ProfilePicLink:       athlete.ProfilePicLink,
			ProfilePicLinkMedium: athlete.ProfilePicLinkMedium,
			UpdatedAt:            athlete.UpdatedAt.Time,
			HugelCount:           int(detailedLoad.HugelCount),
		},
		Load: modelsdk.AthleteLoad{
			AthleteID:                  load.AthleteID,
			LastBackloadActivityStart:  load.LastBackloadActivityStart.Time,
			LastLoadAttempt:            load.LastLoadAttempt.Time,
			LastLoadIncomplete:         load.LastLoadIncomplete,
			LastLoadError:              load.LastLoadError,
			ActivitesLoadedLastAttempt: load.ActivitesLoadedLastAttempt,
			EarliestActivity:           load.EarliestActivity.Time,
			EarliestActivityID:         load.EarliestActivityID,
			EarliestActivityDone:       load.EarliestActivityDone,
		},
		TotalActivities:  total,
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
		UpdatedAt:            athlete.UpdatedAt.Time,
		HugelCount:           int(full.HugelCount),
	})
}

func (api *API) missingSegments(rw http.ResponseWriter, r *http.Request) {
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

	missing, err := api.Opts.DB.MissingHugelSegments(ctx, actID)
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, modelsdk.Response{
			Message: "Failed to fetch missing segments",
			Detail:  err.Error(),
		})
		return
	}

	httpapi.Write(ctx, rw, http.StatusOK, missing)
}

func (api *API) manualFetchActivity(rw http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
	)

	athleteID, err := strconv.ParseInt(chi.URLParam(r, "athlete_id"), 10, 64)
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusBadRequest, modelsdk.Response{
			Message: "Invalid athlete ID",
			Detail:  err.Error(),
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

	unique, err := api.RiverManager.EnqueueFetchActivity(ctx, river2.FetchActivityArgs{
		Source:         database.ActivityDetailSourceManual,
		ActivityID:     actID,
		AthleteID:      athleteID,
		HugelPotential: true,
		OnHugelDates:   true,
	}, river2.PriorityHighest)
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusBadRequest, modelsdk.Response{
			Message: "Enqueue fetch",
			Detail:  err.Error(),
		})
		return
	}

	httpapi.Write(ctx, rw, http.StatusOK, modelsdk.Response{
		Message: fmt.Sprintf("Enqueued %d", actID),
		Detail:  fmt.Sprintf("unique: %t", unique),
	})
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
		StartDate:      activity.StartDate.Time,
		StartDateLocal: activity.StartDateLocal.Time,
		Timezone:       activity.Timezone,
	}
}

func convertHugelAthleteActivities(activities []database.AthleteHugelActivitesRow) []modelsdk.AthleteHugelActivity {
	sdk := make([]modelsdk.AthleteHugelActivity, 0, len(activities))
	for _, act := range activities {
		sdk = append(sdk, convertHugelAthleteActivity(act))
	}
	return sdk
}

func convertHugelAthleteActivity(activity database.AthleteHugelActivitesRow) modelsdk.AthleteHugelActivity {
	return modelsdk.AthleteHugelActivity{
		Summary:          convertActivitySummary(activity.ActivitySummary),
		Efforts:          convertHugelSegmentEfforts(activity.HugelActivity.Efforts),
		TotalTimeSeconds: activity.HugelActivity.TotalTimeSeconds,
	}
}

func convertHugelSegmentEfforts(dbEfforts []database.HugelSegmentEffort) []modelsdk.SegmentEffort {
	var efforts []modelsdk.SegmentEffort
	for _, e := range dbEfforts {
		efforts = append(efforts, modelsdk.SegmentEffort{
			ActivityID:   modelsdk.StringInt(e.ActivityID),
			EffortID:     modelsdk.StringInt(e.EffortID),
			StartDate:    e.StartDate,
			SegmentID:    modelsdk.StringInt(e.SegmentID),
			ElapsedTime:  int64(e.ElapsedTime),
			MovingTime:   int64(e.MovingTime),
			DeviceWatts:  e.DeviceWatts,
			AverageWatts: e.AverageWatts,
		})
	}
	return efforts
}
