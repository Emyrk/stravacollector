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
			ByPeriod: time.Minute * 25,
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
	skipped := 0
	failed := 0
	for _, athlete := range need {
		skip, err := w.mgr.EnqueueForwardLoad(ctx, athlete.AthleteForwardLoad.AthleteID)
		if err != nil {
			out[fmt.Sprintf("https://www.strava.com/athletes/%d", athlete.AthleteLogin.AthleteID)] = err.Error()
			failed++
		}
		if err == nil && skip {
			skipped++
		}
	}

	out["quantity"] = fmt.Sprintf("%d athletes included", len(need))
	out["actual"] = fmt.Sprintf("%d load jobs started", len(need)-failed-skipped)
	out["skipped"] = fmt.Sprintf("%d athletes skipped", skipped)
	out["failed"] = fmt.Sprintf("%d athletes failed", failed)
	_ = river.RecordOutput(ctx, out)
	return nil
}
