BEGIN;

ALTER TABLE activity_summary
	ADD COLUMN average_heartrate double precision NOT NULL DEFAULT 0,
	ADD COLUMN max_heartrate double precision NOT NULL DEFAULT 0;

COMMIT;