BEGIN;

ALTER TABLE athletes RENAME TO athlete_logins;

ALTER TABLE athlete_logins
	RENAME COLUMN id TO athlete_id;

ALTER TABLE athlete_logins
	RENAME COLUMN premium TO summit;

ALTER TABLE athlete_logins
    ADD COLUMN id uuid NOT NULL DEFAULT gen_random_uuid();

ALTER TABLE athlete_logins
    DROP COLUMN username,
	DROP COLUMN firstname,
	DROP COLUMN lastname,
	DROP COLUMN sex,
	DROP COLUMN raw;

CREATE TABLE athletes (
	id bigint NOT NULL,
	summit boolean NOT NULL,
	username text NOT NULL,
	firstname text NOT NULL,
	lastname text NOT NULL,
	sex text NOT NULL,
	city text NOT NULL,
	state text NOT NULL,
	country text NOT NULL,
	follow_count integer NOT NULL,
	friend_count integer NOT NULL,
	measurement_preference text NOT NULL,
	ftp double precision NOT NULL,
	weight double precision NOT NULL,
	clubs json NOT NULL,
	created_at timestamp with time zone NOT NULL,
	updated_at timestamp with time zone NOT NULL,
	fetched_at timestamp with time zone NOT NULL
);

COMMENT ON COLUMN athletes.measurement_preference IS 'feet or meters';

COMMIT;