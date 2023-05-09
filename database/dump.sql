-- Code generated by 'make coderd/database/generate'. DO NOT EDIT.

CREATE TYPE activity_detail_source AS ENUM (
    'webhook',
    'backload',
    'requested',
    'manual',
    'unknown'
);

COMMENT ON TYPE activity_detail_source IS 'The source of the activity fetching.';

CREATE TABLE activity_detail (
    id bigint NOT NULL,
    athlete_id bigint NOT NULL,
    start_latlng double precision[] NOT NULL,
    end_latlng double precision[] NOT NULL,
    from_accepted_tag boolean NOT NULL,
    average_cadence double precision NOT NULL,
    average_temp double precision NOT NULL,
    average_watts double precision NOT NULL,
    weighted_average_watts double precision NOT NULL,
    kilojoules double precision NOT NULL,
    max_watts double precision NOT NULL,
    elev_high double precision NOT NULL,
    elev_low double precision NOT NULL,
    suffer_score integer NOT NULL,
    calories double precision NOT NULL,
    embed_token text NOT NULL,
    segment_leaderboard_opt_out boolean NOT NULL,
    leaderboard_opt_out boolean NOT NULL,
    num_segment_efforts integer NOT NULL,
    premium_fetch boolean NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    map_id text NOT NULL,
    source activity_detail_source DEFAULT 'unknown'::activity_detail_source NOT NULL
);

COMMENT ON COLUMN activity_detail.premium_fetch IS 'Owner of the activity has premium account at the time of the fetch.';

COMMENT ON COLUMN activity_detail.updated_at IS 'The time at which the activity was last updated by the collector';

CREATE TABLE activity_summary (
    id bigint NOT NULL,
    athlete_id bigint NOT NULL,
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
    updated_at timestamp with time zone NOT NULL,
    average_heartrate double precision DEFAULT 0 NOT NULL,
    max_heartrate double precision DEFAULT 0 NOT NULL
);

COMMENT ON TABLE activity_summary IS 'Activity is missing many detailed fields';

CREATE TABLE athletes (
    id bigint NOT NULL,
    summit boolean NOT NULL,
    username text NOT NULL,
    firstname text NOT NULL,
    lastname text NOT NULL,
    sex text NOT NULL,
    city text NOT NULL,
    state text NOT NULL,
    country text NOT NULL,
    follow_count integer NOT NULL,
    friend_count integer NOT NULL,
    measurement_preference text NOT NULL,
    ftp double precision NOT NULL,
    weight double precision NOT NULL,
    clubs json NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    fetched_at timestamp with time zone NOT NULL,
    profile_pic_link text DEFAULT ''::text NOT NULL,
    profile_pic_link_medium text DEFAULT ''::text NOT NULL
);

COMMENT ON COLUMN athletes.measurement_preference IS 'feet or meters';

CREATE TABLE competitive_routes (
    name text NOT NULL,
    display_name text NOT NULL,
    description text NOT NULL,
    segments bigint[] NOT NULL
);

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
    pr_rank integer,
    updated_at timestamp with time zone NOT NULL,
    activities_id bigint DEFAULT 0 NOT NULL
);

COMMENT ON COLUMN segment_efforts.distance IS 'Distance is in meters';

COMMENT ON COLUMN segment_efforts.activities_id IS 'FK to activities table';

