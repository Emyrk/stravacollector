package api

import (
	"fmt"
	"net/http"

	"github.com/Emyrk/strava/api/httpmw"

	"github.com/go-chi/chi/v5"
	"golang.org/x/oauth2"
)

type Options struct {
	OAuthCfg *oauth2.Config
}

type API struct {
	OAuthCfg *oauth2.Config
	Handler  http.Handler
}

func New(opts Options) (*API, error) {
	if opts.OAuthCfg == nil {
		return nil, fmt.Errorf("missing oauth2 config")
	}

	api := &API{
		OAuthCfg: opts.OAuthCfg,
	}
	api.Handler = api.Routes()

	return api, nil
}

func (api *API) Routes() http.Handler {
	r := chi.NewRouter()

	r.Route("/oauth2", func(r chi.Router) {
		r.Use(httpmw.ExtractOauth2(api.OAuthCfg, nil))
		r.Get("/callback", api.stravaOAuth2)
	})

	return r
}
