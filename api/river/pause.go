package river

import (
	"context"
	"fmt"
	"time"

	"github.com/riverqueue/river"
)

func (mgr *Manager) Pause(until time.Time, reason, queue string) error {
	q, err := mgr.cli.QueueGet(mgr.appCtx, queue)
	if err == nil && q.PausedAt != nil {
		// Already paused, no need to pause it again.
		return nil
	}

	err = mgr.cli.QueuePause(mgr.appCtx, queue, &river.QueuePauseOpts{})
	if err != nil {
		return fmt.Errorf("could not pause queue: %w", err)
	}

	_, err = mgr.EnqueueResume(until, reason, queue)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) EnqueueResume(until time.Time, reason, queue string, opts ...func(j *river.InsertOpts)) (bool, error) {
	iopts := &river.InsertOpts{
		ScheduledAt: until,
	}
	for _, opt := range opts {
		opt(iopts)
	}

	fi, err := m.cli.Insert(m.appCtx, ResumeArgs{
		Queues: []string{queue},
		// 10s debounce to prevent dupe hits on scheduled jobs.
		RandomID: until.Unix() / 10,
		Reason:   reason,
	}, iopts)

	skipped := false
	if fi != nil {
		skipped = fi.UniqueSkippedAsDuplicate
	}

	return !skipped, err
}

type ResumeArgs struct {
	Queues []string `json:"queue"`
	// RandomID is used to prevent dupe hits on scheduled jobs.
	// The cron jobs will want to be ignored when a dupe is hit.
	RandomID int64  `json:"random_id,omitzero"`
	Reason   string `json:"reason,omitempty"`
}

func (ResumeArgs) Kind() string { return "resume" }
func (ResumeArgs) InsertOpts() river.InsertOpts {
	return river.InsertOpts{
		Queue: riverControlQueue,
		UniqueOpts: river.UniqueOpts{
			ByArgs:   true,
			ByPeriod: time.Minute * 30,
		},
	}
}

type ResumeWorker struct {
	mgr *Manager
	river.WorkerDefaults[ResumeArgs]
}

func (w *ResumeWorker) Work(ctx context.Context, job *river.Job[ResumeArgs]) error {
	out := make(map[string]string)
	for _, q := range job.Args.Queues {
		queue, err := w.mgr.cli.QueueGet(ctx, q)
		if err != nil {
			_ = river.RecordOutput(ctx, fmt.Sprintf("could not get queue %q: %s", q, err.Error()))
			return nil
		}
		out[q] = fmt.Sprintf("Queue paused = %t", queue.PausedAt != nil)
		_ = w.mgr.cli.QueueResume(ctx, q, &river.QueuePauseOpts{})
	}

	_ = river.RecordOutput(ctx, out)
	return nil
}
