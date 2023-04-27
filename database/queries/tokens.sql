-- name: InsertAPIToken :one
INSERT INTO
    api_tokens (
		created_at,
		updated_at,
		last_used_at,
		id,
		name,
		athlete_id,
		hashed_token,
		expires_at,
		lifetime_seconds
	)
VALUES (Now(), Now(), Now(), $1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: RenewToken :one
UPDATE api_tokens
SET
	updated_at = Now(),
	last_used_at = Now(),
	expires_at = @expires_at,
	lifetime_seconds = @lifetime_seconds
WHERE
    id = @id
RETURNING *;

-- name: DeleteToken :exec
DELETE FROM api_tokens
WHERE
	id = @id;

-- name: DeleteExpiredTokens :exec
DELETE FROM api_tokens
WHERE
	expires_at < Now();

-- name: GetToken :one
SELECT
	*
FROM
    api_tokens
WHERE
	id = $1;