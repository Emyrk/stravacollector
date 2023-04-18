package api

import (
	"encoding/json"
	"net/http"

	"github.com/Emyrk/strava/api/httpapi"
	"github.com/Emyrk/strava/api/httpmw"
	"github.com/Emyrk/strava/database"
	"github.com/Emyrk/strava/strava"
	"golang.org/x/oauth2"
)

func (api *API) stravaOAuth2(rw http.ResponseWriter, r *http.Request) {
	var (
		ctx    = r.Context()
		state  = httpmw.OAuth2(r)
		logger = api.Opts.Logger
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

	raw, _ := json.Marshal(athlete)
	_, err = api.Opts.DB.UpsertAthlete(ctx, database.UpsertAthleteParams{
		ID:                athlete.ID,
		Premium:           athlete.Premium,
		Username:          athlete.Username,
		Firstname:         athlete.Firstname,
		Lastname:          athlete.Lastname,
		Sex:               athlete.Sex,
		ProviderID:        api.Opts.OAuthCfg.ClientID,
		OauthAccessToken:  state.Token.AccessToken,
		OauthRefreshToken: state.Token.RefreshToken,
		OauthExpiry:       state.Token.Expiry,
		Raw:               string(raw),
	})
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, httpapi.Response{
			Message: "Failed to store athlete",
			Detail:  err.Error(),
		})
		return
	}

	logger.Info().
		Str("username", athlete.Username).
		Str("firstname", athlete.Firstname).
		Str("lastname", athlete.Lastname).
		Int64("id", athlete.ID).
		Msg("Authenticated Athlete")

	httpapi.Write(ctx, rw, http.StatusOK, athlete)
}
