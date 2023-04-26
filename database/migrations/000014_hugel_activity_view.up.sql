BEGIN;

CREATE TABLE IF NOT EXISTS competitive_routes(
    name text PRIMARY KEY NOT NULL,
    display_name text NOT NULL,
    description text NOT NULL,
    segments bigint[] NOT NULL
);

INSERT INTO competitive_routes(name, display_name, description, segments)
VALUES (
        'das-hugel',
        'Das Hugel',
        'The legendary Das Hugel route.',
        ARRAY[629046, 6744304, 910045, 628785, 629546, 628842]
)
ON CONFLICT (name) DO UPDATE SET
	display_name = EXCLUDED.display_name,
	description = EXCLUDED.description,
	segments = EXCLUDED.segments
;

CREATE OR REPLACE VIEW hugel_activities AS
	SELECT
		*
	FROM
		(
			SELECT
				activities_id,
				-- segment_ids is all the segments this activity has efforts on.
				-- Only segments in the provided list are considered.
				array_agg(segment_id) AS segment_ids,
				-- Sum is the total time of all the efforts.
				sum(elapsed_time),
				-- A json struct containing each effort details.
				json_agg(
					json_build_object(
						'effort_id', id,
						'start_date', start_date,
						'segment_id', segment_id,
						'elapsed_time', elapsed_time,
						'moving_time', moving_time,
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
						segment_id = any(ARRAY(SELECT segments FROM competitive_routes WHERE name = 'das-hugel'))
					ORDER BY
						activities_id, segment_id, elapsed_time ASC
				) as hugel_efforts
				-- Each activity will now be represented by a single aggregated row
			GROUP BY
				activities_id
		) AS merged
	WHERE
		segment_ids @> ARRAY(SELECT segments FROM competitive_routes WHERE name = 'das-hugel')
;

COMMENT ON VIEW hugel_activities IS 'This view contains all activities that classify as a "hugel" and their best efforts on each segment.';

COMMIT;