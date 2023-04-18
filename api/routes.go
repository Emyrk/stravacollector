package api

import (
	"net/http"

	"github.com/rs/zerolog"

	"github.com/Emyrk/strava/database"

	"github.com/Emyrk/strava/api/httpmw"

	"github.com/go-chi/chi/v5"
	"golang.org/x/oauth2"
)

type Options struct {
	OAuthCfg *oauth2.Config
	DB       database.Store
	Logger   zerolog.Logger
}

type API struct {
	Opts    *Options
	Handler http.Handler
}

func New(opts Options) (*API, error) {
	api := &API{
		Opts: &opts,
	}
	api.Handler = api.Routes()

	return api, nil
}

func (api *API) Routes() http.Handler {
	r := chi.NewRouter()

	r.Route("/oauth2", func(r chi.Router) {
		r.Use(httpmw.ExtractOauth2(api.Opts.OAuthCfg, nil))
		r.Get("/callback", api.stravaOAuth2)
	})

	return r
}
