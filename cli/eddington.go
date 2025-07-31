package cli

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/Emyrk/strava/api/webhooks"
	"github.com/Emyrk/strava/database"
	"github.com/rs/zerolog"
)

func eddington(ctx context.Context, db database.Store, logger zerolog.Logger) error {
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
