BEGIN;

-- Preserve the wrapper and its dependents while replacing its backing matview.
-- The temporary redirect is invisible outside this transaction.
CREATE OR REPLACE VIEW hugel_activities AS SELECT * FROM hugel_activities_2025;
DROP MATERIALIZED VIEW hugel_activities_2026;

-- Cat Mountain alternatives share one required scoring slot. Keep the actual
-- winning effort identity in JSON for links, splits, and auditability.
CREATE MATERIALIZED VIEW hugel_activities_2026 AS
WITH normalized AS (
    SELECT segment_efforts.*,
           CASE WHEN segment_id = 784585 THEN 891091::bigint
                ELSE segment_id END AS scoring_segment_id
    FROM segment_efforts
), best_efforts AS (
    SELECT DISTINCT ON (activities_id, scoring_segment_id) *
    FROM normalized
    WHERE scoring_segment_id = ANY (
        (SELECT segments FROM competitive_routes WHERE name = 'das-hugel-2026')::bigint[])
    ORDER BY activities_id, scoring_segment_id, elapsed_time, id
), merged AS (
    SELECT activities_id AS activity_id,
           athlete_id,
           array_agg(scoring_segment_id ORDER BY scoring_segment_id) AS segment_ids,
           sum(elapsed_time) AS total_time_seconds,
           json_agg(json_build_object(
               'activity_id', activities_id, 'effort_id', id,
               'start_date', start_date, 'segment_id', segment_id,
               'elapsed_time', elapsed_time, 'moving_time', moving_time,
               'device_watts', device_watts, 'average_watts', average_watts
           ) ORDER BY scoring_segment_id) AS efforts
    FROM best_efforts
    GROUP BY activities_id, athlete_id
)
SELECT activity_id, athlete_id, segment_ids, total_time_seconds, efforts
FROM merged
WHERE segment_ids @> (
    SELECT segments FROM competitive_routes WHERE name = 'das-hugel-2026');

CREATE OR REPLACE VIEW hugel_activities AS SELECT * FROM hugel_activities_2026;

COMMIT;