CREATE VIEW hugel_activities AS
 SELECT merged.activity_id,
    merged.athlete_id,
    merged.segment_ids,
    merged.total_time_seconds,
    merged.efforts
   FROM ( SELECT hugel_efforts.activities_id AS activity_id,
            hugel_efforts.athlete_id,
            array_agg(hugel_efforts.segment_id) AS segment_ids,
            sum(hugel_efforts.elapsed_time) AS total_time_seconds,
            json_agg(json_build_object('activity_id', hugel_efforts.activities_id, 'effort_id', hugel_efforts.id, 'start_date', hugel_efforts.start_date, 'segment_id', hugel_efforts.segment_id, 'elapsed_time', hugel_efforts.elapsed_time, 'moving_time', hugel_efforts.moving_time, 'device_watts', hugel_efforts.device_watts, 'average_watts', hugel_efforts.average_watts)) AS efforts
           FROM ( SELECT DISTINCT ON (segment_efforts.activities_id, segment_efforts.segment_id) segment_efforts.id,
                    segment_efforts.athlete_id,
                    segment_efforts.segment_id,
                    segment_efforts.name,
                    segment_efforts.elapsed_time,
                    segment_efforts.moving_time,
                    segment_efforts.start_date,
                    segment_efforts.start_date_local,
                    segment_efforts.distance,
                    segment_efforts.start_index,
                    segment_efforts.end_index,
                    segment_efforts.device_watts,
                    segment_efforts.average_watts,
                    segment_efforts.kom_rank,
                    segment_efforts.pr_rank,
                    segment_efforts.updated_at,
                    segment_efforts.activities_id
                   FROM segment_efforts
                  WHERE (segment_efforts.segment_id = ANY (ARRAY( SELECT competitive_routes.segments
                           FROM competitive_routes
                          WHERE (competitive_routes.name = 'das-hugel'::text))))
                  ORDER BY segment_efforts.activities_id, segment_efforts.segment_id, segment_efforts.elapsed_time) hugel_efforts
          GROUP BY hugel_efforts.activities_id, hugel_efforts.athlete_id) merged
  WHERE (merged.segment_ids @> ARRAY( SELECT competitive_routes.segments
           FROM competitive_routes
          WHERE (competitive_routes.name = 'das-hugel'::text)));

COMMENT ON VIEW hugel_activities IS 'This view contains all activities that classify as a "hugel" and their best efforts on each segment.';

CREATE VIEW athlete_hugel_count AS
 SELECT hugel_activities.athlete_id,
    count(*) AS count
   FROM (public.athletes
     JOIN hugel_activities ON ((athletes.id = hugel_activities.athlete_id)))
  GROUP BY hugel_activities.athlete_id;

CREATE TABLE athlete_load (
    athlete_id bigint NOT NULL,
    last_backload_activity_start timestamp with time zone NOT NULL,
    last_load_attempt timestamp with time zone NOT NULL,
    last_load_incomplete boolean NOT NULL,
    last_load_error text NOT NULL,
    activites_loaded_last_attempt integer NOT NULL,
    earliest_activity timestamp with time zone DEFAULT (now())::timestamp without time zone NOT NULL,
    earliest_activity_done boolean DEFAULT false NOT NULL
);

COMMENT ON TABLE athlete_load IS 'Tracks loading athlete activities. Must be an authenticated athlete.';

COMMENT ON COLUMN athlete_load.last_backload_activity_start IS 'Timestamp start of the last activity loaded. Future ones are not loaded.';

COMMENT ON COLUMN athlete_load.last_load_attempt IS 'Timestamp of the last time the athlete was attempted to be loaded.';

COMMENT ON COLUMN athlete_load.last_load_incomplete IS 'True if the last load was incomplete and needs more work to catch up.';

COMMENT ON COLUMN athlete_load.earliest_activity IS 'The earliest activity found for the athlete';

COMMENT ON COLUMN athlete_load.earliest_activity_done IS 'Loading backwards is done';

CREATE TABLE athlete_logins (
    athlete_id bigint NOT NULL,
    summit boolean NOT NULL,
    provider_id text NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL,
    oauth_access_token text NOT NULL,
    oauth_refresh_token text NOT NULL,
    oauth_expiry timestamp with time zone NOT NULL,
    oauth_token_type text DEFAULT ''::text NOT NULL,
    id uuid DEFAULT gen_random_uuid() NOT NULL
);

COMMENT ON COLUMN athlete_logins.provider_id IS 'Oauth app client ID';

