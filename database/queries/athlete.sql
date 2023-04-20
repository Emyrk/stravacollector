-- name: GetAthletes :many
SELECT * FROM athletes;

-- name: GetAthlete :one
SELECT * FROM athletes WHERE id = $1;

-- name: UpsertAthlete :one
INSERT INTO
    athletes(
		created_at, updated_at,
             id,
             premium, username, firstname, lastname, sex,
             provider_id, oauth_access_token, oauth_refresh_token, oauth_expiry, oauth_token_type,
             raw
	)
VALUES
    (Now(), Now(), $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
ON CONFLICT
	(id)
DO UPDATE SET
	updated_at = Now(),
	premium = $2,
	username = $3,
	firstname = $4,
	lastname = $5,
	sex = $6,
	provider_id = $7,
	oauth_access_token = $8,
	oauth_refresh_token = $9,
	oauth_expiry = $10,
	oauth_token_type = $11,
	raw = $11
RETURNING *;


