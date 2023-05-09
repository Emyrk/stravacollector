-- name: StarSegments :one
INSERT INTO
	starred_segments(
		updated_at,
		athlete_id,
		segment_id,
	    starred
	)
SELECT
	Now() AS updated_at,
	unnest(@athlete_id::bigint[]) AS athlete_id,
	unnest(@segment_id::bigint[]) AS segment_id,
	unnest(@starred::boolean[]) AS starred
ON CONFLICT
	(athlete_id, segment_id)
	DO UPDATE SET
		updated_at = Now(),
		starred = EXCLUDED.starred
RETURNING *;
;

-- name: LoadedSegments :many
SELECT id, fetched_at FROM segments;

-- name: GetSegments :many
SELECT
    sqlc.embed(segments), sqlc.embed(maps)
FROM
    segments
LEFT JOIN
	maps ON segments.map_id = maps.id
WHERE segments.id = ANY(@segment_ids::bigint[])
;

-- name: GetPersonalSegments :many
-- For authenticated users
SELECT
	sqlc.embed(segments),
	sqlc.embed(maps),
	COALESCE(starred_segments.starred, false) as starred,

	-- SegmentEffort
	COALESCE(best_effort.id, -1) as best_effort_id,
	COALESCE(best_effort.elapsed_time, -1) as best_effort_elapsed_time,
	COALESCE(best_effort.moving_time, -1) as best_effort_moving_time,
	COALESCE(best_effort.start_date, '0001-01-01 00:00:00+00'::timestamp) as best_effort_start_date,
	COALESCE(best_effort.start_date_local, '0001-01-01 00:00:00+00'::timestamp) as best_effort_start_date_local,
	COALESCE(best_effort.device_watts, false) as best_effort_device_watts,
	COALESCE(best_effort.average_watts, -1) as best_effort_average_watts,
	COALESCE(best_effort.activities_id, -1) as best_effort_activities_id
FROM
	segments
LEFT JOIN
	maps ON segments.map_id = maps.id
LEFT JOIN
	-- Only for the authenticated user
	starred_segments
	    ON segments.id = starred_segments.segment_id AND starred_segments.athlete_id = @athlete_id
LEFT JOIN LATERAL
	(
	    SELECT DISTINCT ON (segment_efforts.athlete_id, segment_efforts.segment_id)
			*
	    FROM
	        segment_efforts
	    WHERE
	        athlete_id = @athlete_id AND
	        segment_id = segments.id
    	ORDER BY
			segment_efforts.athlete_id, segment_efforts.segment_id, elapsed_time ASC
	) best_effort ON best_effort.segment_id = segments.id
WHERE segments.id = ANY(@segment_ids::bigint[])
;

-- name: test :many
SELECT
	*
FROM
	segments
		LEFT JOIN
	maps ON segments.map_id = maps.id
		LEFT JOIN LATERAL (
		SELECT DISTINCT ON (segment_efforts.athlete_id, segment_efforts.segment_id)
			*
		FROM
			segment_efforts
		WHERE
				athlete_id = 20563755 AND
				segment_id = segments.id
		ORDER BY
			segment_efforts.athlete_id, segment_efforts.segment_id, elapsed_time ASC
		) segment_effort ON segment_effort.segment_id = segments.id

WHERE segments.id = ANY(ARRAY[628842]);

-- name: UpsertSegmentEffort :one
INSERT INTO
	segment_efforts(
		updated_at,
		id, athlete_id, segment_id, name, elapsed_time,
		moving_time, start_date, start_date_local, distance,
		start_index, end_index, device_watts, average_watts,
		kom_rank, pr_rank, activities_id
	)
VALUES
	(Now(), $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
ON CONFLICT
	(id)
	DO UPDATE SET
		updated_at = Now(),
		athlete_id = $2,
		segment_id = $3,
		name = $4,
		elapsed_time = $5,
		moving_time = $6,
		start_date = $7,
		start_date_local = $8,
		distance = $9,
		start_index = $10,
		end_index = $11,
		device_watts = $12,
		average_watts = $13,
		kom_rank = $14,
		pr_rank = $15,
		activities_id = $16
	RETURNING *;


-- name: UpsertSegment :one
INSERT INTO
	segments(
	id, name, activity_type, distance, average_grade,
	maximum_grade, elevation_high, elevation_low, start_latlng, end_latlng,
	elevation_profile, climb_category, city, state, country, private, hazardous,
	created_at, updated_at, total_elevation_gain, map_id, total_effort_count,
	total_athlete_count, total_star_count, fetched_at
	)
VALUES
	($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17,
	 $18, $19, $20, $21, $22, $23, $24, Now())
ON CONFLICT
	(id)
	DO UPDATE SET
	name = CASE WHEN $2 != '' THEN $2 ELSE segments.name END,
	activity_type = CASE WHEN $3 != '' THEN $3 ELSE segments.activity_type END,
	distance = CASE WHEN $4 != 0 THEN $4 ELSE segments.distance END,
	average_grade = CASE WHEN $5 != 0 THEN $5 ELSE segments.average_grade END,
	maximum_grade = CASE WHEN $6 != 0 THEN $6 ELSE segments.maximum_grade END,
	elevation_high = CASE WHEN $7 != 0 THEN $7 ELSE segments.elevation_high END,
	elevation_low = CASE WHEN $8 != 0 THEN $8 ELSE segments.elevation_low END,
	start_latlng = $9,
	end_latlng = $10,
	elevation_profile = CASE WHEN $11 != '' THEN $11 ELSE segments.elevation_profile END,
	climb_category = CASE WHEN $12 != 0 THEN $12 ELSE segments.climb_category END,
	city = CASE WHEN $13 != '' THEN $13 ELSE segments.city END,
	state = CASE WHEN $14 != '' THEN $14 ELSE segments.state END,
	country = CASE WHEN $15 != '' THEN $15 ELSE segments.country END,
	private = $16,
	hazardous = $17,
	created_at = CASE WHEN $18 != '0001-01-01 00:00:00+00' THEN $18 ELSE segments.created_at END,
	updated_at = CASE WHEN $19 != '0001-01-01 00:00:00+00' THEN $18 ELSE segments.updated_at END,
	total_elevation_gain = CASE WHEN $20 != 0 THEN $20 ELSE segments.total_elevation_gain END,
	map_id = CASE WHEN $21 != '' THEN $21 ELSE segments.map_id END,
	total_effort_count = CASE WHEN $22 != 0 THEN $22 ELSE segments.total_effort_count END,
	total_athlete_count = CASE WHEN $23 != 0 THEN $23 ELSE segments.total_athlete_count END,
	total_star_count = CASE WHEN $24 != 0 THEN $24 ELSE segments.total_star_count END,
	fetched_at = Now()

RETURNING *;
