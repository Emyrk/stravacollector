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