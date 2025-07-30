package river

import (
	"context"
	"fmt"
	"time"

	"github.com/Emyrk/strava/database"
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

	edds := EddingtonNumbers{}
	for _, act := range acts {
		edds.Add(int(database.DistanceToMiles(act.Distance)))
	}

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

	return nil
}

// EddingtonNumbers is a slice of integers representing the number of rides over the
// index in miles in distance.
type EddingtonNumbers []int32

func (e EddingtonNumbers) Current() int32 {
	for need, have := range e {
		need = need + 1 // 1-indexed
		if int32(need) > have {
			return int32(need) - 1
		}
	}
	return e[len(e)-1]
}

func (e *EddingtonNumbers) Add(value int) {
	if value < 0 {
		return
	}
	if *e == nil {
		*e = make(EddingtonNumbers, 0, value)
	}
	if value > len(*e) {
		*e = append(*e, make(EddingtonNumbers, value-len(*e))...)
	}
	for i := 0; i < value; i++ {
		(*e)[i]++
	}
}
