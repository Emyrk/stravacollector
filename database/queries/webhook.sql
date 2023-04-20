-- name: InsertWebhookDump :exec
INSERT INTO
	webhook_dump(
	id, recorded_at, raw
)
VALUES
	(gen_random_uuid(), Now(), @raw_json);