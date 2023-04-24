-- name: UpsertMap :one
INSERT INTO
	maps(
	updated_at, id, polyline, summary_polyline
)
VALUES
	(Now(), $1, $2, $3)
ON CONFLICT
	(id)
	DO UPDATE SET
		updated_at = Now(),
		polyline = $2,
		summary_polyline = $3
RETURNING *;
