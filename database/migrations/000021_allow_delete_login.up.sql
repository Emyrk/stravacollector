BEGIN;

ALTER TABLE athlete_load
	DROP CONSTRAINT IF EXISTS athlete_load_athlete_id_fkey;

ALTER TABLE athlete_load
    ADD CONSTRAINT athlete_load_athlete_id_fkey FOREIGN KEY (athlete_id) REFERENCES athletes(id) ON DELETE CASCADE;

COMMIT;