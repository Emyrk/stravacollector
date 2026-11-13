package database

import (
	"context"
	"fmt"
	"strings"
)

type manualQuerier interface {
	YearlyHugelLeaderboard(ctx context.Context, arg YearlyHugelLeaderboardParams) ([]HugelLeaderboardRow, error)
}

type YearlyHugelLeaderboardParams struct {
	HugelLeaderboardParams
	RouteYear int
	Lite      bool
}

const hugel2026LiteScoringYear = 2025

func hugelActivitiesView(routeYear int, lite bool) string {
	view := "hugel_activities"
	if lite {
		view = "lite_hugel_activities"
	}

	switch routeYear {
	case 2023, 2024, 2025:
		return view + "_" + fmt.Sprint(routeYear)
	case 2026:
		// The confirmed 2026 Lite route retains the 2025 scoring segments.
		if lite {
			return view + "_" + fmt.Sprint(hugel2026LiteScoringYear)
		}
		return "hugel_activities_2026"
	default:
		return view
	}
}

func (q *sqlQuerier) YearlyHugelLeaderboard(ctx context.Context, arg YearlyHugelLeaderboardParams) ([]HugelLeaderboardRow, error) {
	query := strings.ReplaceAll(hugelLeaderboard, "hugel_activities", hugelActivitiesView(arg.RouteYear, arg.Lite))

	// Yearly windows are instants, independent of the PostgreSQL session timezone.
	// Keep the generated generic leaderboard's timestamp API unchanged.
	query = strings.ReplaceAll(query, ":: timestamp", ":: timestamptz")
	query = strings.ReplaceAll(query, "start_date <=", "start_date <")

	rows, err := q.db.Query(ctx, query, Timestamptz(arg.After.Time), Timestamptz(arg.Before.Time), arg.AthleteID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []HugelLeaderboardRow
	for rows.Next() {
		var i HugelLeaderboardRow
		if err := rows.Scan(
			&i.BestTime,
			&i.Rank,
			&i.ActivityID,
			&i.AthleteID,
			&i.TotalTimeSeconds,
			&i.Efforts,
			&i.Name,
			&i.DeviceWatts,
			&i.Distance,
			&i.MovingTime,
			&i.ElapsedTime,
			&i.TotalElevationGain,
			&i.StartDate,
			&i.AchievementCount,
			&i.AverageHeartrate,
			&i.AverageSpeed,
			&i.SufferScore,
			&i.AverageWatts,
			&i.AverageCadence,
			&i.Firstname,
			&i.Lastname,
			&i.Username,
			&i.ProfilePicLink,
			&i.Sex,
			&i.HugelCount,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
