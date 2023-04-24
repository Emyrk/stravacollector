BEGIN;

ALTER TABLE segment_efforts
    ADD COLUMN activities_id BIGINT NOT NULL DEFAULT 0;

ALTER TABLE ONLY segment_efforts
	ADD CONSTRAINT segment_efforts_activities_id_fk FOREIGN KEY (activities_id) REFERENCES activities(id)
	ON DELETE CASCADE
;

COMMENT ON COLUMN segment_efforts.activities_id IS 'FK to activities table';

COMMIT;