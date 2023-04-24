BEGIN;

CREATE TABLE activity_summary(
	id bigint NOT NULL PRIMARY KEY,
	athlete_id bigint
	    REFERENCES athletes(id) ON DELETE CASCADE
	    NOT NULL,
	upload_id bigint NOT NULL,
	external_id text NOT NULL,

	name text NOT NULL,
	distance double precision NOT NULL,
	moving_time double precision NOT NULL,
	elapsed_time double precision NOT NULL,
	total_elevation_gain double precision NOT NULL,
	activity_type text NOT NULL,
	sport_type text NOT NULL,
	workout_type integer NOT NULL,
	start_date timestamp with time zone NOT NULL,
	start_date_local timestamp with time zone NOT NULL,
	timezone text NOT NULL,
	utc_offset double precision NOT NULL,
	achievement_count integer NOT NULL,
	kudos_count integer NOT NULL,
	comment_count integer NOT NULL,
	athlete_count integer NOT NULL,
	photo_count integer NOT NULL,
	map_id text NOT NULL,
	trainer boolean NOT NULL,
	commute boolean NOT NULL,
	manual boolean NOT NULL,
	private boolean NOT NULL,
	flagged boolean NOT NULL,
	gear_id text NOT NULL,
	average_speed double precision NOT NULL,
	max_speed double precision NOT NULL,
	device_watts boolean NOT NULL,
	has_heartrate boolean NOT NULL,
	pr_count integer NOT NULL,
	total_photo_count integer NOT NULL,
	calories double precision NOT NULL,
	updated_at timestamp with time zone NOT NULL
);

CREATE TABLE maps(
	id text NOT NULL PRIMARY KEY,
	polyline text NOT NULL,
	summary_polyline text NOT NULL,
	updated_at timestamp with time zone NOT NULL
);

ALTER TABLE activities
	DROP COLUMN name,
    DROP COLUMN upload_id,
    DROP COLUMN external_id,
-- 	DROP COLUMN distance,
	DROP COLUMN moving_time,
	DROP COLUMN elapsed_time,
	DROP COLUMN total_elevation_gain,
	DROP COLUMN activity_type,
	DROP COLUMN sport_type,
	DROP COLUMN workout_type,
	DROP COLUMN start_date,
	DROP COLUMN start_date_local,
	DROP COLUMN timezone,
	DROP COLUMN utc_offset,
	DROP COLUMN achievement_count,
	DROP COLUMN kudos_count,
	DROP COLUMN comment_count,
	DROP COLUMN athlete_count,
	DROP COLUMN photo_count,
	DROP COLUMN map_id,
	DROP COLUMN trainer,
	DROP COLUMN commute,
	DROP COLUMN manual,
	DROP COLUMN private,
	DROP COLUMN flagged,
	DROP COLUMN gear_id,
	DROP COLUMN average_speed,
	DROP COLUMN max_speed,
	DROP COLUMN device_watts,
	DROP COLUMN has_heartrate,
	DROP COLUMN pr_count,
	DROP COLUMN total_photo_count,
	DROP COLUMN calories,
    DROP COLUMN map_polyline,
    DROP COLUMN map_summary_polyline,
	ADD COLUMN map_id text NOT NULL
;

COMMENT ON TABLE activity_summary IS 'Activity is missing many detailed fields';

ALTER TABLE activities RENAME TO activity_detail;

ALTER TABLE ONLY activity_detail
	ADD CONSTRAINT activity_detail_id_fk FOREIGN KEY (id) REFERENCES activity_summary(id)
	ON DELETE CASCADE ;

ALTER TABLE ONLY activity_detail
	ADD CONSTRAINT activity_detail_map_id_fk FOREIGN KEY (map_id) REFERENCES maps(id);


COMMIT;