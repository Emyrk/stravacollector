BEGIN;

CREATE TYPE activity_detail_source AS ENUM ('webhook', 'backload', 'requested', 'manual', 'unknown');
COMMENT ON TYPE activity_detail_source IS 'The source of the activity fetching.';

ALTER TABLE activity_detail
    ADD COLUMN source activity_detail_source NOT NULL default 'unknown';

UPDATE activity_detail SET source = 'backload';

COMMIT;