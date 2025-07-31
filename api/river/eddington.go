package river

import (
	"context"
	"fmt"
	"time"

	"github.com/Emyrk/strava/database"
	"github.com/Emyrk/strava/internal/eddington"
	"github.com/riverqueue/river"
)

func (m *Manager) EnqueueEddington(athlete int64, opts ...func(j *river.InsertOpts)) (bool, error) {
	iopts := &river.InsertOpts{}
	for _, opt := range opts {
		opt(iopts)
	}

	fi, err := m.cli.Insert(m.appCtx, EddingtonArgs{
		AthleteID: athlete,
	}, iopts)

	skipped := false
	if fi != nil {
		skipped = fi.UniqueSkippedAsDuplicate
	}

	return !skipped, err
}

type EddingtonArgs struct {
	AthleteID int64
}

func (EddingtonArgs) Kind() string { return "eddington" }
func (EddingtonArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue:       riverDatabaseQueue,
		MaxAttempts: 3,
		UniqueOpts: river.UniqueOpts{
			ByArgs:   true,
			ByPeriod: time.Minute * 1,
		},
	}
}

type EddingtonWorker struct {
	mgr *Manager
	river.WorkerDefaults[EddingtonArgs]
}

func (w *EddingtonWorker) Work(ctx context.Context, job *river.Job[EddingtonArgs]) error {
	acts, err := w.mgr.db.EddingtonActivities(ctx, job.Args.AthleteID)
	if err != nil {
		return fmt.Errorf("fetching activities: %w", err)
	}

	edds := eddington.FromActivities(acts)

	_, err = w.mgr.db.UpsertAthleteEddington(ctx, database.UpsertAthleteEddingtonParams{
		AthleteID:        job.Args.AthleteID,
		MilesHistogram:   edds,
		CurrentEddington: edds.Current(),
		LastCalculated:   database.Timestamptz(time.Now()),
		TotalActivities:  int32(len(acts)),
	})
	if err != nil {
		return fmt.Errorf("upserting athlete eddington: %w", err)
	}

	_ = river.RecordOutput(ctx, map[string]interface{}{
		"number": edds.Current(),
		"total":  len(acts),
	})
	return nil
}
