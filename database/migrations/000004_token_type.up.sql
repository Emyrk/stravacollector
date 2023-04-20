BEGIN;

ALTER TABLE athletes ADD COLUMN oauth_token_type text NOT NULL DEFAULT '';

COMMIT;