ALTER TABLE segments ADD COLUMN friendly_name TEXT NOT NULL DEFAULT '';

COMMENT ON COLUMN segments.friendly_name IS 'Human friendly name for the segment';