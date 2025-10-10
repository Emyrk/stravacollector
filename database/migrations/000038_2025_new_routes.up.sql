BEGIN;

ALTER MATERIALIZED VIEW hugel_activities RENAME TO hugel_activities_2024;
ALTER MATERIALIZED VIEW lite_hugel_activities RENAME TO lite_hugel_activities_2024;

COMMIT;