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

-- name: UpsertMapSummary :one
INSERT INTO
	maps(
	updated_at, polyline, id, summary_polyline
)
VALUES
	(Now(), '', $1, $2)
ON CONFLICT
	(id)
	DO UPDATE SET
	  updated_at = Now(),
	  summary_polyline = $2
RETURNING *;

