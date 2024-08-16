CREATE INDEX activity_summary_start_date_idx ON activity_summary(start_date);
CREATE INDEX segment_efforts_distinct_idx ON segment_efforts(activities_id, segment_id);
CREATE INDEX segment_efforts_elapsed_time_idx ON segment_efforts(elapsed_time);