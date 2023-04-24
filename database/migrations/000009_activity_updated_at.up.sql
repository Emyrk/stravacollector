BEGIN;

ALTER TABLE activities
ADD COLUMN updated_at timestamp with time zone;

ALTER TABLE segment_efforts
	ADD COLUMN updated_at timestamp with time zone;

COMMENT ON COLUMN activities.updated_at IS 'The time at which the activity was last updated by the collector';

-- ALTER TABLE activities DROP CONSTRAINT activities_athlete_id_fkey;
-- ALTER TABLE segment_efforts DROP CONSTRAINT activities_athletes_id_fk;

ALTER TABLE ONLY activities
	ADD CONSTRAINT activities_athletes_id_fk FOREIGN KEY (athlete_id) REFERENCES athletes(id);

ALTER TABLE ONLY segment_efforts
	ADD CONSTRAINT segment_efforts_athletes_id_fk FOREIGN KEY (athlete_id) REFERENCES athletes(id);

ALTER TABLE activities
	DROP COLUMN IF EXISTS num_efforts;

ALTER TABLE ONLY segment_efforts
	ADD CONSTRAINT segment_efforts_pk PRIMARY KEY (id);

COMMIT;