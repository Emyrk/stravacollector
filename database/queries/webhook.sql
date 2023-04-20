-- name: InsertWebhookDump :one
INSERT INTO
	webhook_dump(
	id, recorded_at, raw
)
VALUES
	(gen_random_uuid(), Now(), @raw_json)
RETURNING *
;