BEGIN;

ALTER TABLE athlete_load
	ADD COLUMN earliest_activity_id BIGINT DEFAULT 0 NOT NULL;

COMMIT;