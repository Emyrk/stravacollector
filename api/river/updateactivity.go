package river

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Emyrk/strava/api/webhooks"
	"github.com/Emyrk/strava/database"
	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
)

const (
	updateActivityField = "update_activity"
	deleteActivityJob   = "delete_activity"
	updateAthleteJob    = "update_athlete"
)

func (m *Manager) enqueueHook(ctx context.Context, typ string, event webhooks.WebhookEvent, opts ...func(j *river.InsertOpts)) (bool, error) {
	iopts := &river.InsertOpts{
		Priority: PriorityDefault,
	}
	for _, opt := range opts {
		opt(iopts)
	}

	fi, err := m.cli.Insert(ctx, UpdateActivityArgs{
		Type:         typ,
		WebhookEvent: event,
	}, iopts)

	skipped := false
	if fi != nil {
		skipped = fi.UniqueSkippedAsDuplicate
	}

	return !skipped, err
}

func (m *Manager) EnqueueUpdateActivity(ctx context.Context, event webhooks.WebhookEvent, opts ...func(j *river.InsertOpts)) (bool, error) {
	return m.enqueueHook(ctx, updateActivityField, event, opts...)
}

func (m *Manager) EnqueueDeleteActivity(ctx context.Context, event webhooks.WebhookEvent, opts ...func(j *river.InsertOpts)) (bool, error) {
	return m.enqueueHook(ctx, deleteActivityJob, event, opts...)
}

func (m *Manager) EnqueueUpdateAthlete(ctx context.Context, event webhooks.WebhookEvent, opts ...func(j *river.InsertOpts)) (bool, error) {
	return m.enqueueHook(ctx, updateAthleteJob, event, opts...)
}

type UpdateActivityArgs struct {
	Type string
	webhooks.WebhookEvent
}

func (UpdateActivityArgs) Kind() string { return "update_activity" }
func (UpdateActivityArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue:       riverDatabaseQueue,
		MaxAttempts: 3,
		UniqueOpts: river.UniqueOpts{
			ByArgs:   true,
			ByPeriod: time.Minute * 5,
		},
	}
}

type UpdateActivityWorker struct {
	mgr *Manager
	river.WorkerDefaults[UpdateActivityArgs]
}

func (*UpdateActivityWorker) Middleware(job *rivertype.JobRow) []rivertype.WorkerMiddleware {
	return []rivertype.WorkerMiddleware{}
}

func (w *UpdateActivityWorker) Work(ctx context.Context, job *river.Job[UpdateActivityArgs]) error {
	switch job.Args.Type {
	case updateActivityField:
		return w.Update(ctx, job)
	case deleteActivityJob:
		return w.Delete(ctx, job)
	case updateAthleteJob:
		return w.UpdateAthlete(ctx, job)
	default:
		return fmt.Errorf("unknown job type: %s", job.Args.Type)
	}
}

func (w *UpdateActivityWorker) Update(ctx context.Context, job *river.Job[UpdateActivityArgs]) error {
	logger := jobLogFields[UpdateActivityArgs](w.mgr.logger, job)
	args := job.Args

	// This updates an activity.
	err := w.mgr.db.InTx(func(store database.Store) error {
		_, err := store.GetActivitySummary(ctx, args.ObjectID)
		if errors.Is(err, sql.ErrNoRows) {
			logger.Warn().
				Str("activity_id", fmt.Sprintf("%d", args.ObjectID)).
				Err(err).
				Msg("activity not found, update activity job abandoned")
			_ = river.RecordOutput(ctx, "activity not found, nothing to update")
			return nil
		}
		if err != nil {
			return fmt.Errorf("get activity: %w", err)
		}

		//// The update is older than the activity fetch, ignore it.
		// if act.UpdatedAt.Unix() > args.EventTime {
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
			case "type":
				err := store.UpdateActivityType(ctx, database.UpdateActivityTypeParams{
					ID:   args.ObjectID,
					Type: v,
				})
				if err != nil {
					return fmt.Errorf("update activity type: %w", err)
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

func (w *UpdateActivityWorker) Delete(ctx context.Context, job *river.Job[UpdateActivityArgs]) error {
	args := job.Args

	_, err := w.mgr.db.DeleteActivity(ctx, args.ObjectID)
	if errors.Is(err, sql.ErrNoRows) {
		_ = river.RecordOutput(ctx, "activity not found, nothing to delete")
		return nil
	}
	if err != nil {
		return err
	}
	return nil
}

func (w *UpdateActivityWorker) UpdateAthlete(ctx context.Context, job *river.Job[UpdateActivityArgs]) error {
	logger := jobLogFields[UpdateActivityArgs](w.mgr.logger, job)
	args := job.Args

	if args.ObjectType != "athlete" {
		_ = river.RecordOutput(ctx, "not an athlete update, skipping")
		return nil
	}

	err := w.mgr.db.InTx(func(store database.Store) error {
		for key, v := range args.Updates {
			switch key {
			case "authorized":
				if authed, err := strconv.ParseBool(v); !authed && err == nil {
					err := store.DeleteAthleteLogin(ctx, args.ObjectID)
					if err != nil && !errors.Is(err, sql.ErrNoRows) {
						return fmt.Errorf("delete athlete login: %w", err)
					}
				}
			default:
				logger.Error().
					Int64("athelete_id", args.ObjectID).
					Str("key", key).
					Str("value", v).
					Msg("unknown athlete update")
				_ = river.RecordOutput(ctx, fmt.Sprintf("unknown athlete update: %s", key))
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
