package river

import (
	"context"

	"github.com/riverqueue/river"
)

type AthleteLoadArgs struct {
}

func (AthleteLoadArgs) Kind() string { return "athlete_historical_load" }

type AthleteLoadWorker struct {
	mgr *Manager
	river.WorkerDefaults[AthleteLoadArgs]
}

func (w *AthleteLoadWorker) Work(ctx context.Context, job *river.Job[AthleteLoadArgs]) error {
	var _ = w.mgr

	return nil
}
