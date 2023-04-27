package api

import (
	"net/http"

	"github.com/Emyrk/strava/api/httpapi"
	"github.com/Emyrk/strava/api/modelsdk"

	"github.com/Emyrk/strava/api/httpmw"
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

	httpapi.Write(ctx, rw, http.StatusOK, modelsdk.AthleteLogin{
		AthleteID: id,
		Summit:    login.Summit,
	})
}
