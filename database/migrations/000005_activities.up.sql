BEGIN;

CREATE TABLE activities (
	id bigint NOT NULL PRIMARY KEY,
	athlete_id bigint NOT NULL,
	upload_id bigint NOT NULL,
	external_id text NOT NULL,
	name text NOT NULL,
	moving_time double precision NOT NULL,
	elapsed_time double precision NOT NULL,
	total_elevation_gain double precision NOT NULL,
	activity_type text NOT NULL,
	sport_type text NOT NULL,
	start_date timestamp with time zone NOT NULL,
	start_date_local timestamp with time zone NOT NULL,
	timezone text NOT NULL,
	UTC_offset integer NOT NULL,
	start_latlng double precision[] NOT NULL,
	end_latlng double precision[] NOT NULL,
	achievement_count integer NOT NULL,
	kudos_count integer NOT NULL,
	comment_count integer NOT NULL,
	athlete_count integer NOT NULL,
	photo_count integer NOT NULL,
	map_id text NOT NULL,
	map_polyline text NOT NULL,
	map_summary_polyline text NOT NULL,
	trainer boolean NOT NULL,
	commute boolean NOT NULL,
	manual boolean NOT NULL,
	private boolean NOT NULL,
	flagged boolean NOT NULL,
	gear_id text NOT NULL,
	from_accepted_tag boolean NOT NULL,
	average_speed double precision NOT NULL,
	max_speed double precision NOT NULL,
	average_cadence double precision NOT NULL,
	average_temp double precision NOT NULL,
	average_watts double precision NOT NULL,
	weighted_average_watts double precision NOT NULL,
	kilojoules double precision NOT NULL,
	device_watts boolean NOT NULL,
	has_heartrate boolean NOT NULL,
	max_watts double precision NOT NULL,
	elev_high double precision NOT NULL,
	elev_low double precision NOT NULL,
	pr_count integer NOT NULL,
	total_photo_count integer NOT NULL,
-- 	has_kudoed boolean NOT NULL,
	workout_type integer NOT NULL,
	suffer_score integer NOT NULL,
-- 	description text NOT NULL,
	calories double precision NOT NULL,
	embed_token text NOT NULL,
	segment_leaderboard_opt_out boolean NOT NULL,
	leaderboard_opt_out boolean NOT NULL,
	num_segment_efforts integer NOT NULL,

	-- Custom
	premium_fetch boolean NOT NULL
);

COMMENT ON COLUMN activities.external_id IS 'External ID refers to external source of the activity.';
COMMENT ON COLUMN activities.premium_fetch IS 'Owner of the activity has premium account at the time of the fetch.';

DROP TABLE athlete_efforts;
CREATE TABLE segment_efforts (
	id bigint NOT NULL,
	athlete_id bigint NOT NULL,
	segment_id bigint NOT NULL,
	name text NOT NULL,
	elapsed_time double precision NOT NULL,
	moving_time double precision NOT NULL,
	start_date timestamp with time zone NOT NULL,
	start_date_local timestamp with time zone NOT NULL,
	distance double precision NOT NULL,
	start_index integer NOT NULL,
	end_index integer NOT NULL,
	device_watts boolean NOT NULL,
	average_watts double precision NOT NULL,
	kom_rank integer,
	pr_rank integer
);

COMMENT ON COLUMN segment_efforts.distance IS 'Distance is in meters';

COMMIT;
