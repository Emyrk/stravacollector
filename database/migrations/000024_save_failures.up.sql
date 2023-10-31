BEGIN;

CREATE TABLE failed_jobs (
	id uuid NOT NULL PRIMARY KEY,
	recorded_at timestamp NOT NULL,
	raw text NOT NULL
);

COMMENT ON TABLE failed_jobs IS 'A table to store failed job information for potential debugging.';
COMMENT ON COLUMN failed_jobs.id IS 'Some random uuid';
COMMENT ON COLUMN failed_jobs.raw IS 'Some text. Probably a JSON string.';

ALTER TABLE athlete_load ADD COLUMN next_load_not_before timestamp with time zone NOT NULL DEFAULT now();

COMMIT;