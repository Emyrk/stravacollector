-- name: BestRouteEfforts :many
-- BestRouteEfforts returns all activities that have efforts on all the provided segments.
-- The returned activities include the best effort for each segment.
-- This isn't used in the app, but is the foundation for the hugel view.
SELECT
	*
FROM
	(
		SELECT
			hugel_efforts.activities_id,
			-- segment_ids is all the segments this activity has efforts on.
			-- Only segments in the provided list are considered.
			array_agg(segment_id) AS segment_ids,
			-- Sum is the total time of all the efforts.
			sum(elapsed_time),
			-- A json struct containing each effort details.
			json_agg(
					json_build_object(
					    	'effort_id', id,
							'segment_id', segment_id,
							'elapsed_time', elapsed_time,
							'device_watts', device_watts,
					    	'average_watts', average_watts
						)
				)
		FROM
			(
				-- This query returns only the best effort per (segment_id, activity_id)
				SELECT DISTINCT ON (activities_id, segment_id)
					*
				FROM
					segment_efforts
				WHERE
				    -- ARRAY[629046, 6744304, 910045, 628785, 629546, 628842]
					segment_id = any(@expected_segments :: bigint[])
				ORDER BY
					activities_id, segment_id, elapsed_time ASC
			) as hugel_efforts
			-- Each activity will now be represented by a single aggregated row
		GROUP BY
			hugel_efforts.activities_id
	) AS merged
WHERE
	segment_ids @> @expected_segments :: bigint[]
;

-- name: AllCompetitiveRoutes :many
SELECT * FROM competitive_routes;