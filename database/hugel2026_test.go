package database_test

import (
	"context"
	"os"
	"testing"

	"github.com/Emyrk/strava/database"
	"github.com/Emyrk/strava/database/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

// This test uses an isolated PostgreSQL database, never a shared application DB.
func TestHugel2026Migration(t *testing.T) {
	url, cleanup, err := postgres.Open()
	require.NoError(t, err)
	t.Cleanup(cleanup)
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, url)
	require.NoError(t, err)
	t.Cleanup(pool.Close)
	exec := func(query string, args ...any) {
		t.Helper()
		_, err := pool.Exec(ctx, query, args...)
		require.NoError(t, err)
	}
	migration := func(direction string) {
		t.Helper()
		sql, err := os.ReadFile("migrations/000042_hugel_2026." + direction + ".sql")
		require.NoError(t, err)
		exec(string(sql))
	}
	migration("down")
	old := []int64{629046, 405437, 33454575, 28704989, 995700, 6744304}
	want := []int64{629046, 811022, 42516497, 6744304, 3145259}
	exec(`UPDATE competitive_routes SET segments=$1 WHERE name='das-hugel'`, old)
	// Lite is an independently maintained route. The migration must not change it.
	exec(`INSERT INTO competitive_routes VALUES ('lite-das-hugel', 'Lite', 'Unchanged for 2026', ARRAY[629046]::bigint[])
       ON CONFLICT (name) DO UPDATE SET segments=EXCLUDED.segments`)
	var oldLiteView string
	require.NoError(t, pool.QueryRow(ctx, `SELECT pg_get_viewdef('lite_hugel_activities_2025')`).Scan(&oldLiteView))
	migration("up")
	routeSegments := func(name string) []int64 {
		t.Helper()
		var got []int64
		require.NoError(t, pool.QueryRow(ctx, `SELECT segments FROM competitive_routes WHERE name=$1`, name).Scan(&got))
		return got
	}
	require.Equal(t, old, routeSegments("das-hugel"))
	require.Equal(t, want, routeSegments("das-hugel-2026"))
	require.Equal(t, []int64{629046}, routeSegments("lite-das-hugel"))
	var newLiteView string
	require.NoError(t, pool.QueryRow(ctx, `SELECT pg_get_viewdef('lite_hugel_activities_2025')`).Scan(&newLiteView))
	require.Equal(t, oldLiteView, newLiteView)
	var metadataMatches bool
	require.NoError(t, pool.QueryRow(ctx, `SELECT a.display_name=b.display_name AND a.description=b.description
  FROM competitive_routes a, competitive_routes b WHERE a.name='das-hugel' AND b.name='das-hugel-2026'`).Scan(&metadataMatches))
	require.True(t, metadataMatches)

	// Only segment aggregation is under test. Remove parent-record constraints in
	// this disposable DB so fixtures need not populate unrelated Strava metadata.
	exec(`ALTER TABLE segment_efforts DROP CONSTRAINT segment_efforts_activities_id_fk,
       DROP CONSTRAINT segment_efforts_athletes_id_fk`)
	effortID := int64(0)
	addEfforts := func(activity int64, segments []int64, elapsed float64) {
		t.Helper()
		for _, segment := range segments {
			effortID++
			exec(`INSERT INTO segment_efforts
    (id, athlete_id, segment_id, name, elapsed_time, moving_time, start_date,
     start_date_local, distance, start_index, end_index, device_watts, average_watts, updated_at, activities_id)
    VALUES ($1, 1, $2, '', $3, $3, '2026-11-01', '2026-11-01', 1, 0, 1, false, 0, '2026-11-01', $4)`, effortID, segment, elapsed, activity)
		}
	}
	addEfforts(25, old, 10)
	addEfforts(26, want, 10)
	addEfforts(26, []int64{811022}, 20)    // Slower duplicate must not count.
	addEfforts(27, want[:len(want)-1], 10) // Missing Bell Mountain is not a full finish.
	db := database.New(pool)
	require.NoError(t, db.RefreshHugel2025Activities(ctx))
	require.NoError(t, db.RefreshHugel2026Activities(ctx))
	require.NoError(t, db.RefreshHugelLite2025Activities(ctx))
	for view, expected := range map[string][]int64{
		"hugel_activities_2025":      {25},
		"hugel_activities_2026":      {26},
		"hugel_activities":           {26},
		"lite_hugel_activities_2025": {26, 27},
		"lite_hugel_activities":      {26, 27},
	} {
		var ids []int64
		require.NoError(t, pool.QueryRow(ctx, `SELECT array_agg(activity_id ORDER BY activity_id) FROM `+view).Scan(&ids))
		require.Equal(t, expected, ids, view)
	}
	var total float64
	require.NoError(t, pool.QueryRow(ctx, `SELECT total_time_seconds FROM hugel_activities_2026 WHERE activity_id=26`).Scan(&total))
	require.Equal(t, float64(len(want)*10), total)
	missing, err := db.MissingSegments(ctx, 27)
	require.NoError(t, err)
	require.Equal(t, []string{"3145259"}, missing)
	migration("down")
	require.Equal(t, old, routeSegments("das-hugel"))
	var count int
	require.NoError(t, pool.QueryRow(ctx, `SELECT count(*) FROM competitive_routes WHERE name='das-hugel-2026'`).Scan(&count))
	require.Zero(t, count)
	var currentID int64
	require.NoError(t, pool.QueryRow(ctx, `SELECT activity_id FROM hugel_activities`).Scan(&currentID))
	require.Equal(t, int64(25), currentID)
}

// All migrations must leave the current leaderboard readable without a worker run.
func TestHugel2026PopulatedAfterMigrations(t *testing.T) {
	url, cleanup, err := postgres.Open()
	require.NoError(t, err)
	t.Cleanup(cleanup)
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, url)
	require.NoError(t, err)
	t.Cleanup(pool.Close)
	var populated bool
	require.NoError(t, pool.QueryRow(ctx, `SELECT ispopulated FROM pg_matviews WHERE matviewname='hugel_activities_2026'`).Scan(&populated))
	require.True(t, populated)
	// Also repair a deployment that already ran migration 42 without a refresh.
	_, err = pool.Exec(ctx, `REFRESH MATERIALIZED VIEW hugel_activities_2026 WITH NO DATA`)
	require.NoError(t, err)
	sql, err := os.ReadFile("migrations/000043_populate_hugel_2026.up.sql")
	require.NoError(t, err)
	_, err = pool.Exec(ctx, string(sql))
	require.NoError(t, err)
	var count int64
	require.NoError(t, pool.QueryRow(ctx, `SELECT count(*) FROM hugel_activities`).Scan(&count))
}