CREATE TABLE gue_jobs (
    job_id text NOT NULL,
    priority smallint NOT NULL,
    run_at timestamp with time zone NOT NULL,
    job_type text NOT NULL,
    args bytea NOT NULL,
    error_count integer DEFAULT 0 NOT NULL,
    last_error text,
    queue text NOT NULL,
    created_at timestamp with time zone NOT NULL,
    updated_at timestamp with time zone NOT NULL
);

CREATE TABLE maps (
    id text NOT NULL,
    polyline text NOT NULL,
    summary_polyline text NOT NULL,
    updated_at timestamp with time zone NOT NULL
);

CREATE TABLE segments (
    id bigint NOT NULL,
    name text NOT NULL,
    activity_type text NOT NULL,
    distance double precision NOT NULL,
    average_grade double precision NOT NULL,
    maximum_grade double precision NOT NULL,
    elevation_high double precision NOT NULL,
    elevation_low double precision NOT NULL,
    start_latlng double precision[] NOT NULL,
    end_latlng double precision[] NOT NULL,
    elevation_profile text NOT NULL,
    climb_category integer NOT NULL,
    city text NOT NULL,
    state text NOT NULL,
    country text NOT NULL,
    private boolean NOT NULL,
    hazardous boolean NOT NULL,
    created_at timestamp without time zone NOT NULL,
    updated_at timestamp without time zone NOT NULL,
    total_elevation_gain double precision NOT NULL,
    map_id text NOT NULL,
    total_effort_count integer NOT NULL,
    total_athlete_count integer NOT NULL,
    total_star_count integer NOT NULL,
    fetched_at timestamp without time zone NOT NULL
);

COMMENT ON COLUMN segments.elevation_profile IS 'A small image of the elevation profile of this segment.';

COMMENT ON COLUMN segments.fetched_at IS 'The time at which this segment was fetched from the Strava API.';

CREATE TABLE starred_segments (
    athlete_id bigint NOT NULL,
    segment_id bigint NOT NULL,
    starred boolean NOT NULL,
    updated_at timestamp with time zone NOT NULL
);

CREATE VIEW super_hugel_activities AS
 SELECT merged.athlete_id,
    merged.segment_ids,
    merged.total_time_seconds,
    merged.efforts
   FROM ( SELECT hugel_efforts.athlete_id,
            array_agg(hugel_efforts.segment_id) AS segment_ids,
            sum(hugel_efforts.elapsed_time) AS total_time_seconds,
            json_agg(json_build_object('activity_id', hugel_efforts.activities_id, 'effort_id', hugel_efforts.id, 'start_date', hugel_efforts.start_date, 'segment_id', hugel_efforts.segment_id, 'elapsed_time', hugel_efforts.elapsed_time, 'moving_time', hugel_efforts.moving_time, 'device_watts', hugel_efforts.device_watts, 'average_watts', hugel_efforts.average_watts)) AS efforts
           FROM ( SELECT DISTINCT ON (segment_efforts.athlete_id, segment_efforts.segment_id) segment_efforts.id,
                    segment_efforts.athlete_id,
                    segment_efforts.segment_id,
                    segment_efforts.name,
                    segment_efforts.elapsed_time,
                    segment_efforts.moving_time,
                    segment_efforts.start_date,
                    segment_efforts.start_date_local,
                    segment_efforts.distance,
                    segment_efforts.start_index,
                    segment_efforts.end_index,
                    segment_efforts.device_watts,
                    segment_efforts.average_watts,
                    segment_efforts.kom_rank,
                    segment_efforts.pr_rank,
                    segment_efforts.updated_at,
                    segment_efforts.activities_id
                   FROM segment_efforts
                  WHERE (segment_efforts.segment_id = ANY (ARRAY( SELECT competitive_routes.segments
                           FROM competitive_routes
                          WHERE (competitive_routes.name = 'das-hugel'::text))))
                  ORDER BY segment_efforts.athlete_id, segment_efforts.segment_id, segment_efforts.elapsed_time) hugel_efforts
          GROUP BY hugel_efforts.athlete_id) merged
  WHERE (merged.segment_ids @> ARRAY( SELECT competitive_routes.segments
           FROM competitive_routes
          WHERE (competitive_routes.name = 'das-hugel'::text)));

