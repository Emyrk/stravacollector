-- name: UpsertMapData :one
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
		polyline =
		    CASE
				WHEN $2 != '' THEN $2
		        ELSE maps.polyline
			END,
		summary_polyline =
			CASE
				WHEN $3 != '' THEN $3
				ELSE maps.summary_polyline
			END
RETURNING *;

