package webhooks

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/Emyrk/strava/api/httpapi"
	"github.com/Emyrk/strava/strava/stravawebhook"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"golang.org/x/oauth2"
)

type ActivityEvents struct {
	OauthConfig *oauth2.Config
	AccessURL   *url.URL
	Callback    *url.URL
	VerifyToken string
	Logger      zerolog.Logger

	ID int
}

func NewActivityEvents(logger zerolog.Logger, cfg *oauth2.Config, accessURL *url.URL) *ActivityEvents {
	vData := make([]byte, 32)
	_, err := rand.Read(vData)
	if err != nil {
		panic(err)
	}
	callback := *accessURL
	callback.Path = "/webhooks/strava"

	return &ActivityEvents{
		OauthConfig: cfg,
		AccessURL:   accessURL,
		VerifyToken: hex.EncodeToString(vData),
		Callback:    &callback,
		Logger:      logger,
	}
}

func (a *ActivityEvents) Setup(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	fmt.Println(a.OauthConfig.ClientID, a.OauthConfig.ClientSecret)
	//err := stravawebhook.ViewWebhook(ctx, a.OauthConfig.ClientID, a.OauthConfig.ClientSecret)
	//if err != nil {
	//	return fmt.Errorf("error viewing webhook: %w", err)
	//}

	err := stravawebhook.CreateWebhook(ctx, a.OauthConfig.ClientID, a.OauthConfig.ClientSecret, a.Callback.String(), a.VerifyToken)
	if err != nil {
		return fmt.Errorf("error creating webhook: %w", err)
	}
	return nil
}

func (a *ActivityEvents) Close() {

}

func (a *ActivityEvents) handleWebhook(rw http.ResponseWriter, r *http.Request) {
	d, _ := io.ReadAll(r.Body)
	fmt.Println(string(d))
}

func (a *ActivityEvents) Attach(r chi.Router) {
	r.Route("/webhooks/strava", func(r chi.Router) {
		r.Use(
			func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
					token := r.URL.Query().Get("hub.verify_token")
					if token != a.VerifyToken {
						rw.WriteHeader(http.StatusUnauthorized)
						return
					}
					next.ServeHTTP(rw, r)
				})
			},
			func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
					challenge := r.URL.Query().Get("hub.challenge")
					if challenge != "" {
						httpapi.Write(r.Context(), rw, http.StatusOK, struct {
							Challenge string `json:"hub.challenge"`
						}{
							Challenge: challenge,
						})
						a.Logger.Info().Msg("Strava webhook challenge returned")
						return
					}
					next.ServeHTTP(rw, r)
				})
			},
		)
		r.Get("/", a.handleWebhook)
	})
}
