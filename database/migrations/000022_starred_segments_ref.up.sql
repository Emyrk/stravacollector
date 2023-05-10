BEGIN;

ALTER TABLE starred_segments
	ADD CONSTRAINT starred_segments_segment_id_fkey FOREIGN KEY (segment_id) REFERENCES segments(id) ON DELETE CASCADE;

COMMIT;