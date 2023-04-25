BEGIN;

ALTER TABLE athlete_load
	ADD COLUMN earliest_activity timestamp with time zone NOT NULL DEFAULT NOW()::timestamp without time zone;

COMMENT ON COLUMN athlete_load.earliest_activity IS 'The earliest activity found for the athlete';

ALTER TABLE athlete_load
	ADD COLUMN earliest_activity_done boolean NOT NULL DEFAULT false;

COMMENT ON COLUMN athlete_load.earliest_activity_done IS 'Loading backwards is done';

COMMIT;