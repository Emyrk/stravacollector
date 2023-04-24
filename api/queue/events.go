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
				m.Logger.Warn().
					Interface("event", event).
					Msg("Webhook event to an athlete not handled")
			default:
				m.Logger.Warn().
					Str("object_type", event.ObjectType).
					Msg("Webhook event not supported")
			}
		}
	}
}

func (m *Manager) newActivity(ctx context.Context, event webhooks.WebhookEvent) {
	var err error
	switch event.AspectType {
	case "create":
		err = m.EnqueueFetchActivity(ctx, event.OwnerID, event.ObjectID)
	case "update":
		err = m.EnqueueUpdateActivity(ctx, event)
	case "delete":
		m.Logger.Info().
			Interface("deleted", event.Updates).
			Msg("'Delete' webhook event to an activity")
	default:
		m.Logger.Warn().
			Str("aspect_type", event.AspectType).
			Msg("Webhook event not supported")
	}
	if err != nil {
		m.Logger.Error().
			Err(err).
			Str("aspect_type", event.AspectType).
			Int64("owner_id", event.OwnerID).
			Int64("activity_id", event.ObjectID).
			Msg("error enqueueing activity")
	}
}
