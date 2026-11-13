package database_test

import (
	"context"
	"testing"
	"time"

	"github.com/Emyrk/strava/database"
	"github.com/Emyrk/strava/database/postgres"
	"github.com/Emyrk/strava/internal/hugeldate"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

func TestYearlyHugelLeaderboardWindow(t *testing.T) {
	t.Run("event", func(t *testing.T) {
		testYearlyWindow(t, hugeldate.Year2026.Start, hugeldate.Year2026.End)
	})
	t.Run("daylight-saving-transition", func(t *testing.T) {
		central := hugeldate.CentralTimeZone
		testYearlyWindow(t, time.Date(2026, time.November, 1, 0, 0, 0, 0, central), time.Date(2026, time.December, 1, 0, 0, 0, 0, central))
	})
}

func testYearlyWindow(t *testing.T, start, end time.Time) {
	t.Helper()
	url, cleanup, err := postgres.Open()
	require.NoError(t, err)
	t.Cleanup(cleanup)
	ctx := context.Background()
	config, err := pgxpool.ParseConfig(url)
	require.NoError(t, err)
	config.MaxConns = 1
	conn, err := pgxpool.NewWithConfig(ctx, config)
	require.NoError(t, err)
	t.Cleanup(conn.Close)
	exec := func(sql string, args ...any) {
		t.Helper()
		_, err := conn.Exec(ctx, sql, args...)
		require.NoError(t, err)
	}
	// Temporary relations isolate filtering from route aggregation and Strava metadata.
	// Each finish belongs to a different athlete so DISTINCT ON cannot hide a boundary.
	exec(`CREATE TEMP TABLE window_finishes (id bigint, start_date timestamptz)`)
	for i, instant := range []time.Time{
		start.Add(-time.Microsecond), start, start.Add(time.Microsecond),
		end.Add(-time.Microsecond), end, end.Add(time.Microsecond),
	} {
		exec(`INSERT INTO window_finishes VALUES ($1, $2)`, i+1, instant)
	}
	exec(`CREATE TEMP TABLE hugel_activities_2026 AS
 SELECT id AS activity_id, id AS athlete_id, ARRAY[]::bigint[] AS segment_ids,
 id AS total_time_seconds, '[]'::jsonb AS efforts FROM window_finishes;
 CREATE TEMP TABLE activity_summary AS
 SELECT id, start_date, ''::text AS name, false AS device_watts,
 0::float8 AS distance, 0::float8 AS moving_time, 0::float8 AS elapsed_time,
 0::float8 AS total_elevation_gain, 0::int AS achievement_count,
 0::float8 AS average_heartrate, 0::float8 AS average_speed FROM window_finishes;
 CREATE TEMP TABLE activity_detail AS
 SELECT id, 0::int AS suffer_score, 0::float8 AS average_watts,
 0::float8 AS average_cadence FROM window_finishes;
 CREATE TEMP TABLE athletes AS
 SELECT id, ''::text AS firstname, ''::text AS lastname, ''::text AS username,
 ''::text AS profile_pic_link, ''::text AS sex FROM window_finishes;
 CREATE TEMP TABLE athlete_hugel_count AS
 SELECT id AS athlete_id, 1::bigint AS count FROM window_finishes`)
	db := database.New(conn)
	for _, zone := range []string{"UTC", "America/Chicago", "Asia/Tokyo"} {
		t.Run(zone, func(t *testing.T) {
			// One connection ensures SET and the leaderboard query use the same session.
			_, err := conn.Exec(ctx, `SELECT set_config('TimeZone', $1, false)`, zone)
			require.NoError(t, err)
			rows, err := db.YearlyHugelLeaderboard(ctx, database.YearlyHugelLeaderboardParams{
				RouteYear: 2026,
				HugelLeaderboardParams: database.HugelLeaderboardParams{
					After: database.Timestamp(start), Before: database.Timestamp(end), AthleteID: int64(0),
				},
			})
			require.NoError(t, err)
			var ids []int64
			for _, row := range rows {
				ids = append(ids, row.ActivityID)
				require.Equal(t, int64(2), row.BestTime)
			}
			require.Equal(t, []int64{2, 3, 4}, ids)
		})
	}
}
