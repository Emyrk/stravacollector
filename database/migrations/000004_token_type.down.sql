BEGIN;

ALTER TABLE athletes DROP COLUMN oauth_token_type;

COMMIT;