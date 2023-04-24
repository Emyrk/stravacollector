-- name: UpsertSegmentEffort :one
INSERT INTO
	segment_efforts(
		updated_at,
		id, athlete_id, segment_id, name, elapsed_time,
		moving_time, start_date, start_date_local, distance,
		start_index, end_index, device_watts, average_watts,
		kom_rank, pr_rank, activities_id
	)
VALUES
	(Now(), $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
ON CONFLICT
	(id)
	DO UPDATE SET
		updated_at = Now(),
		athlete_id = $2,
		segment_id = $3,
		name = $4,
		elapsed_time = $5,
		moving_time = $6,
		start_date = $7,
		start_date_local = $8,
		distance = $9,
		start_index = $10,
		end_index = $11,
		device_watts = $12,
		average_watts = $13,
		kom_rank = $14,
		pr_rank = $15,
		activities_id = $16
	RETURNING *;