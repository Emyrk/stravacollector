BEGIN;

CREATE TABLE athlete_eddingtons(
	athlete_id BIGINT
		PRIMARY KEY
		REFERENCES athletes(id) ON DELETE CASCADE
		NOT NULL,

	miles_histogram integer[] NOT NULL DEFAULT '{}' :: integer[], -- [ride_count]
	current_eddington integer NOT NULL DEFAULT 0,
	last_calculated TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
	total_activities INTEGER NOT NULL DEFAULT 0
);

COMMIT;