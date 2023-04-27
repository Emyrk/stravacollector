BEGIN;

CREATE TABLE api_tokens(
    id uuid PRIMARY KEY NOT NULL,
    name TEXT NOT NULL,
	athlete_id BIGINT NOT NULL
        REFERENCES athlete_logins(athlete_id)
        ON DELETE CASCADE,
    hashed_token TEXT NOT NULL,
	created_at timestamp with time zone NOT NULL,
	updated_at timestamp with time zone NOT NULL,
	last_used_at timestamp with time zone NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    -- 7 days
	lifetime_seconds bigint DEFAULT 604800 NOT NULL
);

COMMENT ON COLUMN api_tokens.lifetime_seconds IS 'The amount of time to renew the token for.';

COMMIT;