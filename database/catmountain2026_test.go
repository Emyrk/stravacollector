package database_test

import (
	"context"
	"os"
	"strconv"
	"testing"

	"github.com/Emyrk/strava/database"
	"github.com/Emyrk/strava/database/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

func TestCatMountain2026Alternative(t *testing.T) {
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
		sql, err := os.ReadFile("migrations/000044_cat_mountain_alternative." + direction + ".sql")
		require.NoError(t, err)
		exec(string(sql))
	}
	// Test populated upgrades from the old definition, not just empty fresh installs.
	migration("down")
	exec(`UPDATE competitive_routes SET segments=ARRAY[891091,3145259,629046]::bigint[] WHERE name='das-hugel-2026'`)
	// Only scoring is under test. Avoid unrelated Strava athlete/activity metadata.
	exec(`ALTER TABLE segment_efforts DROP CONSTRAINT segment_efforts_activities_id_fk,
 DROP CONSTRAINT segment_efforts_athletes_id_fk;
 ALTER TABLE segments DROP CONSTRAINT segments_map_id_fkey`)
	effortID := int64(0)
	add := func(activity, segment int64, elapsed float64) int64 {
		t.Helper()
		effortID++
		exec(`INSERT INTO segment_efforts
 (id, athlete_id, segment_id, name, elapsed_time, moving_time, start_date,
 start_date_local, distance, start_index, end_index, device_watts, average_watts, updated_at, activities_id)
 VALUES ($1, 1, $2, '', $3, $3, '2026-11-01', '2026-11-01', 1, 0, 1, false, 0, '2026-11-01', $4)`, effortID, segment, elapsed, activity)
		return effortID
	}
	for activity := int64(1); activity <= 8; activity++ {
		add(activity, 629046, 10)
		if activity != 6 {
			add(activity, 3145259, 10)
		}
	}
	originalID := add(1, 891091, 10)
	alternateID := add(2, 784585, 11)
	add(3, 891091, 20)
	alternateFastID := add(3, 784585, 9)
	originalFastID := add(4, 891091, 8)
	add(4, 784585, 20)
	// 5 has neither Cat Mountain; 6 has alternate but lacks Bell Mountain.
	add(6, 784585, 10)
	add(7, 891091, 30)
	add(7, 784585, 20)
	duplicateID := add(7, 784585, 7)
	add(7, 784585, 7) // Same elapsed: smaller effort ID must win deterministically.
	tieID := add(8, 784585, 6)
	add(8, 891091, 6) // Tie across alternatives also uses effort ID, not segment ID.
	// Include a historical full finish so unchanged history is not an empty check.
	var historicalRoute []int64
	require.NoError(t, pool.QueryRow(ctx, `SELECT segments FROM competitive_routes WHERE name='das-hugel'`).Scan(&historicalRoute))
	for _, segment := range historicalRoute {
		add(99, segment, 10)
	}
	// Lite remains canonical-only, even when a full-route alternative is accepted.
	exec(`INSERT INTO competitive_routes (name,display_name,description,segments)
 VALUES ('lite-das-hugel','Lite','Historical fixture',ARRAY[891091,629046]::bigint[])
 ON CONFLICT (name) DO UPDATE SET segments=EXCLUDED.segments`)
	add(98, 891091, 10)
	add(98, 629046, 10)
	// MissingHugelSegments returns segment metadata, so seed its required records.
	exec(`INSERT INTO segments
 (id,name,activity_type,distance,average_grade,maximum_grade,elevation_high,elevation_low,
 start_latlng,end_latlng,elevation_profile,climb_category,city,state,country,private,hazardous,
 created_at,updated_at,total_elevation_gain,map_id,total_effort_count,total_athlete_count,total_star_count,fetched_at)
 SELECT id,'','Ride',0,0,0,0,0,ARRAY[]::float8[],ARRAY[]::float8[],'',0,'','','',false,false,
 '2026-01-01','2026-01-01',0,'',0,0,0,'2026-01-01'
 FROM unnest(ARRAY[891091,3145259,629046]::bigint[]) AS id`)
	db := database.New(pool)
	require.NoError(t, db.RefreshHugel2025Activities(ctx))
	require.NoError(t, db.RefreshHugelLite2025Activities(ctx))
	snapshot := func(view string) string {
		t.Helper()
		var data string
		require.NoError(t, pool.QueryRow(ctx, `SELECT COALESCE(json_agg(v ORDER BY activity_id)::text,'[]') FROM `+view+` v`).Scan(&data))
		return data
	}
	oldFull, oldLite := snapshot("hugel_activities_2025"), snapshot("lite_hugel_activities_2025")
	require.NotEqual(t, "[]", oldFull)
	require.NotEqual(t, "[]", oldLite)
	// Dependencies and historical view definitions must survive up and down.
	views := []string{"hugel_activities", "athlete_hugel_count", "hugel_activities_2025", "lite_hugel_activities_2025", "lite_hugel_activities"}
	oids := make(map[string]uint32)
	for _, view := range views {
		var oid uint32
		require.NoError(t, pool.QueryRow(ctx, `SELECT $1::regclass::oid`, view).Scan(&oid))
		oids[view] = oid
	}
	checkPreserved := func() {
		t.Helper()
		for view, expected := range oids {
			var oid uint32
			require.NoError(t, pool.QueryRow(ctx, `SELECT $1::regclass::oid`, view).Scan(&oid))
			require.Equal(t, expected, oid, view)
		}
		require.NoError(t, db.RefreshHugel2025Activities(ctx))
		require.NoError(t, db.RefreshHugelLite2025Activities(ctx))
		require.Equal(t, oldFull, snapshot("hugel_activities_2025"))
		require.Equal(t, oldLite, snapshot("lite_hugel_activities_2025"))
		var alternateLiteCount int
		require.NoError(t, pool.QueryRow(ctx, `SELECT count(*) FROM lite_hugel_activities_2025 WHERE activity_id=2`).Scan(&alternateLiteCount))
		require.Zero(t, alternateLiteCount, "alternate-only must not qualify for historical Lite")
		var count int
		require.NoError(t, pool.QueryRow(ctx, `SELECT count(*) FROM athlete_hugel_count`).Scan(&count))
		var populated bool
		require.NoError(t, pool.QueryRow(ctx, `SELECT ispopulated FROM pg_matviews WHERE matviewname='hugel_activities_2026'`).Scan(&populated))
		require.True(t, populated)
		var route []int64
		require.NoError(t, pool.QueryRow(ctx, `SELECT segments FROM competitive_routes WHERE name='das-hugel-2026'`).Scan(&route))
		require.Equal(t, []int64{891091, 3145259, 629046}, route)
	}
	migration("up")
	checkPreserved()
	var qualified []int64
	require.NoError(t, pool.QueryRow(ctx, `SELECT array_agg(activity_id ORDER BY activity_id) FROM hugel_activities WHERE activity_id <= 8`).Scan(&qualified))
	require.Equal(t, []int64{1, 2, 3, 4, 7, 8}, qualified)
	for _, tc := range []struct {
		activity, segment, effort int64
		total                     float64
	}{
		{1, 891091, originalID, 30}, {2, 784585, alternateID, 31},
		{3, 784585, alternateFastID, 29}, {4, 891091, originalFastID, 28},
		{7, 784585, duplicateID, 27}, {8, 784585, tieID, 26},
	} {
		var total float64
		var segments []int64
		var efforts database.HugelSegmentEfforts
		require.NoError(t, pool.QueryRow(ctx, `SELECT total_time_seconds,segment_ids,efforts FROM hugel_activities_2026 WHERE activity_id=$1`, tc.activity).Scan(&total, &segments, &efforts))
		require.Equal(t, tc.total, total)
		require.ElementsMatch(t, []int64{891091, 3145259, 629046}, segments)
		require.Len(t, efforts, 3)
		catCount := 0
		for _, effort := range efforts {
			if effort.SegmentID == 891091 || effort.SegmentID == 784585 {
				catCount++
				require.Equal(t, tc.segment, int64(effort.SegmentID))
				require.Equal(t, tc.effort, effort.EffortID)
			}
		}
		require.Equal(t, 1, catCount)
		missing, err := db.MissingSegments(ctx, tc.activity)
		require.NoError(t, err)
		require.Empty(t, missing)
		details, err := db.MissingHugelSegments(ctx, tc.activity)
		require.NoError(t, err)
		require.Empty(t, details)
	}
	for activity, segment := range map[int64]int64{5: 891091, 6: 3145259} {
		missing, err := db.MissingSegments(ctx, activity)
		require.NoError(t, err)
		require.Equal(t, []string{strconv.FormatInt(segment, 10)}, missing)
		details, err := db.MissingHugelSegments(ctx, activity)
		require.NoError(t, err)
		require.Len(t, details, 1)
		require.Equal(t, segment, details[0].ID)
	}
	// An alternative is not an extra scoring segment when this route does not
	// require the canonical slot (as in the repository's default seed route).
	exec(`UPDATE competitive_routes SET segments=ARRAY[3145259,629046]::bigint[] WHERE name='das-hugel-2026'`)
	require.NoError(t, db.RefreshHugel2026Activities(ctx))
	var withoutCatTotal float64
	var withoutCatEfforts database.HugelSegmentEfforts
	require.NoError(t, pool.QueryRow(ctx, `SELECT total_time_seconds,efforts FROM hugel_activities_2026 WHERE activity_id=3`).Scan(&withoutCatTotal, &withoutCatEfforts))
	require.Equal(t, float64(20), withoutCatTotal)
	require.Len(t, withoutCatEfforts, 2)
	for _, effort := range withoutCatEfforts {
		require.NotContains(t, []int{891091, 784585}, effort.SegmentID)
	}
	exec(`UPDATE competitive_routes SET segments=ARRAY[891091,3145259,629046]::bigint[] WHERE name='das-hugel-2026'`)

	migration("down")
	checkPreserved()
	require.NoError(t, pool.QueryRow(ctx, `SELECT array_agg(activity_id ORDER BY activity_id) FROM hugel_activities WHERE activity_id <= 8`).Scan(&qualified))
	require.Equal(t, []int64{1, 3, 4, 7, 8}, qualified)
	var oldTotal float64
	require.NoError(t, pool.QueryRow(ctx, `SELECT total_time_seconds FROM hugel_activities_2026 WHERE activity_id=3`).Scan(&oldTotal))
	require.Equal(t, float64(40), oldTotal)
	migration("up")
	checkPreserved()
}
