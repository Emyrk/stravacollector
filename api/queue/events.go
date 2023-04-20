package queue

import (
	"context"

	"github.com/Emyrk/strava/api/webhooks"
)

func (m *Manager) HandleWebhookEvents(ctx context.Context, c <-chan *webhooks.WebhookEvent) {
	for {
		select {
		case <-ctx.Done():
			return
		case event := <-c:
			switch event.ObjectType {
			case "activity":
				m.newActivity(ctx, *event)
			case "athlete":
				// Ignore these for now.
			default:
				m.Logger.Warn().
					Str("object_type", event.ObjectType).
					Msg("Webhook event not supported")
			}
		}
	}
}

func (m *Manager) newActivity(ctx context.Context, event webhooks.WebhookEvent) {
	switch event.AspectType {
	case "create":
		actID := event.ObjectID
		err := m.EnqueueFetchActivity(ctx, event.OwnerID, actID)
		if err != nil {
			m.Logger.Error().
				Err(err).
				Int64("owner_id", event.OwnerID).
				Int64("activity_id", actID).
				Msg("error enqueueing activity")
		}
	case "update":
		// Updates to events we probably don't care about?
		m.Logger.Info().
			Interface("updated", event.Updates).
			Msg("'Update' webhook event to an activity")
	case "delete":
	}
}
