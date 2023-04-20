BEGIN;

CREATE TABLE webhook_dump (
	id uuid NOT NULL PRIMARY KEY,
	recorded_at timestamp NOT NULL,
	raw text NOT NULL
);

COMMIT;