package api

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"golang.org/x/oauth2"

	"github.com/Emyrk/strava/api/httpapi"
	"github.com/Emyrk/strava/api/modelsdk"
	"github.com/Emyrk/strava/database"
	"github.com/Emyrk/strava/database/gencache"
	"github.com/Emyrk/strava/strava"
)

// refreshAthleteResponse excludes leaderboard-derived counts and OAuth credentials.
type refreshAthleteResponse struct {
	AthleteID      modelsdk.StringInt `json:"athlete_id"`
	Username       string             `json:"username"`
	Firstname      string             `json:"firstname"`
	Lastname       string             `json:"lastname"`
	Sex            string             `json:"sex"`
	ProfilePicLink string             `json:"profile_pic_link"`
}

func (api *API) refreshAthlete(rw http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
	defer cancel()
	fail := func(status int, message string) { httpapi.Write(ctx, rw, status, modelsdk.Response{Message: message}) }
	id, err := strconv.ParseInt(chi.URLParam(r, "athlete_id"), 10, 64)
	if err != nil || id <= 0 {
		fail(http.StatusBadRequest, "Invalid athlete ID")
		return
	}
	login, err := api.Opts.DB.GetAthleteLogin(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		fail(http.StatusNotFound, "Athlete has no stored Strava authorization; they must reconnect")
		return
	}
	if err != nil {
		fail(http.StatusInternalServerError, "Failed to load athlete authorization")
		return
	}
	token, err := api.OAuthConfig.TokenSource(ctx, login.OAuthToken()).Token()
	if err != nil {
		fail(http.StatusBadGateway, "Failed to refresh Strava authorization; the athlete may need to reconnect")
		return
	}
	// Persist rotated credentials before the profile request, even if that request fails.
	if token.AccessToken != login.OauthAccessToken || token.RefreshToken != login.OauthRefreshToken || !token.Expiry.Equal(login.OauthExpiry.Time) || token.TokenType != login.OauthTokenType {
		_, err = api.Opts.DB.UpsertAthleteLogin(ctx, database.UpsertAthleteLoginParams{
			AthleteID: id, Summit: login.Summit, ProviderID: login.ProviderID,
			OauthAccessToken: token.AccessToken, OauthRefreshToken: token.RefreshToken,
			OauthExpiry: database.Timestamptz(token.Expiry), OauthTokenType: token.TokenType,
		})
		if err != nil {
			fail(http.StatusInternalServerError, "Failed to save refreshed Strava authorization")
			return
		}
	}
	client := strava.NewOAuthClient(oauth2.NewClient(ctx, oauth2.StaticTokenSource(token)))
	athlete, err := client.GetAuthenticatedAthelete(ctx)
	if err != nil {
		status := http.StatusBadGateway
		if se := strava.IsAPIError(err); se != nil && se.Response != nil && se.Response.StatusCode == http.StatusTooManyRequests {
			status = http.StatusTooManyRequests
		}
		fail(status, "Failed to fetch Strava profile; authorization may need reconnecting or Strava may be unavailable")
		return
	}
	if athlete.ID != id {
		fail(http.StatusBadGateway, "Strava returned a different athlete; profile was not updated")
		return
	}
	if _, err := upsertStravaAthlete(ctx, api.Opts.DB, athlete); err != nil {
		fail(http.StatusInternalServerError, "Failed to save Strava profile")
		return
	}
	api.invalidateLeaderboardCaches()
	httpapi.Write(ctx, rw, http.StatusOK, refreshAthleteResponse{
		AthleteID: modelsdk.StringInt(id), Username: athlete.Username,
		Firstname: athlete.Firstname, Lastname: athlete.Lastname, Sex: athlete.Sex,
		ProfilePicLink: athlete.Profile,
	})
}

func (api *API) invalidateLeaderboardCaches() {
	for _, cache := range []*gencache.LazyCache[[]database.HugelLeaderboardRow]{
		api.HugelBoard2023Cache, api.HugelBoard2024Cache, api.HugelBoard2024LiteCache,
		api.HugelBoard2025Cache, api.HugelBoard2025LiteCache, api.HugelBoard2026Cache, api.HugelBoard2026LiteCache,
	} {
		if cache != nil {
			cache.Invalidate()
		}
	}
	if api.SuperHugelBoardCache != nil {
		api.SuperHugelBoardCache.Invalidate()
	}
}
