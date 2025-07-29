package river

import (
	"context"
	"fmt"
	"time"

	"github.com/riverqueue/river"
)

func (mgr *Manager) Pause(until time.Time, queue string) error {
	err := mgr.cli.QueuePause(mgr.appCtx, queue, &river.QueuePauseOpts{})
	if err != nil {
		return fmt.Errorf("could not pause queue: %w", err)
	}

	_, err = mgr.EnqueueResume(until, queue)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) EnqueueResume(until time.Time, queue string, opts ...func(j *river.InsertOpts)) (bool, error) {
	iopts := &river.InsertOpts{
		ScheduledAt: until,
	}
	for _, opt := range opts {
		opt(iopts)
	}

	fi, err := m.cli.Insert(m.appCtx, ResumeArgs{
		Queue: queue,
	}, iopts)

	skipped := false
	if fi != nil {
		skipped = fi.UniqueSkippedAsDuplicate
	}

	return !skipped, err
}

type ResumeArgs struct {
	Queue string `json:"queue"`
}

func (ResumeArgs) Kind() string { return "resume" }
func (ResumeArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue:      riverControlQueue,
		UniqueOpts: river.UniqueOpts{},
	}
}

type ResumeWorker struct {
	mgr *Manager
	river.WorkerDefaults[ResumeArgs]
}

func (w *ResumeWorker) Work(ctx context.Context, job *river.Job[ResumeArgs]) error {
	q := job.Args.Queue
	queue, err := w.mgr.cli.QueueGet(ctx, q)
	if err != nil {
		_ = river.RecordOutput(ctx, fmt.Sprintf("could not get queue %q: %s", q, err.Error()))
		return nil
	}

	_ = river.RecordOutput(ctx, fmt.Sprintf("Queue paused = %t", queue.PausedAt != nil))
	return w.mgr.cli.QueueResume(ctx, q, &river.QueuePauseOpts{})
}
