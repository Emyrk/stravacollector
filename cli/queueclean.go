package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/Emyrk/strava/api/river"
	"github.com/Emyrk/strava/database"
	"github.com/prometheus/client_golang/prometheus"
	river2 "github.com/riverqueue/river"
	"github.com/rs/zerolog"
)

func queueClean(ctx context.Context, db database.Store, dbURL string, logger zerolog.Logger) error {
	riverManager, err := river.New(ctx, river.Options{
		DBURL:      dbURL,
		Logger:     logger.With().Str("component", "river").Logger(),
		DB:         db,
		Registry:   prometheus.NewRegistry(),
		InsertOnly: true,
	})
	if err != nil {
		return fmt.Errorf("create river manager: %w", err)
	}
	defer riverManager.Close(ctx)

	var total int
	var cursor *river2.JobListCursor
	for {
		fmt.Println("Deleting jobs...")
		params := river2.NewJobListParams().
			Kinds("resume")
		if cursor != nil {
			params = params.After(cursor)
		}
		jobs, err := riverManager.Cli().JobList(
			ctx,
			params,
		)
		if err != nil {
			return fmt.Errorf("list jobs: %w", err)
		}

		fmt.Printf("Found %d jobs to delete\n", len(jobs.Jobs))
		cursor = jobs.LastCursor
		start := time.Now()
		for _, job := range jobs.Jobs {
			fmt.Print(".")
			_, err := riverManager.Cli().JobDelete(ctx, job.ID)
			if err != nil {
				return fmt.Errorf("delete job %d: %w", job.ID, err)
			}
			total++
		}
		end := time.Since(start)
		fmt.Println("")
		fmt.Printf("Deleted %d more jobs in %.2fs (%.2fs/job), %d total\n", len(jobs.Jobs), end.Seconds(), end.Seconds()/float64(len(jobs.Jobs)), total)
		if len(jobs.Jobs) == 0 {
			break
		}
	}

	return nil
}
