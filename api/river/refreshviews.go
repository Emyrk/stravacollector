package river

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
)

func (m *Manager) EnqueueRefreshViews(ctx context.Context, opts ...func(j *river.InsertOpts)) (bool, error) {
	iopts := &river.InsertOpts{}
	for _, opt := range opts {
		opt(iopts)
	}

	fi, err := m.cli.Insert(ctx, RefreshViewsArgs{}, iopts)

	skipped := false
	if fi != nil {
		skipped = fi.UniqueSkippedAsDuplicate
	}

	return !skipped, err
}

type RefreshViewsArgs struct {
}

func (RefreshViewsArgs) Kind() string { return "refresh_views" }
func (RefreshViewsArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue:    riverDatabaseQueue,
		Priority: PriorityHighest,
		UniqueOpts: river.UniqueOpts{
			ByArgs:   true,
			ByPeriod: time.Minute * 45,
		},
	}
}

type RefreshViewsWorker struct {
	mgr *Manager
	river.WorkerDefaults[RefreshViewsArgs]
}

func (*RefreshViewsWorker) Middleware(job *rivertype.JobRow) []rivertype.WorkerMiddleware {
	return []rivertype.WorkerMiddleware{}
}

func (w *RefreshViewsWorker) Work(ctx context.Context, job *river.Job[RefreshViewsArgs]) error {
	logger := jobLogFields(w.mgr.logger, job)

	wg := sync.WaitGroup{}
	start := time.Now()

	var hugelDone, hugel2023Done, superDone time.Duration
	var hugelErr, hugel2023Err, superErr, hugelLiteErr error

	wg.Add(1)
	go func() {
		hugelErr = w.mgr.db.RefreshHugelActivities(ctx)

		hugelLiteErr = w.mgr.db.RefreshHugelLiteActivities(ctx)
		hugelDone = time.Since(start)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		superErr = w.mgr.db.RefreshSuperHugelActivities(ctx)
		superDone = time.Since(start)
		wg.Done()
	}()

	wg.Add(1)
	go func() {
		hugel2023Err = w.mgr.db.RefreshHugel2023Activities(ctx)
		hugel2023Done = time.Since(start)
		wg.Done()
	}()

	wg.Wait()

	logger.Info().
		AnErr("super_err", superErr).
		AnErr("hugel_err", hugelErr).
		AnErr("hugel2023_err", hugel2023Err).
		AnErr("hugel_lite_err", hugelLiteErr).
		Str("super_duration", fmt.Sprintf("%.3fs", superDone.Seconds())).
		Str("hugel_duration", fmt.Sprintf("%.3fs", hugelDone.Seconds())).
		Str("hugel2023_duration", fmt.Sprintf("%.3fs", hugel2023Done.Seconds())).
		Msg("refresh views")

	_ = river.RecordOutput(ctx, map[string]any{
		"super_err":          superErr,
		"hugel_err":          hugelErr,
		"hugel2023_err":      hugel2023Err,
		"hugel_lite_err":     hugelLiteErr,
		"super_duration":     fmt.Sprintf("%.3fs", superDone.Seconds()),
		"hugel_duration":     fmt.Sprintf("%.3fs", hugelDone.Seconds()),
		"hugel2023_duration": fmt.Sprintf("%.3fs", hugel2023Done.Seconds()),
	})
	return nil
}
