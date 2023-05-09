package queue

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/vgarvardt/gue/v5"

	"github.com/Emyrk/strava/api/webhooks"
	"github.com/Emyrk/strava/database"
)

func (m *Manager) EnqueueUpdateActivity(ctx context.Context, event webhooks.WebhookEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}

	return m.Client.Enqueue(ctx, &gue.Job{
		Type:  updateActivityField,
		Queue: stravaUpdateHookQueue,
		Args:  data,
	})
}

func (m *Manager) EnqueueDeleteActivity(ctx context.Context, event webhooks.WebhookEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}

	return m.Client.Enqueue(ctx, &gue.Job{
		Type:  deleteActivityJob,
		Queue: stravaUpdateHookQueue,
		Args:  data,
	})
}

func (m *Manager) deleteActivity(ctx context.Context, j *gue.Job) error {
	logger := jobLogFields(m.Logger, j)

	var args webhooks.WebhookEvent
	err := json.Unmarshal(j.Args, &args)
	if err != nil {
		logger.Error().Err(err).Msg("json unmarshal, update activity job abandoned")
		return nil
	}

	_, err = m.DB.DeleteActivity(ctx, args.ObjectID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	}
	return fmt.Errorf("delete activity: %w", err)
}

func (m *Manager) updateActivity(ctx context.Context, j *gue.Job) error {
	logger := jobLogFields(m.Logger, j)

	var args webhooks.WebhookEvent
	err := json.Unmarshal(j.Args, &args)
	if err != nil {
		logger.Error().Err(err).Msg("json unmarshal, update activity job abandoned")
		return nil
	}

	// This updates an activity.
	err = m.DB.InTx(func(store database.Store) error {
		_, err := store.GetActivitySummary(ctx, args.ObjectID)
		if errors.Is(err, sql.ErrNoRows) {
			logger.Warn().Err(err).Msg("activity not found, update activity job abandoned")
			return nil
		}
		if err != nil {
			return fmt.Errorf("get activity: %w", err)
		}
		//// The update is older than the activity fetch, ignore it.
		//if act.UpdatedAt.Unix() > args.EventTime {
		//	return nil
		//}

		for k, v := range args.Updates {
			switch k {
			case "title":
				err := store.UpdateActivityName(ctx, database.UpdateActivityNameParams{
					ID:   args.ObjectID,
					Name: v,
				})
				if err != nil {
					return fmt.Errorf("update activity name: %w", err)
				}
			default:
				return fmt.Errorf("unknown update type: %s", k)
			}
		}
		return nil
	}, nil)
	if err != nil {
		return fmt.Errorf("update activity: %w", err)
	}

	return nil
}
