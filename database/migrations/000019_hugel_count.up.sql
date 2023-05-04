BEGIN;

CREATE OR REPLACE VIEW athlete_hugel_count AS
	SELECT
		athlete_id, count(*) AS count
	FROM
		athletes
			INNER JOIN
		hugel_activities
		ON athletes.id = hugel_activities.athlete_id
	GROUP BY athlete_id;

COMMIT;