BEGIN;

CREATE TABLE athletes (
	id integer PRIMARY KEY
);

CREATE TABLE segments (
	id integer PRIMARY KEY,
	name text NOT NULL
);

CREATE TABLE athlete_efforts (
	 id uuid PRIMARY KEY,
	 athlete_id integer REFERENCES athletes(id) NOT NULL,
	 segment_id integer REFERENCES segments(id) NOT NULL,
	 last_checked timestamp NOT NULL,
	 best_effort_id integer NOT NULL,
	 kom_rank integer NOT NULL
);

COMMIT;