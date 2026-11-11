BEGIN;
CREATE OR REPLACE VIEW hugel_activities AS SELECT * FROM hugel_activities_2025;
DROP MATERIALIZED VIEW hugel_activities_2026;
DELETE FROM competitive_routes WHERE name = 'das-hugel-2026';
COMMIT;
