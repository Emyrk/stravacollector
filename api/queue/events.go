package queue

import (
	"context"

	"github.com/vgarvardt/gue/v5"

	"github.com/Emyrk/strava/database"

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
				m.newAthlete(ctx, *event)
			default:
				m.Logger.Warn().
					Str("object_type", event.ObjectType).
					Msg("Webhook event not supported")
			}
		}
	}
}

func (m *Manager) newAthlete(ctx context.Context, event webhooks.WebhookEvent) {
	var qErr error
	switch event.AspectType {
	case "create":
		m.Logger.Warn().
			Interface("event", event).
			Msg("Webhook create event to an athlete not handled")
	case "update":
		qErr = m.EnqueueUpdateAthlete(ctx, event)
	case "delete":
		m.Logger.Warn().
			Interface("event", event).
			Msg("Webhook delete event to an athlete not handled")
	default:
		m.Logger.Warn().
			Str("aspect_type", event.AspectType).
			Msg("Webhook event not supported")
	}
	if qErr != nil {
		m.Logger.Error().
			Err(qErr).
			Str("aspect_type", event.AspectType).
			Int64("owner_id", event.OwnerID).
			Int64("activity_id", event.ObjectID).
			Msg("error enqueueing activity")
	}
}

func (m *Manager) newActivity(ctx context.Context, event webhooks.WebhookEvent) {
	var qErr error
	switch event.AspectType {
	case "create":
		// Set a low priority for webhooked events.
		priority := gue.JobPriorityDefault - 10000
		// Hugel potential is always there for new events. This is kinda unfortunate, but
		// the webhook gives us no intel into the event.
		qErr = m.EnqueueFetchActivity(ctx, database.ActivityDetailSourceWebhook, event.OwnerID, event.ObjectID, true, priority)
	case "update":
		qErr = m.EnqueueUpdateActivity(ctx, event)
	case "delete":
		m.Logger.Info().
			Interface("deleted", event.Updates).
			Msg("'Delete' webhook event to an activity")
	default:
		m.Logger.Warn().
			Str("aspect_type", event.AspectType).
			Msg("Webhook event not supported")
	}
	if qErr != nil {
		m.Logger.Error().
			Err(qErr).
			Str("aspect_type", event.AspectType).
			Int64("owner_id", event.OwnerID).
			Int64("activity_id", event.ObjectID).
			Msg("error enqueueing activity")
	}
}