CREATE TABLE webhook_dump (
    id uuid NOT NULL,
    recorded_at timestamp without time zone NOT NULL,
    raw text NOT NULL
);

ALTER TABLE ONLY activity_detail
    ADD CONSTRAINT activities_pkey PRIMARY KEY (id);

ALTER TABLE ONLY activity_summary
    ADD CONSTRAINT activity_summary_pkey PRIMARY KEY (id);

ALTER TABLE ONLY athlete_load
    ADD CONSTRAINT athlete_load_pkey PRIMARY KEY (athlete_id);

ALTER TABLE ONLY athlete_logins
    ADD CONSTRAINT athletes_pkey PRIMARY KEY (athlete_id);

ALTER TABLE ONLY athletes
    ADD CONSTRAINT athletes_pkey1 PRIMARY KEY (id);

ALTER TABLE ONLY competitive_routes
    ADD CONSTRAINT competitive_routes_pkey PRIMARY KEY (name);

ALTER TABLE ONLY gue_jobs
    ADD CONSTRAINT gue_jobs_pkey PRIMARY KEY (job_id);

ALTER TABLE ONLY maps
    ADD CONSTRAINT maps_pkey PRIMARY KEY (id);

ALTER TABLE ONLY segment_efforts
    ADD CONSTRAINT segment_efforts_pk PRIMARY KEY (id);

ALTER TABLE ONLY segments
    ADD CONSTRAINT segments_pkey PRIMARY KEY (id);

ALTER TABLE ONLY starred_segments
    ADD CONSTRAINT starred_segments_pkey PRIMARY KEY (athlete_id, segment_id);

ALTER TABLE ONLY webhook_dump
    ADD CONSTRAINT webhook_dump_pkey PRIMARY KEY (id);

CREATE INDEX idx_gue_jobs_selector ON gue_jobs USING btree (queue, run_at, priority);

ALTER TABLE ONLY activity_detail
    ADD CONSTRAINT activities_athletes_id_fk FOREIGN KEY (athlete_id) REFERENCES athletes(id);

ALTER TABLE ONLY activity_detail
    ADD CONSTRAINT activity_detail_id_fk FOREIGN KEY (id) REFERENCES activity_summary(id) ON DELETE CASCADE;

ALTER TABLE ONLY activity_detail
    ADD CONSTRAINT activity_detail_map_id_fk FOREIGN KEY (map_id) REFERENCES maps(id);

ALTER TABLE ONLY activity_summary
    ADD CONSTRAINT activity_summary_athlete_id_fkey FOREIGN KEY (athlete_id) REFERENCES athletes(id) ON DELETE CASCADE;

ALTER TABLE ONLY athlete_load
    ADD CONSTRAINT athlete_load_athlete_id_fkey FOREIGN KEY (athlete_id) REFERENCES athlete_logins(athlete_id) ON DELETE CASCADE;

ALTER TABLE ONLY segment_efforts
    ADD CONSTRAINT segment_efforts_activities_id_fk FOREIGN KEY (activities_id) REFERENCES activity_detail(id) ON DELETE CASCADE;

ALTER TABLE ONLY segment_efforts
    ADD CONSTRAINT segment_efforts_athletes_id_fk FOREIGN KEY (athlete_id) REFERENCES athletes(id);

ALTER TABLE ONLY segments
    ADD CONSTRAINT segments_map_id_fkey FOREIGN KEY (map_id) REFERENCES maps(id);

ALTER TABLE ONLY starred_segments
    ADD CONSTRAINT starred_segments_athlete_id_fkey FOREIGN KEY (athlete_id) REFERENCES athletes(id);

