BEGIN;

DROP TABLE segments;

CREATE TABLE segments(
    id BIGINT PRIMARY KEY NOT NULL,
    name TEXT NOT NULL,
    activity_type TEXT NOT NULL,
    distance DOUBLE PRECISION NOT NULL,
    average_grade DOUBLE PRECISION NOT NULL,
    maximum_grade DOUBLE PRECISION NOT NULL,
    elevation_high DOUBLE PRECISION NOT NULL,
    elevation_low DOUBLE PRECISION NOT NULL,
    start_latlng DOUBLE PRECISION[] NOT NULL,
    end_latlng DOUBLE PRECISION[] NOT NULL,
    elevation_profile TEXT NOT NULL,
    climb_category INTEGER NOT NULL,
    city TEXT NOT NULL,
    state TEXT NOT NULL,
    country TEXT NOT NULL,
    private BOOLEAN NOT NULL,
    hazardous BOOLEAN NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    total_elevation_gain DOUBLE PRECISION NOT NULL,
    map_id TEXT NOT NULL REFERENCES maps(id),
    total_effort_count INTEGER NOT NULL,
    total_athlete_count INTEGER NOT NULL,
    total_star_count INTEGER NOT NULL,

    --
    fetched_at TIMESTAMP NOT NULL
);


COMMENT ON COLUMN segments.elevation_profile IS 'A small image of the elevation profile of this segment.';
COMMENT ON COLUMN segments.fetched_at IS 'The time at which this segment was fetched from the Strava API.';
COMMIT;