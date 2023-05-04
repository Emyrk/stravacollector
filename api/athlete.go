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
