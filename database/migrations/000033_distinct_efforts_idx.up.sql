CREATE INDEX IF NOT EXISTS segment_efforts_distinct_effort_idx ON segment_efforts (athlete_id, segment_id, elapsed_time);

COMMENT ON INDEX segment_efforts_distinct_effort_idx IS 'Index to support GetBestPersonalSegmentEffort query';