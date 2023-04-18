BEGIN;

CREATE TABLE athletes (
    -- Strava Details
	id BIGINT PRIMARY KEY,
	premium boolean NOT NULL,
	username text NOT NULL,
	firstname text NOT NULL,
	lastname text NOT NULL,
	sex text NOT NULL,
	-- Authentication
	provider_id text NOT NULL,
	created_at timestamp with time zone NOT NULL,
	updated_at timestamp with time zone NOT NULL,
	oauth_access_token text NOT NULL,
	oauth_refresh_token text NOT NULL,
	oauth_expiry timestamp with time zone NOT NULL,
    -- Raw
    raw text NOT NULL
);

COMMENT ON COLUMN athletes.provider_id IS 'Oauth app client ID';

CREATE TABLE segments (
	id integer PRIMARY KEY,
	name text NOT NULL
);

CREATE TABLE athlete_efforts (
	 id uuid PRIMARY KEY,
	 athlete_id integer REFERENCES athletes(id) NOT NULL,
	 segment_id integer REFERENCES segments(id) NOT NULL,
	 last_checked timestamp NOT NULL,
	 best_effort_id integer NOT NULL,
	 kom_rank integer NOT NULL
);

COMMIT;