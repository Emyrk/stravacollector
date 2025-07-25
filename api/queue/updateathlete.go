package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/vgarvardt/gue/v5"

	"github.com/Emyrk/strava/api/webhooks"
	"github.com/Emyrk/strava/database"
)

func (m *Manager) EnqueueUpdateAthlete(ctx context.Context, event *webhooks.WebhookEvent) error {
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("json marshal: %w", err)
	}

	return m.Client.Enqueue(ctx, &gue.Job{
		Type:  updateAthleteJob,
		Queue: stravaUpdateHookQueue,
		Args:  data,
	})
}

func (m *Manager) updateAthlete(ctx context.Context, j *gue.Job) error {
	logger := jobLogFields(m.Logger, j)

	var args webhooks.WebhookEvent
	err := json.Unmarshal(j.Args, &args)
	if err != nil {
		logger.Error().Err(err).Msg("json unmarshal, update activity job abandoned")
		return nil
	}

	if args.ObjectType != "athlete" {
		return nil
	}

	err = m.DB.InTx(func(store database.Store) error {
		for key, v := range args.Updates {
			switch key {
			case "authorized":
				if authed, err := strconv.ParseBool(v); !authed && err == nil {
					err := m.DB.DeleteAthleteLogin(ctx, args.ObjectID)
					if err != nil {
						return fmt.Errorf("delete athlete login: %w", err)
					}
				}
			default:
				logger.Error().
					Int64("athelete_id", args.ObjectID).
					Str("key", key).
					Str("value", v).
					Msg("unknown athlete update")
				return fmt.Errorf("unknown athlete update: %s", key)
			}
		}
		return nil
	}, nil)
	if err != nil {
		return fmt.Errorf("update athlete tx: %w", err)
	}

	return nil
}
