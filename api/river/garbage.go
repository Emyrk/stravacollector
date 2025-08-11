package river

import (
	"context"
	"fmt"
	"time"

	"github.com/riverqueue/river"
	"github.com/riverqueue/river/rivertype"
)

const (
	timeout = time.Minute * 5
)

func (m *Manager) EnqueueGarbageCollect(opts ...func(j *river.InsertOpts)) (bool, error) {
	iopts := &river.InsertOpts{}
	for _, opt := range opts {
		opt(iopts)
	}

	fi, err := m.cli.Insert(m.appCtx, GarbageCollectArgs{}, iopts)

	skipped := false
	if fi != nil {
		skipped = fi.UniqueSkippedAsDuplicate
	}

	return !skipped, err
}

type GarbageCollectArgs struct {
}

func (GarbageCollectArgs) Kind() string { return "GarbageCollect" }
func (GarbageCollectArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue:       riverDatabaseQueue,
		MaxAttempts: 2,

		UniqueOpts: river.UniqueOpts{
			ByArgs:   true,
			ByPeriod: time.Hour * 23, // Run once a day
		},
	}
}

type GarbageCollectWorker struct {
	mgr *Manager
	river.WorkerDefaults[GarbageCollectArgs]
}

func (*GarbageCollectWorker) Middleware(job *rivertype.JobRow) []rivertype.WorkerMiddleware {
	return []rivertype.WorkerMiddleware{}
}

func (w GarbageCollectWorker) Timeout(*river.Job[GarbageCollectArgs]) time.Duration {
	return timeout
}

func (w *GarbageCollectWorker) Work(ctx context.Context, _ *river.Job[GarbageCollectArgs]) error {
	var total int
	var cursor *river.JobListCursor
	stat := time.Now()
	deleteBefore := time.Now().Add(time.Hour * 24 * -1)

GarbageCollectLoop:
	for {
		if time.Since(stat) > timeout-time.Minute {
			// Exit if we've been running for too long.
			// Otherwise, this just gets cancelled by the river manager.
			break
		}
		params := river.NewJobListParams().
			Kinds(
				ResumeArgs{}.Kind(),
				ReloadSegmentsArgs{}.Kind(),
			).
			//States(rivertype.JobStateCompleted).
			OrderBy(river.JobListOrderByFinalizedAt, river.SortOrderAsc)
		if cursor != nil {
			params = params.After(cursor)
		}
		jobs, err := w.mgr.cli.JobList(
			ctx,
			params,
		)
		if err != nil {
			return fmt.Errorf("list jobs: %w", err)
		}

		cursor = jobs.LastCursor
		for _, job := range jobs.Jobs {
			if job.FinalizedAt != nil {
				break GarbageCollectLoop
			}

			if job.FinalizedAt.After(deleteBefore) {
				break GarbageCollectLoop
			}

			_, err := w.mgr.cli.JobDelete(ctx, job.ID)
			if err != nil {
				return fmt.Errorf("delete job %d: %w", job.ID, err)
			}
			total++
		}

		if len(jobs.Jobs) == 0 {
			break
		}
	}

	_ = river.RecordOutput(ctx, map[string]interface{}{
		"total": total,
	})
	return nil
}
