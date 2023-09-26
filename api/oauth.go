package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/Emyrk/strava/api/modelsdk"

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
		httpapi.Write(ctx, rw, http.StatusInternalServerError, modelsdk.Response{
			Message: "Failed to get authenticated athlete",
			Detail:  err.Error(),
		})
		return
	}

	doLoad := false
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
			return fmt.Errorf("upsert login for %d: %w", athlete.ID, err)
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
			ProfilePicLink:        athlete.Profile,
			ProfilePicLinkMedium:  athlete.ProfileMedium,
		})
		if err != nil {
			return fmt.Errorf("upsert athlete: %w", err)
		}

		// Insert a load if we don't have one
		_, err := store.GetAthleteLoad(ctx, athlete.ID)
		if err != nil && errors.Is(err, sql.ErrNoRows) {
			doLoad = true
			// No load means we need to insert one
			_, err = store.UpsertAthleteLoad(ctx, database.UpsertAthleteLoadParams{
				AthleteID:                  athlete.ID,
				LastBackloadActivityStart:  time.Time{},
				LastLoadAttempt:            time.Time{},
				LastLoadIncomplete:         false,
				LastLoadError:              "",
				ActivitesLoadedLastAttempt: 0,
				// Start from the future
				EarliestActivity:     time.Now().Add(time.Hour * 360),
				EarliestActivityDone: false,
			})
			if err != nil {
				return fmt.Errorf("upsert load: %w", err)
			}
		} else if err != nil {
			return fmt.Errorf("get load: %w", err)
		}

		return nil
	}, nil)
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, modelsdk.Response{
			Message: "Failed to store athlete",
			Detail:  err.Error(),
		})
		return
	}

	session, err := api.Auth.CreateSession(ctx, athlete.ID)
	if err != nil {
		httpapi.Write(ctx, rw, http.StatusInternalServerError, modelsdk.Response{
			Message: "Failed to create session",
			Detail:  err.Error(),
		})
		return
	}

	http.SetCookie(rw, &http.Cookie{
		Name:     httpmw.StravaAuthJWTCookie,
		Path:     "/",
		Value:    session,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		// 7 Days
		MaxAge: 7 * 24 * 60 * 60,
	})

	logger.Info().
		Str("username", athlete.Username).
		Str("firstname", athlete.Firstname).
		Str("lastname", athlete.Lastname).
		Int64("id", athlete.ID).
		Str("redirect", state.Redirect).
		Bool("do_load", doLoad).
		Msg("Authenticated Athlete")

	http.Redirect(rw, r, state.Redirect, http.StatusSeeOther)
}

func (api *API) logout(rw http.ResponseWriter, r *http.Request) {
	http.SetCookie(rw, &http.Cookie{
		Name:   httpmw.StravaAuthJWTCookie,
		Path:   "/",
		MaxAge: -1,
	})

	// TODO: Pull redirect url from query param if present.
	http.Redirect(rw, r, "/signed-out", http.StatusSeeOther)
}
