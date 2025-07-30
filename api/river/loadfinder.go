package river

import (
	"context"
	"fmt"
	"time"

	"github.com/riverqueue/river"
)

func (m *Manager) EnqueueLoadFinder(until time.Time, opts ...func(j *river.InsertOpts)) (bool, error) {
	iopts := &river.InsertOpts{
		ScheduledAt: until,
	}
	for _, opt := range opts {
		opt(iopts)
	}

	fi, err := m.cli.Insert(m.appCtx, LoadFinderArgs{}, iopts)

	skipped := false
	if fi != nil {
		skipped = fi.UniqueSkippedAsDuplicate
	}

	return !skipped, err
}

type LoadFinderArgs struct {
}

func (LoadFinderArgs) Kind() string { return "load_finder" }
func (LoadFinderArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		MaxAttempts: 3,
		Queue:       riverDatabaseQueue,
		UniqueOpts: river.UniqueOpts{
			ByArgs:   true,
			ByPeriod: time.Minute * 30,
		},
	}
}

type LoadFinderWorker struct {
	mgr *Manager
	river.WorkerDefaults[LoadFinderArgs]
}

func (w *LoadFinderWorker) Work(ctx context.Context, job *river.Job[LoadFinderArgs]) error {
	need, err := w.mgr.db.GetAthleteNeedsForwardLoad(ctx)
	if err != nil {
		return err
	}

	out := make(map[string]string)
	for _, athlete := range need {
		_, err := w.mgr.EnqueueForwardLoad(ctx, athlete.AthleteForwardLoad.AthleteID)
		if err != nil {
			out[fmt.Sprintf("https://www.strava.com/athletes/%d", athlete.AthleteLogin.AthleteID)] = err.Error()
		}
	}

	out["quantity"] = fmt.Sprintf("%d athletes included", len(need))
	_ = river.RecordOutput(ctx, out)
	return nil
}
