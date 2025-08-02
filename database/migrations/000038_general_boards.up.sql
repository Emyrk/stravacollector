BEGIN;

-- The new pkey will include all columns that are used to differentiate the routes.
ALTER TABLE competitive_routes DROP CONSTRAINT competitive_routes_pkey;

ALTER TABLE
	competitive_routes
	ADD COLUMN year INT NOT NULL DEFAULT 2023,
	ADD COLUMN course TEXT NOT NULL DEFAULT 'full',
	ADD PRIMARY KEY (name, course, year)
;

COMMENT ON COLUMN competitive_routes.course
	IS 'The course name for the competitive route is used to differentiate between different versions of the same route, such as "lite" or "full".';
;


INSERT INTO
	competitive_routes (name, display_name, description, year, course, segments)
VALUES
	-- 2024
	('tour-das-hugel', 'Tour Das Hügel', 'The legendary Tour Das Hügel route.', 2024, 'full', ARRAY[629046,609560,629546,628842,33454575,740569,617537,995700,628845,10355605,28704989,628782,681737,891091,626149,10721603,626148,1776842,405437,674506,648953,37363955,705265]),
	('tour-das-hugel', 'Tour Das Hügel Lite', 'The short Tour Das Hügel route.', 2024, 'lite', ARRAY[10721603,629046,626148,681737,628782,648953,674506]),
	-- 2023
	('tour-das-hugel', 'Tour Das Hügel', 'The legendary Tour Das Hügel route.',2023, 'full', ARRAY[629046,609560,629546,628842,33454575,740569,617537,995700,628845,1469968,10355605,28704989,628782,681737,891091,626149,10721603,626148,1776842,405437])
;

CREATE OR REPLACE FUNCTION build_route_view_name(
	route_name TEXT,
	route_year INT,
	route_course TEXT
) RETURNS TEXT AS $$
BEGIN
	RETURN format(
			'%s_%s_%s-activities',
			replace(lower(route_name), '-', '_'),
			route_year,
			replace(lower(route_course), '-', '_')
		   );
END;
$$ LANGUAGE plpgsql IMMUTABLE;


CREATE OR REPLACE PROCEDURE create_route_activities_view(route_name TEXT, route_year INT, route_course TEXT)
	LANGUAGE plpgsql
AS $$
DECLARE
	view_name TEXT := build_route_view_name(route_name, route_year, route_course);
    sql TEXT;
BEGIN
	sql := format($fmt$
        CREATE MATERIALIZED VIEW IF NOT EXISTS %I AS
        SELECT *
        FROM (
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
            FROM (
            	-- This query returns only the best effort per (segment_id, activity_id)
				SELECT DISTINCT ON (activities_id, segment_id)
					*
				FROM
					segment_efforts
				WHERE
				segment_id = ANY (
					SELECT
						segments
					FROM
						competitive_routes
					WHERE
						name = %L
						AND year = %L
						AND course = %L
				)
                ORDER BY
                	activities_id, segment_id, elapsed_time ASC
			) as route_efforts
			-- Each activity will now be represented by a single aggregated row
			GROUP BY
				(activities_id, athlete_id)
        ) AS merged
        WHERE segment_ids @> ARRAY(
            SELECT
            	segments
            FROM
            	competitive_routes
            WHERE
            	name = %L
				AND year = %L
				AND course = %L
        )
    $fmt$, view_name,
		route_name, route_year, route_course,
		route_name, route_year, route_course);

	EXECUTE sql;
END;
$$;


CALL create_route_activities_view('tour-das-hugel', 2024, 'full');

COMMIT;