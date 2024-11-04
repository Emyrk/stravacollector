package database

import (
	"context"
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

func (q *sqlQuerier) YearlyHugelLeaderboard(ctx context.Context, arg YearlyHugelLeaderboardParams) ([]HugelLeaderboardRow, error) {
	query := hugelLeaderboard

	if arg.RouteYear == 2023 {
		query = strings.Replace(query, "hugel_activities", "hugel_activities_2023", -1)
	}

	if arg.Lite {
		query = strings.Replace(query, "hugel_activities", "lite_hugel_activities", -1)
	}

	rows, err := q.db.QueryContext(ctx, query, arg.AthleteID, arg.After, arg.Before)
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
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
