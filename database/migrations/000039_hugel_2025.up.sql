BEGIN;

-- -- Make materialized
CREATE MATERIALIZED VIEW lite_hugel_activities AS
SELECT
	*
FROM
	(
		SELECT
			activities_id AS activity_id,
			athlete_id,
			-- segment_ids is all the segments this activity has efforts on.
			-- Only segments in the provided list are considered.
			array_agg(segment_id) :: BIGINT[] AS segment_ids,
			-- Sum is the total time of all the efforts.
			sum(elapsed_time) AS total_time_seconds,
			-- A json struct containing each effort details.
			json_agg(
					json_build_object(
							'activity_id', activities_id,
							'effort_id', id,
							'start_date', start_date,
							'segment_id', segment_id,
							'elapsed_time', elapsed_time,
							'moving_time', moving_time,
							'device_watts', device_watts,
							'average_watts', average_watts
					)
			) AS efforts
		FROM
			(
				-- This query returns only the best effort per (segment_id, activity_id)
				SELECT DISTINCT ON (activities_id, segment_id)
					*
				FROM
					segment_efforts
				WHERE
					segment_id = any(ARRAY(SELECT segments FROM competitive_routes WHERE name = 'lite-das-hugel'))
				ORDER BY
					activities_id, segment_id, elapsed_time ASC
			) as hugel_efforts
		-- Each activity will now be represented by a single aggregated row
		GROUP BY
			(activities_id, athlete_id)
	) AS merged
WHERE
	segment_ids @> ARRAY(SELECT segments FROM competitive_routes WHERE name = 'lite-das-hugel')
;

CREATE MATERIALIZED VIEW hugel_activities AS
SELECT merged.activity_id,
	   merged.athlete_id,
	   merged.segment_ids,
	   merged.total_time_seconds,
	   merged.efforts
FROM ( SELECT hugel_efforts.activities_id AS activity_id,
			  hugel_efforts.athlete_id,
			  array_agg(hugel_efforts.segment_id) AS segment_ids,
			  sum(hugel_efforts.elapsed_time) AS total_time_seconds,
			  json_agg(json_build_object('activity_id', hugel_efforts.activities_id, 'effort_id', hugel_efforts.id, 'start_date', hugel_efforts.start_date, 'segment_id', hugel_efforts.segment_id, 'elapsed_time', hugel_efforts.elapsed_time, 'moving_time', hugel_efforts.moving_time, 'device_watts', hugel_efforts.device_watts, 'average_watts', hugel_efforts.average_watts)) AS efforts
	   FROM ( SELECT DISTINCT ON (segment_efforts.activities_id, segment_efforts.segment_id) segment_efforts.id,
																							 segment_efforts.athlete_id,
																							 segment_efforts.segment_id,
																							 segment_efforts.name,
																							 segment_efforts.elapsed_time,
																							 segment_efforts.moving_time,
																							 segment_efforts.start_date,
																							 segment_efforts.start_date_local,
																							 segment_efforts.distance,
																							 segment_efforts.start_index,
																							 segment_efforts.end_index,
																							 segment_efforts.device_watts,
																							 segment_efforts.average_watts,
																							 segment_efforts.kom_rank,
																							 segment_efforts.pr_rank,
																							 segment_efforts.updated_at,
																							 segment_efforts.activities_id
			  FROM segment_efforts
			  WHERE (segment_efforts.segment_id = ANY (ARRAY( SELECT competitive_routes.segments
															  FROM competitive_routes
															  WHERE (competitive_routes.name = 'das-hugel'::text))))
			  ORDER BY segment_efforts.activities_id, segment_efforts.segment_id, segment_efforts.elapsed_time) hugel_efforts
	   GROUP BY hugel_efforts.activities_id, hugel_efforts.athlete_id) merged
WHERE (merged.segment_ids @> ARRAY( SELECT competitive_routes.segments
									FROM competitive_routes
									WHERE (competitive_routes.name = 'das-hugel'::text)))
WITH NO DATA;

COMMIT;