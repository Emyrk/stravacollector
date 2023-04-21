-- name: GetAthleteLogin :one
SELECT * FROM athlete_logins WHERE athlete_id = @athlete_id;

-- name: UpsertAthleteLogin :one
INSERT INTO
	athlete_logins(
		created_at, updated_at, id,
             athlete_id, summit, provider_id, oauth_access_token,
             oauth_refresh_token, oauth_expiry, oauth_token_type
	)
VALUES
    (Now(), Now(), gen_random_uuid(), $1, $2, $3, $4, $5, $6, $7)
ON CONFLICT
	(athlete_id)
DO UPDATE SET
	updated_at = Now(),
	summit = $2,
	provider_id = $3,
	oauth_access_token = $4,
	oauth_refresh_token = $5,
	oauth_expiry = $6,
	oauth_token_type = $7
RETURNING *;

-- name: UpsertAthlete :one
INSERT INTO
	athletes(
	fetched_at, id, created_at, updated_at,
		summit, username, firstname, lastname, sex, city, state, country,
		follow_count, friend_count, measurement_preference, ftp, weight, clubs
)
VALUES
	(Now(), $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17)
ON CONFLICT
	(id)
	DO UPDATE SET
		fetched_at = Now(),
		created_at = $2,
		updated_at = $3,
		summit = $4,
		username = $5,
		firstname = $6,
		lastname = $7,
		sex = $8,
		city = $9,
		state = $10,
		country = $11,
		follow_count = $12,
		friend_count = $13,
		measurement_preference = $14,
		ftp = $15,
		weight = $16,
		clubs = $17

RETURNING *;


