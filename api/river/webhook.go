package river

import (
	"context"
	"time"

	"github.com/Emyrk/strava/api/webhooks"
	"github.com/Emyrk/strava/database"
	"github.com/riverqueue/river"
)

func (m *Manager) HandleWebhookEvents(ctx context.Context, c <-chan *webhooks.Handled[webhooks.WebhookEvent]) {
	for {
		select {
		case <-ctx.Done():
			return
		case h := <-c:
			if h == nil {
				continue
			}
			event := h.Data
			switch event.ObjectType {
			case "activity":
				m.newActivity(ctx, event)
			case "athlete":
				m.newAthlete(ctx, event)
			default:
				m.logger.Warn().
					Str("object_type", event.ObjectType).
					Msg("Webhook event not supported")
			}

			// Tell strava that we have processed the event.
			h.MarkDone()
		}
	}
}

func (m *Manager) newAthlete(ctx context.Context, event webhooks.WebhookEvent) {
	var qErr error
	switch event.AspectType {
	case "create":
		m.logger.Warn().
			Interface("event", event).
			Msg("Webhook create event to an athlete not handled")
	case "update":
		_, qErr = m.EnqueueUpdateAthlete(ctx, event)
	case "delete":
		m.logger.Warn().
			Interface("event", event).
			Msg("Webhook delete event to an athlete not handled")
	default:
		m.logger.Warn().
			Str("aspect_type", event.AspectType).
			Msg("Webhook event not supported")
	}
	if qErr != nil {
		m.logger.Error().
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
		priority := PriorityLow
		// Hugel potential is always there for new events. This is kinda unfortunate, but
		// the webhook gives us no intel into the event.
		_, qErr = m.EnqueueFetchActivity(ctx, database.ActivityDetailSourceWebhook, event.OwnerID, event.ObjectID, true, true, priority, func(j *river.InsertOpts) {
			// When syncing a new activity, let strava first load all the segments. Strava populates segments async,
			// and if we query them too soon, we get an empty list.
			j.ScheduledAt = time.Now().Add(time.Minute * 30)
		})
	case "update":
		_, qErr = m.EnqueueUpdateActivity(ctx, event)
	case "delete":
		m.logger.Info().
			Interface("deleted", event.Updates).
			Msg("'Delete' webhook event to an activity")
	default:
		m.logger.Warn().
			Str("aspect_type", event.AspectType).
			Msg("Webhook event not supported")
	}
	if qErr != nil {
		m.logger.Error().
			Err(qErr).
			Str("aspect_type", event.AspectType).
			Int64("owner_id", event.OwnerID).
			Int64("activity_id", event.ObjectID).
			Msg("error enqueueing activity")
	}
}
