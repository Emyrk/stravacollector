-- name: RefreshHugelActivities :exec
REFRESH MATERIALIZED VIEW hugel_activities;

-- name: RefreshHugelLiteActivities :exec
REFRESH MATERIALIZED VIEW lite_hugel_activities;

-- name: RefreshHugel2023Activities :exec
REFRESH MATERIALIZED VIEW hugel_activities_2023;

-- name: RefreshSuperHugelActivities :exec
REFRESH MATERIALIZED VIEW super_hugel_activities;


-- name: AthleteHugelActivites :many
SELECT
    sqlc.embed(hugel_activities),
    sqlc.embed(activity_summary)
FROM
    hugel_activities
INNER JOIN
	activity_summary ON hugel_activities.activity_id = activity_summary.id
WHERE
	hugel_activities.athlete_id = @athlete_id;


-- name: HugelLeaderboard :many
SELECT
	(SELECT min(total_time_seconds) FROM hugel_activities) :: BIGINT AS best_time,
	ROW_NUMBER() over(ORDER BY total_time_seconds ASC) AS rank,
	athlete_bests.activity_id,
	athlete_bests.athlete_id,
	athlete_bests.total_time_seconds,
	athlete_bests.efforts,

	activity_summary.name,
	activity_summary.device_watts,
	activity_summary.distance,
	activity_summary.moving_time,
	activity_summary.elapsed_time,
	activity_summary.total_elevation_gain,
	activity_summary.start_date,
	activity_summary.achievement_count,
	activity_summary.average_heartrate,
	activity_summary.average_speed,

	activity_detail.suffer_score,
	activity_detail.average_watts,
	activity_detail.average_cadence,

	athletes.firstname,
	athletes.lastname,
	athletes.username,
	athletes.profile_pic_link,
	athletes.sex,
	COALESCE(hugel_count.count, 0) AS hugel_count
FROM
	(
		SELECT DISTINCT ON (hugel_activities.athlete_id)
			hugel_activities.*
		FROM
			hugel_activities
		INNER JOIN
			activity_summary ON hugel_activities.activity_id = activity_summary.id
		WHERE
		CASE WHEN
			@after :: timestamp != '0001-01-01 00:00:00Z'
			AND @before :: timestamp != '0001-01-01 00:00:00Z' THEN
			activity_summary.start_date >= @after :: timestamp AND activity_summary.start_date <= @before :: timestamp
		ELSE TRUE END
		ORDER BY
			hugel_activities.athlete_id, hugel_activities.total_time_seconds ASC
	) AS athlete_bests
INNER JOIN
	athletes ON athlete_bests.athlete_id = athletes.id
LEFT JOIN athlete_hugel_count AS hugel_count
	ON hugel_count.athlete_id = athlete_bests.athlete_id
INNER JOIN
	activity_summary ON athlete_bests.activity_id = activity_summary.id
INNER JOIN
	activity_detail ON athlete_bests.activity_id = activity_detail.id
WHERE
    CASE WHEN @athlete_id > 0 THEN athlete_bests.athlete_id = @athlete_id ELSE TRUE END
    AND
	CASE WHEN
    	@after :: timestamp != '0001-01-01 00:00:00Z'
    		AND @before :: timestamp != '0001-01-01 00:00:00Z' THEN
			activity_summary.start_date >= @after :: timestamp AND activity_summary.start_date <= @before :: timestamp
    ELSE TRUE END
	ORDER BY
		athlete_bests.total_time_seconds ASC
;

-- name: SuperHugelLeaderboard :many
SELECT
	(SELECT min(total_time_seconds) FROM super_hugel_activities) :: BIGINT AS best_time,
	ROW_NUMBER() over(ORDER BY total_time_seconds ASC) AS rank,
	athlete_bests.athlete_id,
	athlete_bests.total_time_seconds,
	athlete_bests.efforts,

	athletes.firstname,
	athletes.lastname,
	athletes.username,
	athletes.profile_pic_link,
	athletes.sex,
	hugel_count.count AS hugel_count
FROM
	(
		SELECT DISTINCT ON (athlete_id)
			*
		FROM
			super_hugel_activities
		ORDER BY
			athlete_id, total_time_seconds ASC
	) AS athlete_bests
INNER JOIN
	athletes ON athlete_bests.athlete_id = athletes.id
INNER JOIN
	athlete_hugel_count AS hugel_count
		ON hugel_count.athlete_id = athlete_bests.athlete_id
WHERE
	CASE WHEN @athlete_id > 0 THEN athlete_bests.athlete_id = @athlete_id ELSE TRUE END
ORDER BY
	athlete_bests.total_time_seconds ASC
;

-- name: GetCompetitiveRoute :one
SELECT
	competitive_routes.name, competitive_routes.display_name, competitive_routes.description, (
	SELECT
		json_agg(
			json_build_object(
				'id',segments.id,
				'name',segments.name
			)
		) AS segment_summaries
	FROM
		segments
	WHERE
		id = ANY(competitive_routes.segments)
)
FROM
	competitive_routes
WHERE
	competitive_routes.name = @route_name
LIMIT 1;


-- name: MissingHugelSegments :many
SELECT
	*
FROM
	segments
WHERE
	id = ANY(
		select unnest(segments) as data
		from (SELECT segments FROM competitive_routes WHERE name = 'das-hugel') as hugel
		except
		select segment_id as data
		from segment_efforts WHERE
		    activities_id = @activity_id
	);
