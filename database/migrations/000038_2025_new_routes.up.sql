BEGIN;

ALTER MATERIALIZED VIEW hugel_activities RENAME TO hugel_activities_2024;
ALTER MATERIALIZED VIEW hugel_activities_lite RENAME TO hugel_activities_lite_2024;

COMMIT;