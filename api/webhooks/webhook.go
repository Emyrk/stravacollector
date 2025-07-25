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
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog"
	"golang.org/x/oauth2"

	"github.com/Emyrk/strava/api/httpapi"
	"github.com/Emyrk/strava/database"
	"github.com/Emyrk/strava/strava/stravalimit"
	"github.com/Emyrk/strava/strava/stravawebhook"
)

type ActivityEvents struct {
	OauthConfig *oauth2.Config
	AccessURL   *url.URL
	Callback    *url.URL
	VerifyToken string
	Logger      zerolog.Logger
	DB          database.Store

	eventQueue chan *WebhookEvent

	webhookCount *prometheus.GaugeVec

	ID int
}

func NewActivityEvents(logger zerolog.Logger, cfg *oauth2.Config, db database.Store, accessURL *url.URL, verifyToken string, registry prometheus.Registerer) *ActivityEvents {
	if verifyToken == "" {
		vData := make([]byte, 32)
		_, err := rand.Read(vData)
		if err != nil {
			panic(err)
		}
		hex.EncodeToString(vData)
	}
	callback := *accessURL
	callback.Path = "/webhooks/strava"
	factory := promauto.With(registry)

	return &ActivityEvents{
		OauthConfig: cfg,
		AccessURL:   accessURL,
		VerifyToken: verifyToken,
		Callback:    &callback,
		Logger:      logger,
		DB:          db,
		eventQueue:  make(chan *WebhookEvent, 100),
		webhookCount: factory.NewGaugeVec(prometheus.GaugeOpts{
			Namespace:   "strava",
			Subsystem:   "api_webhooks",
			Name:        "webhook_count",
			Help:        "Number of webhooks received",
			ConstLabels: nil,
		}, []string{"type"}),
	}
}

func (a *ActivityEvents) EventQueue() <-chan *WebhookEvent {
	return a.eventQueue
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
		return fmt.Errorf("error creating webhook (%s): %w", a.Callback.String(), err)
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
	EventTime int64 `json:"event_time"`
	// ObjectID is the activity's ID or athlete's ID.
	ObjectID int64 `json:"object_id"`
	// ObjectType either "activity" or "athlete."
	ObjectType string `json:"object_type"`
	// OwnerID is the athlete ID
	OwnerID        int64             `json:"owner_id"`
	SubscriptionID int               `json:"subscription_id"`
	Updates        map[string]string `json:"updates"`

	once sync.Once
	done chan struct{}
}

func (e *WebhookEvent) MarkDone() {
	e.once.Do(func() {
		close(e.done)
	})
}

func (e *WebhookEvent) WaitDone() {
	<-e.done
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

	event.done = make(chan struct{})
	a.eventQueue <- &event
	a.webhookCount.WithLabelValues(event.AspectType).Inc()

	event.WaitDone()
	_, _ = rw.Write([]byte("Thanks!"))
}

func (a *ActivityEvents) Attach(r chi.Router) chi.Router {
	r.Route("/webhooks/strava", func(r chi.Router) {
		r.Use(
			func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
					stravalimit.Update(r.Header)

					challenge := r.URL.Query().Get("hub.challenge")
					a.Logger.Info().
						Str("remote_addr", r.RemoteAddr).
						Str("challenge", challenge).
						Interface("query", r.URL.Query()).
						Msg("Strava webhook received")
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
	a.Logger.Debug().
		Msg("view webhook api call")
	return stravawebhook.ViewWebhook(ctx, a.OauthConfig.ClientID, a.OauthConfig.ClientSecret)
}

func (a *ActivityEvents) CreateWebhook(ctx context.Context) (int, error) {
	a.Logger.Debug().
		Msg("create webhook api call")
	return stravawebhook.CreateWebhook(ctx,
		a.OauthConfig.ClientID,
		a.OauthConfig.ClientSecret,
		a.Callback.String(),
		a.VerifyToken,
	)
}

func (a *ActivityEvents) DeleteWebhook(ctx context.Context, id int) error {
	a.Logger.Debug().
		Msg("delete webhook api call")
	return stravawebhook.DeleteWebhook(ctx,
		a.OauthConfig.ClientID,
		a.OauthConfig.ClientSecret,
		id,
	)
}
