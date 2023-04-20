package webhooks

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/Emyrk/strava/database"

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
	DB          database.Store

	ID int
}

func NewActivityEvents(logger zerolog.Logger, cfg *oauth2.Config, db database.Store, accessURL *url.URL) *ActivityEvents {
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
		DB:          db,
	}
}

func (a *ActivityEvents) Setup(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	hooks, err := a.ViewWebhook(ctx)
	if err == nil {
		for _, h := range hooks {
			// Always reset the hook
			err := a.DeleteWebhook(ctx, h.ID)
			if err != nil {
				return fmt.Errorf("error deleting webhook: %w", err)
			}
		}
	}

	id, err := stravawebhook.CreateWebhook(ctx, a.OauthConfig.ClientID, a.OauthConfig.ClientSecret, a.Callback.String(), a.VerifyToken)
	if err != nil {
		return fmt.Errorf("error creating webhook: %w", err)
	}
	a.ID = id
	return nil
}

func (a *ActivityEvents) Close() {

}

// WebhookEvent is documented https://developers.strava.com/docs/webhooks/
type WebhookEvent struct {
	// AspectType always "create," "update," or "delete."
	AspectType string `json:"aspect_type"`
	// EventTime is unix seconds
	EventTime int `json:"event_time"`
	// ObjectID is the activity's ID or athlete's ID.
	ObjectID int64 `json:"object_id"`
	// ObjectType either "activity" or "athlete."
	ObjectType string `json:"object_type"`
	// OwnerID is the athlete ID
	OwnerID        int64             `json:"owner_id"`
	SubscriptionID int               `json:"subscription_id"`
	Updates        map[string]string `json:"updates"`
}

func (a *ActivityEvents) handleWebhook(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Always save the payload
	d, _ := io.ReadAll(r.Body)
	dump, err := a.DB.InsertWebhookDump(ctx, string(d))
	if err != nil {
		a.Logger.Error().Err(err).Msg("error inserting webhook dump")
	}
	logger := a.Logger.With().Str("id", dump.ID.String()).Logger()

	var event WebhookEvent
	err = json.Unmarshal(d, &event)
	if err != nil {
		logger.Error().
			Str("body", string(d)).
			Err(err).Msg("error unmarshalling webhook event")
		return
	}

	switch event.ObjectType {
	case "activity":
		a.newActivity(event, logger)
	case "athlete":
		// Ignore these for now.
	default:
		logger.Warn().
			Str("object_type", event.ObjectType).
			Msg("Webhook event not supported")
	}
}

func (a *ActivityEvents) newActivity(event WebhookEvent, logger zerolog.Logger) {
	switch event.AspectType {
	case "create":
	case "update":
		// Updates to events we probably don't care about?
		logger.Info().
			Interface("updated", event.Updates).
			Msg("'Update' webhook event to an activity")

	case "delete":
	}
}

func (a *ActivityEvents) Attach(r chi.Router) chi.Router {
	r.Route("/webhooks/strava", func(r chi.Router) {
		r.Use(
			func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
					a.Logger.Info().
						Str("remote_addr", r.RemoteAddr).
						Msg("Strava webhook received")
					challenge := r.URL.Query().Get("hub.challenge")
					if challenge != "" {
						token := r.URL.Query().Get("hub.verify_token")
						if token != a.VerifyToken {
							d, _ := io.ReadAll(r.Body)
							a.Logger.Warn().
								Str("found-token", token).
								Str("expected-token", a.VerifyToken).
								Str("url", r.URL.String()).
								Str("body", string(d)).
								Str("method", r.Method).
								Interface("headers", r.Header).
								Msg("Strava webhook token mismatch")

							rw.WriteHeader(http.StatusUnauthorized)
							return
						}

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
		r.HandleFunc("/", a.handleWebhook)
	})
	return r
}

func (a *ActivityEvents) ViewWebhook(ctx context.Context) ([]stravawebhook.Webhook, error) {
	return stravawebhook.ViewWebhook(ctx, a.OauthConfig.ClientID, a.OauthConfig.ClientSecret)
}

func (a *ActivityEvents) CreateWebhook(ctx context.Context) (int, error) {
	return stravawebhook.CreateWebhook(ctx,
		a.OauthConfig.ClientID,
		a.OauthConfig.ClientSecret,
		a.Callback.String(),
		a.VerifyToken,
	)
}

func (a *ActivityEvents) DeleteWebhook(ctx context.Context, id int) error {
	return stravawebhook.DeleteWebhook(ctx,
		a.OauthConfig.ClientID,
		a.OauthConfig.ClientSecret,
		id,
	)
}
