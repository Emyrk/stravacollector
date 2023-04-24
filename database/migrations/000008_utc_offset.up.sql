BEGIN;

ALTER TABLE activities
ALTER COLUMN utc_offset
TYPE double precision USING utc_offset::double precision;

COMMIT;