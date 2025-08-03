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
	id, distance, total_elevation_gain
FROM
	activity_summary
WHERE
	athlete_id = @athlete_id
  AND lower(activity_type) = ANY(ARRAY['ride', 'virtualride'])
;


-- name: GetAthleteEddington :one
SELECT
	*
FROM
	athlete_eddingtons
WHERE
	athlete_id = @athlete_id
;

-- name: AthletesNeedingEddington :many
SELECT
	athlete_logins.athlete_id, athlete_eddingtons.last_calculated
FROM
	athlete_logins
	LEFT JOIN
		athlete_eddingtons
		ON athlete_eddingtons.athlete_id = athlete_logins.athlete_id
WHERE
	athlete_eddingtons.last_calculated IS NULL -- null is never loaded
	OR athlete_eddingtons.last_calculated < (now() - interval '24hr')
;


-- name: AllEddingtons :many
SELECT
	athlete_eddingtons.athlete_id,
	athlete_eddingtons.current_eddington,
	athlete_eddingtons.total_activities
FROM
	athlete_eddingtons
;