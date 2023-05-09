package api

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"

	"github.com/Emyrk/strava/database"
	"github.com/go-chi/chi/v5"

	"github.com/Emyrk/strava/api/httpapi"
	"github.com/Emyrk/strava/api/httpmw"
	"github.com/Emyrk/strava/api/modelsdk"
)

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
