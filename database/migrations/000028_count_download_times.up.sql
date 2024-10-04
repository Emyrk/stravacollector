ALTER TYPE activity_detail_source ADD VALUE 'zero_segment_refetch';
ALTER TABLE activity_summary ADD COLUMN download_count INT DEFAULT 0 NOT NULL;