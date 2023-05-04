package api

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/Emyrk/strava/api/httpapi"
	"github.com/Emyrk/strava/api/httpmw"
	"github.com/Emyrk/strava/api/modelsdk"
)

func (api *API) whoAmI(rw http.ResponseWriter, r *http.Request) {
	var (
		ctx = r.Context()
		id  = httpmw.AuthenticatedAthleteID(r)
	)

	login, err := api.Opts.DB.GetAthleteLogin(ctx, id)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	athlete, err := api.Opts.DB.GetAthlete(ctx, id)
	if err != nil && errors.Is(err, sql.ErrNoRows) {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	if errors.Is(err, sql.ErrNoRows) {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, modelsdk.Response{
			Message: "Failed to fetch authenticated athlete",
			Detail:  "Please try to log out and login again.",
		})
		return
	}

	httpapi.Write(ctx, rw, http.StatusOK, modelsdk.AthleteSummary{
		AthleteID:            modelsdk.Int64String(id),
		Summit:               login.Summit,
		Username:             athlete.Username,
		Firstname:            athlete.Firstname,
		Lastname:             athlete.Lastname,
		Sex:                  athlete.Sex,
		ProfilePicLink:       athlete.ProfilePicLink,
		ProfilePicLinkMedium: athlete.ProfilePicLinkMedium,
		UpdatedAt:            athlete.UpdatedAt,
	})
}
