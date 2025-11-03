BEGIN;

DROP VIEW IF EXISTS athlete_hugel_count;
DROP VIEW IF EXISTS lite_hugel_activities;
DROP MATERIALIZED VIEW IF EXISTS lite_hugel_activities_2024;
DROP MATERIALIZED VIEW IF EXISTS lite_hugel_activities_2025;

DROP MATERIALIZED VIEW IF EXISTS hugel_activities_2024;
CREATE MATERIALIZED VIEW hugel_activities_2024 AS
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
															  WHERE (competitive_routes.name = 'das-hugel-2024'::text))))
			  ORDER BY segment_efforts.activities_id, segment_efforts.segment_id, segment_efforts.elapsed_time) hugel_efforts
	   GROUP BY hugel_efforts.activities_id, hugel_efforts.athlete_id) merged
WHERE (merged.segment_ids @> ARRAY( SELECT competitive_routes.segments
									FROM competitive_routes
									WHERE (competitive_routes.name = 'das-hugel-2024'::text)))
WITH NO DATA;

COMMENT ON MATERIALIZED VIEW hugel_activities_2024 IS 'This view contains all activities that classify as a "hugel" and their best efforts on each segment.';


CREATE MATERIALIZED VIEW lite_hugel_activities_2024 AS
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
			  WHERE ((NOT (segment_efforts.activities_id IN ( SELECT hugel_activities_2024.activity_id
															  FROM hugel_activities_2024))) AND (segment_efforts.segment_id = ANY (ARRAY( SELECT competitive_routes.segments
																																		  FROM competitive_routes
																																		  WHERE (competitive_routes.name = 'lite-das-hugel-2024'::text)))))
			  ORDER BY segment_efforts.activities_id, segment_efforts.segment_id, segment_efforts.elapsed_time) hugel_efforts
	   GROUP BY hugel_efforts.activities_id, hugel_efforts.athlete_id) merged
WHERE (merged.segment_ids @> ARRAY( SELECT competitive_routes.segments
									FROM competitive_routes
									WHERE (competitive_routes.name = 'lite-das-hugel-2024'::text)))
WITH NO DATA;


CREATE OR REPLACE VIEW athlete_hugel_count AS
SELECT
	athlete_id, count(*) AS count
FROM
	athletes
		INNER JOIN
	hugel_activities
	ON athletes.id = hugel_activities.athlete_id
GROUP BY athlete_id;


CREATE MATERIALIZED VIEW lite_hugel_activities_2025 AS
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
			  WHERE ((NOT (segment_efforts.activities_id IN ( SELECT hugel_activities_2025.activity_id
															  FROM hugel_activities_2025))) AND (segment_efforts.segment_id = ANY (ARRAY( SELECT competitive_routes.segments
																																		  FROM competitive_routes
																																		  WHERE (competitive_routes.name = 'lite-das-hugel'::text)))))
			  ORDER BY segment_efforts.activities_id, segment_efforts.segment_id, segment_efforts.elapsed_time) hugel_efforts
	   GROUP BY hugel_efforts.activities_id, hugel_efforts.athlete_id) merged
WHERE (merged.segment_ids @> ARRAY( SELECT competitive_routes.segments
									FROM competitive_routes
									WHERE (competitive_routes.name = 'lite-das-hugel'::text)))
WITH NO DATA;


CREATE VIEW lite_hugel_activities AS
SELECT lite_hugel_activities_2025.activity_id,
	   lite_hugel_activities_2025.athlete_id,
	   lite_hugel_activities_2025.segment_ids,
	   lite_hugel_activities_2025.total_time_seconds,
	   lite_hugel_activities_2025.efforts
FROM lite_hugel_activities_2025;

COMMIT;