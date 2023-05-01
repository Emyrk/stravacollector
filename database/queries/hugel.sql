-- name: HugelLeaderboard :many
SELECT
	ROW_NUMBER() over(ORDER BY sum ASC), *
FROM
	(
		SELECT DISTINCT ON (athlete_id)
			*
		FROM
			hugel_activities
		ORDER BY
			athlete_id, sum ASC
	) AS athlete_bests
ORDER BY sum ASC
;

-- name: GetCompetitiveRoute :many
SELECT * FROM competitive_routes WHERE name = @route_name;
