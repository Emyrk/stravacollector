SELECT * FROM hugel_leaderboard(true, 0);

CREATE OR REPLACE FUNCTION hugel_leaderboard (super BOOLEAN, filter_athlete_id BIGINT)
	RETURNS TABLE (
					  best_time bigint,
					  rank bigint,
					  activity_id bigint,
					  athlete_id bigint,
					  total_time_seconds double precision,
					  efforts json,
					  name text,
					  distance double precision,
					  moving_time double precision,
					  elapsed_time double precision,
					  total_elevation_gain double precision,
					  start_date timestamp with time zone,
					  firstname text,
					  lastname text,
					  username text,
					  profile_pic_link text,
					  sex text
				  )
AS $$
DECLARE selected_activites hugel_activities;
--         RECORD(
--             segment_ids BIGINT[],
--             athlete_id BIGINT,
--             activity_id BIGINT,
--             total_time_seconds DOUBLE PRECISION,
--             efforts jsonb
--         );
BEGIN
	CREATE TEMP TABLE selected_activites AS SELECT * FROM hugel_activities;
--     SELECT * INTO selected_activites FROM hugel_activities;

	RETURN QUERY SELECT
					 (SELECT min(total_time_seconds) FROM selected_activites) :: BIGINT AS best_time,
					 ROW_NUMBER() over(ORDER BY total_time_seconds ASC) AS rank,
					 athlete_bests.activity_id,
					 athlete_bests.athlete_id,
					 athlete_bests.total_time_seconds,
					 athlete_bests.efforts,

					 activity_summary.name,
					 activity_summary.distance,
					 activity_summary.moving_time,
					 activity_summary.elapsed_time,
					 activity_summary.total_elevation_gain,
					 activity_summary.start_date,

					 athletes.firstname,
					 athletes.lastname,
					 athletes.username,
					 athletes.profile_pic_link,
					 athletes.sex
				 FROM
					 (
						 SELECT DISTINCT ON (athlete_id)
							 *
						 FROM
							 selected_activites
						 ORDER BY
							 athlete_id, total_time_seconds ASC
					 ) AS athlete_bests
						 INNER JOIN
					 athletes ON athlete_bests.athlete_id = athletes.id
						 INNER JOIN
					 activity_summary ON athlete_bests.activity_id = activity_summary.id
				 WHERE
					 CASE WHEN filter_athlete_id > 0 THEN athlete_bests.athlete_id = filter_athlete_id ELSE TRUE END
				 ORDER BY
					 total_time_seconds ASC;
END; $$

	LANGUAGE 'plpgsql';
