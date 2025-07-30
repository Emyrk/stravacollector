-- name: UpsertAthleteEddington :one
INSERT INTO
	athlete_eddingtons(
		athlete_id,
		miles_histogram,
		current_eddington,
		last_calculated,
		total_activities
)
VALUES
	($1, $2, $3, $4, $5)
ON CONFLICT
	(athlete_id)
	DO UPDATE SET
		miles_histogram = $2,
		current_eddington = $3,
		last_calculated = $4,
		total_activities = $5
RETURNING *;

-- name: EddingtonActivities :many
SELECT
	distance, total_elevation_gain
FROM
	activity_summary
WHERE
	athlete_id = @athlete_id
  AND lower(activity_type) = 'ride'
;


-- name: GetAthleteEddington :one
SELECT
	*
FROM
	athlete_eddingtons
WHERE
	athlete_id = @athlete_id
;