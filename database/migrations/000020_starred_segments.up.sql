BEGIN;

CREATE TABLE starred_segments(
	athlete_id bigint NOT NULL
	    REFERENCES athletes(id),
	segment_id bigint NOT NULL,
	starred boolean NOT NULL,
	updated_at timestamp with time zone NOT NULL,
	PRIMARY KEY (athlete_id, segment_id)
);

COMMIT;