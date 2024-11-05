package cli

import (
	"context"
	"fmt"

	"github.com/Emyrk/strava/api/queue"
	"github.com/Emyrk/strava/database"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
	"github.com/vgarvardt/gue/v5"
)

func redownloadHugels(ctx context.Context, db database.Store, dbURL string, logger zerolog.Logger) error {
	var allHugels []database.HugelLeaderboardRow

	m, err := queue.New(ctx, queue.Options{
		DBURL:    dbURL,
		Logger:   logger,
		DB:       db,
		OAuthCfg: nil,
		Registry: prometheus.NewRegistry(),
	})
	if err != nil {
		return fmt.Errorf("failed to create queue manager: %w", err)
	}

	for _, year := range []int{2023, 2024} {
		for _, lite := range []bool{false, true} {
			if year == 2023 && lite {
				continue
			}
			hugels, err := db.YearlyHugelLeaderboard(ctx, database.YearlyHugelLeaderboardParams{
				HugelLeaderboardParams: database.HugelLeaderboardParams{},
				RouteYear:              year,
				Lite:                   lite,
			})
			if err != nil {
				return fmt.Errorf("failed to fetch %d activities: %w", year, err)
			}

			allHugels = append(allHugels, hugels...)
		}
		hugels, err := db.YearlyHugelLeaderboard(ctx, database.YearlyHugelLeaderboardParams{
			HugelLeaderboardParams: database.HugelLeaderboardParams{},
			RouteYear:              year,
			Lite:                   false,
		})
		if err != nil {
			return fmt.Errorf("failed to refresh %d activities: %w", year, err)
		}
		allHugels = append(allHugels, hugels...)
	}

	fmt.Println(len(allHugels))

	for _, h := range allHugels {
		err = m.EnqueueFetchActivity(ctx, database.ActivityDetailSourceManual, h.AthleteID, h.ActivityID, true, true, gue.JobPriorityHighest)
		if err != nil {
			return fmt.Errorf("failed to enqueue fetch activity job: %w", err)
		}
	}

	return nil
}
