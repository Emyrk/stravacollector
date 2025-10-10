package api

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Emyrk/strava/api/river"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"golang.org/x/oauth2"

	"github.com/Emyrk/strava/api/auth"
	"github.com/Emyrk/strava/api/httpapi"
	"github.com/Emyrk/strava/api/httpmw"
	"github.com/Emyrk/strava/api/modelsdk"
	"github.com/Emyrk/strava/api/webhooks"
	"github.com/Emyrk/strava/database"
	"github.com/Emyrk/strava/database/gencache"
	"github.com/Emyrk/strava/internal/hugeldate"
	server "github.com/Emyrk/strava/site"
)

type Options struct {
	OAuth         OAuthOptions
	DB            database.Store
	Logger        zerolog.Logger
	AccessURL     *url.URL
	VerifyToken   string
	SigningKeyPEM []byte
	Registry      *prometheus.Registry
}

type OAuthOptions struct {
	ClientID string
	Secret   string
}

type API struct {
	Opts    *Options
	Handler chi.Router

	Auth         *auth.Authentication
	OAuthConfig  *oauth2.Config
	Events       *webhooks.ActivityEvents
	RiverManager *river.Manager

	SuperHugelBoardCache    *gencache.LazyCache[[]database.SuperHugelLeaderboardRow]
	HugelBoard2023Cache     *gencache.LazyCache[[]database.HugelLeaderboardRow]
	HugelBoard2024Cache     *gencache.LazyCache[[]database.HugelLeaderboardRow]
	HugelBoard2024LiteCache *gencache.LazyCache[[]database.HugelLeaderboardRow]

	HugelRouteCache     *gencache.LazyCache[database.GetCompetitiveRouteRow]
	HugelLiteRouteCache *gencache.LazyCache[database.GetCompetitiveRouteRow]

	// Metrics
	Registry *prometheus.Registry
}

func New(ctx context.Context, opts Options) (*API, error) {
	if opts.Registry == nil {
		opts.Registry = prometheus.NewRegistry()
	}
	api := &API{
		Opts: &opts,
		OAuthConfig: &oauth2.Config{
			ClientID:     opts.OAuth.ClientID,
			ClientSecret: opts.OAuth.Secret,
			Endpoint: oauth2.Endpoint{
				AuthURL:   "https://www.strava.com/oauth/authorize",
				TokenURL:  "https://www.strava.com/oauth/token",
				AuthStyle: 0,
			},
			RedirectURL: fmt.Sprintf("%s/oauth2/callback", strings.TrimSuffix(opts.AccessURL.String(), "/")),
			// Must be comma joined
			Scopes: []string{strings.Join([]string{"read", "read_all", "profile:read_all", "activity:read"}, ",")},
		},
		Registry: opts.Registry,
	}
	ath, err := auth.New(auth.Options{
		Lifetime:  time.Hour * 24 * 7,
		SecretPEM: opts.SigningKeyPEM,
		Issuer:    "Strava-Hugel",
		Registry:  api.Registry,
	})
	if err != nil {
		return nil, fmt.Errorf("create auth: %w", err)
	}
	api.Auth = ath

	api.Events = webhooks.NewActivityEvents(opts.Logger, api.OAuthConfig, api.Opts.DB, opts.AccessURL, opts.VerifyToken, api.Registry)
	r := api.Routes()
	r = api.Events.Attach(r)
	api.Handler = r

	api.SuperHugelBoardCache = gencache.New(ctx, time.Hour*24, func(ctx context.Context) ([]database.SuperHugelLeaderboardRow, error) {
		return api.Opts.DB.SuperHugelLeaderboard(ctx, 0)
	})

	api.HugelBoard2023Cache = gencache.New(ctx, time.Hour*48, func(ctx context.Context) ([]database.HugelLeaderboardRow, error) {
		return api.Opts.DB.YearlyHugelLeaderboard(ctx, database.YearlyHugelLeaderboardParams{
			RouteYear: 2023,
			HugelLeaderboardParams: database.HugelLeaderboardParams{
				AthleteID: -1,
				After:     database.Timestamp(hugeldate.Year2023.Start),
				Before:    database.Timestamp(hugeldate.Year2023.End),
			},
		})
	})
	api.HugelBoard2024Cache = gencache.New(ctx, time.Hour*48, func(ctx context.Context) ([]database.HugelLeaderboardRow, error) {
		return api.Opts.DB.YearlyHugelLeaderboard(ctx, database.YearlyHugelLeaderboardParams{
			RouteYear: 2024,
			HugelLeaderboardParams: database.HugelLeaderboardParams{
				AthleteID: -1,
				After:     database.Timestamp(hugeldate.Year2024.Start),
				Before:    database.Timestamp(hugeldate.Year2024.End),
			},
		})
	})
	api.HugelBoard2024LiteCache = gencache.New(ctx, time.Hour*48, func(ctx context.Context) ([]database.HugelLeaderboardRow, error) {
		return api.Opts.DB.YearlyHugelLeaderboard(ctx, database.YearlyHugelLeaderboardParams{
			RouteYear: 2024,
			Lite:      true,
			HugelLeaderboardParams: database.HugelLeaderboardParams{
				AthleteID: -1,
				After:     database.Timestamp(hugeldate.Year2024.Start),
				Before:    database.Timestamp(hugeldate.Year2024.End),
			},
		})
	})
	api.HugelRouteCache = gencache.New(ctx, time.Hour*4, func(ctx context.Context) (database.GetCompetitiveRouteRow, error) {
		return api.Opts.DB.GetCompetitiveRoute(ctx, "das-hugel")
	})
	api.HugelLiteRouteCache = gencache.New(ctx, time.Hour*4, func(ctx context.Context) (database.GetCompetitiveRouteRow, error) {
		return api.Opts.DB.GetCompetitiveRoute(ctx, "lite-das-hugel")
	})

	return api, nil
}

