package api

import (
	"fmt"
	"net/http"
	"time"

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

	err = api.Opts.DB.InTx(func(store database.Store) error {
		_, err = store.UpsertAthleteLogin(ctx, database.UpsertAthleteLoginParams{
			AthleteID:         athlete.ID,
			Summit:            athlete.Premium || athlete.Summit,
			ProviderID:        api.OAuthConfig.ClientID,
			OauthAccessToken:  state.Token.AccessToken,
			OauthRefreshToken: state.Token.RefreshToken,
			OauthExpiry:       state.Token.Expiry,
			OauthTokenType:    state.Token.TokenType,
		})
		if err != nil {
			return fmt.Errorf("upsert login: %w", err)
		}

		_, err = store.UpsertAthleteLoad(ctx, database.UpsertAthleteLoadParams{
			AthleteID:                 athlete.ID,
			LastBackloadActivityStart: time.Time{},
			LastLoadAttempt:           time.Time{},
			LastLoadIncomplete:        false,
		})
		if err != nil {
			return fmt.Errorf("upsert load: %w", err)
		}

		_, err = store.UpsertAthlete(ctx, database.UpsertAthleteParams{
			ID:                    athlete.ID,
			CreatedAt:             athlete.CreatedAt,
			UpdatedAt:             athlete.UpdatedAt,
			Summit:                athlete.Summit || athlete.Premium,
			Username:              athlete.Username,
			Firstname:             athlete.Firstname,
			Lastname:              athlete.Lastname,
			Sex:                   athlete.Sex,
			City:                  athlete.City,
			State:                 athlete.State,
			Country:               athlete.Country,
			FollowCount:           int32(athlete.FollowerCount),
			FriendCount:           int32(athlete.FriendCount),
			MeasurementPreference: athlete.MeasurementPreference,
			Ftp:                   athlete.Ftp,
			Weight:                athlete.Weight,
			Clubs:                 athlete.Clubs,
		})
		if err != nil {
			return fmt.Errorf("upsert athlete: %w", err)
		}

		return nil
	}, nil)
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
