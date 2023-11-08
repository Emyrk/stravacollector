-- name: InsertFailedJob :one
INSERT INTO
	failed_jobs(
	id, recorded_at, raw
)
VALUES
	(gen_random_uuid(), Now(), @raw_json)
RETURNING *
;

-- name: TotalJobCount :one
SELECT count(*) FROM gue_jobs;