package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/Emyrk/strava/api/river"
	"github.com/Emyrk/strava/database"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
)

func refetchActivities(ctx context.Context, db database.Store, dbURL string, logger zerolog.Logger) error {
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

	acts, err := db.GetActivitySummariesByDate(ctx, pgtype.Timestamptz{
		Time:             time.Now().AddDate(0, 0, -2),
		InfinityModifier: 0,
		Valid:            true,
	})

	if err != nil {
		return fmt.Errorf("get activity summaries: %w", err)
	}

	for i, act := range acts {
		_, err = riverManager.EnqueueFetchActivity(ctx, river.FetchActivityArgs{
			Source:         database.ActivityDetailSourceManual,
			ActivityID:     act.ID,
			AthleteID:      act.AthleteID,
			HugelPotential: true,
			OnHugelDates:   true,
		}, river.PriorityHighest)
		if err != nil {
			return fmt.Errorf("enqueue fetch activity job: %w", err)
		}
		fmt.Print(".")
		if i%10 == 0 {
			fmt.Printf(" %d/%d\n", i, len(acts))
		}
	}
	return nil
}
