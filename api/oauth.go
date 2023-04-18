package api

import (
	"fmt"
	"net/http"

	"github.com/Emyrk/strava/api/httpapi"
	"github.com/Emyrk/strava/api/httpmw"
	"github.com/Emyrk/strava/strava"
	"golang.org/x/oauth2"
)

func (api *API) stravaOAuth2(rw http.ResponseWriter, r *http.Request) {
	var (
		ctx   = r.Context()
		state = httpmw.OAuth2(r)
	)

	oauthClient := oauth2.NewClient(ctx, oauth2.StaticTokenSource(state.Token))
	scli := strava.NewOAuthClient(oauthClient)
	athlete, err := scli.GetAuthenticatedAthelete(ctx)
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, httpapi.Response{
			Message: "Failed to get authenticated athlete",
			Detail:  err.Error(),
		})
		return
	}

	fmt.Println(athlete)
	httpapi.Write(ctx, rw, http.StatusOK, athlete)
}
