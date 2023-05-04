BEGIN;

CREATE OR REPLACE VIEW super_hugel_activities AS
	SELECT
		*
	FROM
		(
			SELECT
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
					-- This query returns only the best effort per (segment_id, athlete_id)
					SELECT DISTINCT ON (athlete_id, segment_id)
						*
					FROM
						segment_efforts
					WHERE
						segment_id = any(ARRAY(SELECT segments FROM competitive_routes WHERE name = 'das-hugel'))
					ORDER BY
						athlete_id, segment_id, elapsed_time ASC
				) as hugel_efforts
				-- Each activity will now be represented by a single aggregated row
			GROUP BY
				(athlete_id)
		) AS merged
	WHERE
		segment_ids @> ARRAY(SELECT segments FROM competitive_routes WHERE name = 'das-hugel')
;

COMMIT;