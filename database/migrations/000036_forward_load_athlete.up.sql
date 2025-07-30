BEGIN;

CREATE TABLE athlete_forward_load(
	athlete_id BIGINT
	 PRIMARY KEY
	 REFERENCES athletes(id) ON DELETE CASCADE
     NOT NULL,


	-- Load params
	activity_time_after TIMESTAMP WITH TIME ZONE NOT NULL,

	-- Metadata
	last_load_complete boolean NOT NULL DEFAULT false,
	last_touched TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	next_load_not_before TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);


COMMENT ON TABLE athlete_forward_load IS 'Tracks loading athlete activities. Must be an authenticated athlete.';
COMMENT ON COLUMN athlete_forward_load.last_touched IS 'Timestamp this row was last updated.';
COMMENT ON COLUMN athlete_forward_load.next_load_not_before IS 'Timestamp when the next load can be attempted.';

COMMIT;