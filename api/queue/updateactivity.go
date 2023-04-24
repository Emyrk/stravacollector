package queue

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/Emyrk/strava/database"

	"github.com/Emyrk/strava/api/webhooks"
	"github.com/vgarvardt/gue/v5"
)

func (m *Manager) EnqueueUpdateActivity(ctx context.Context, event webhooks.WebhookEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}

	return m.Client.Enqueue(ctx, &gue.Job{
		Type:  updateActivityJob,
		Queue: stravaUpdateActivityQueue,
		Args:  data,
	})
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
		_, err := store.GetActivity(ctx, args.ObjectID)
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
					ID:   0,
					Name: "",
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
