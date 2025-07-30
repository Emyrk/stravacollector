package cli

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Emyrk/strava/api/river"
	"github.com/Emyrk/strava/api/webhooks"
	"github.com/Emyrk/strava/database"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
)

func rehook(ctx context.Context, db database.Store, dbURL string, logger zerolog.Logger) error {
	riverManager, err := river.New(ctx, river.Options{
		DBURL:      dbURL,
		Logger:     logger.With().Str("component", "river").Logger(),
		DB:         db,
		Registry:   prometheus.NewRegistry(),
		InsertOnly: true,
	})
	if err != nil {
		return fmt.Errorf("create river manager: %w", err)
	}
	defer riverManager.Close(ctx)

	hooks, err := db.GetDeleteActivityWebhooks(ctx)
	if err != nil {
		return fmt.Errorf("get delete activity webhooks: %w", err)
	}

	for _, hook := range hooks {
		var evt webhooks.WebhookEvent

		err := json.Unmarshal([]byte(hook.Raw), &evt)
		if err != nil {
			return fmt.Errorf("unmarshal webhook event: %w", err)
		}
		riverManager.HandleWebhookEvent(ctx, evt)
		fmt.Printf("Webhook dump %s handled for activity %d\n", hook.ID.String(), evt.ObjectID)
	}

	return nil
}