func (api *API) HugelRoute(ctx context.Context, lite bool) (database.GetCompetitiveRouteRow, error) {
	if lite {
		return api.HugelLiteRouteCache.Load(ctx)
	}
	return api.HugelRouteCache.Load(ctx)
}

// StartWebhook needs to be called after the API is served.
func (api *API) StartWebhook(ctx context.Context, setup bool) (<-chan *webhooks.Handled[webhooks.WebhookEvent], error) {
	if setup {
		err := api.Events.Setup(ctx)
		if err != nil {
			return nil, err
		}
	}
	return api.Events.EventQueue(), nil
}

func (api *API) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(
		httpmw.PrometheusMW(api.Registry),
		httpmw.NoWWW(),
	)

	r.Get("/myhealthz", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})
	r.Get("/static-error", func(w http.ResponseWriter, r *http.Request) {
		httpapi.Write(r.Context(), w, http.StatusInternalServerError, modelsdk.Response{
			Message: "This error is ok, and is used for testing",
			Detail:  "Maybe if you try again really hard it will work!",
		})
	})
	r.Route("/oauth2", func(r chi.Router) {
		r.Use(httpmw.ExtractOauth2(api.Opts.Logger, api.OAuthConfig, nil))
		r.Get("/callback", api.stravaOAuth2)
	})
	r.Route("/api/v1", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Route("/athlete", func(r chi.Router) {
				r.Route("/{athlete_id}/", func(r chi.Router) {
					r.Use(
						httpmw.Authenticated(api.Auth, true),
						httpmw.ExtractAthlete(api.Opts.DB),
					)
					r.Get("/", api.athlete)
					r.Get("/hugels", api.athleteHugels)
					r.Get("/sync-summary", api.syncSummary)
					r.Get("/eddington", api.eddingtonNumber)
				})
			})
			r.Route("/athletes", func(r chi.Router) {
				r.Get("/eddington", api.allEddingtons)
			})
		})
		r.Group(func(r chi.Router) {
			// Authenticated routes
			r.Use(
				httpmw.Authenticated(api.Auth, false),
			)
			r.Get("/whoami", api.whoAmI)
			r.Route("/fetch-activity", func(r chi.Router) {
				r.Use(httpmw.AuthenticatedAsAdmins())
				r.Get("/{athlete_id}-{activity_id}", api.manualFetchActivity)
			})
			r.Route("/forward-load", func(r chi.Router) {
				r.Use(httpmw.AuthenticatedAsAdmins())
				r.Get("/{athlete_id}", api.forwardLoadAthlete)
			})
			r.Route("/eddington", func(r chi.Router) {
				r.Use(httpmw.AuthenticatedAsAdmins())
				r.Get("/{athlete_id}", api.eddington)
			})
			r.Route("/missing", func(r chi.Router) {
				r.Get("/{activity_id}", api.missingSegments)
			})
		})
		r.Group(func(r chi.Router) {
			// Unauthenticated routes
			r.Use(
				httpmw.Authenticated(api.Auth, true),
			)
			r.Get("/superhugelboard", api.superHugelboard)
			r.Get("/hugelboard", api.hugelboard)
			r.Route("/route", func(r chi.Router) {
				r.Get("/{route-name}", api.competitiveRoute)
				r.Get("/{route-name}/verify/{route-id}", api.verifyRoute)
			})
			r.Route("/segments", func(r chi.Router) {
				r.Post("/", api.getSegments)
			})
		})
		r.NotFound(api.apiNotFound)
	})
	r.Get("/logout", api.logout)
	r.NotFound(server.Handler(server.FS()).ServeHTTP)

	return r
}

func (api *API) apiNotFound(w http.ResponseWriter, r *http.Request) {
	httpapi.Write(r.Context(), w, http.StatusNotFound, modelsdk.Response{
		Message: "Not found",
		Detail:  "api route not found",
	})
}
