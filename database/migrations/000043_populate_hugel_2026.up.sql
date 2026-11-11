-- Make the current leaderboard readable before the API starts.
-- Migration 42 created this materialized view WITH NO DATA.
REFRESH MATERIALIZED VIEW hugel_activities_2026;
