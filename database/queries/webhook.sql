-- name: InsertWebhookDump :one
INSERT INTO
	webhook_dump(
	id, recorded_at, raw
)
VALUES
	(gen_random_uuid(), Now(), @raw_json)
RETURNING *
;

-- name: GetDeleteActivityWebhooks :many
SELECT * FROM webhook_dump
WHERE
	raw::json ->> 'aspect_type' = 'delete'
  	AND raw::json ->> 'object_type' = 'activity'
  	AND raw::json ->> 'object_id' IN(
  	    SELECT id::text FROM activity_summary
	)
LIMIT 200
;

-- name: DeleteWebhookDump :exec
DELETE FROM webhook_dump
WHERE
	id = @id
;