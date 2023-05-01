-- name: HugelLeaderboard :many
SELECT
	ROW_NUMBER() over(ORDER BY total_time_seconds ASC) AS rank,
	athlete_bests.activity_id,
	athlete_bests.athlete_id,
	athlete_bests.total_time_seconds,
	athlete_bests.efforts,

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
			hugel_activities
		ORDER BY
			athlete_id, total_time_seconds ASC
	) AS athlete_bests
INNER JOIN
	athletes ON athlete_bests.athlete_id = athletes.id
WHERE
    CASE WHEN @athlete_id > 0 THEN athlete_bests.athlete_id = @athlete_id ELSE TRUE END
ORDER BY
    total_time_seconds ASC
;

-- name: GetCompetitiveRoute :many
SELECT * FROM competitive_routes WHERE name = @route_name;
