-- name: GetAthleteLoad :one
SELECT * FROM athlete_forward_load WHERE athlete_id = @athlete_id;

-- name: UpsertAthleteForwardLoad :one
INSERT INTO
	athlete_forward_load(
		athlete_id,
		activity_time_after,
		last_load_complete,
		last_touched,
		next_load_not_before
)
VALUES
	($1, $2, $3, $4, $5)
ON CONFLICT
	(athlete_id)
	DO UPDATE SET
		athlete_id = $1,
		activity_time_after = $2,
		last_load_complete = $3,
		last_touched = $4,
		next_load_not_before = $5
RETURNING *;
;

-- name: GetAthleteLoadDetailed :one
SELECT
    sqlc.embed(athlete_load),
    sqlc.embed(athletes),
	(SELECT count(*) FROM activity_summary WHERE activity_summary.athlete_id = @athlete_id AND LOWER(activity_summary.activity_type) = 'ride') AS summary_count,
    (SELECT count(*) FROM activity_detail WHERE activity_detail.athlete_id = @athlete_id AND activity_detail.id = ANY(
		SELECT id FROM activity_summary WHERE activity_summary.athlete_id = @athlete_id AND LOWER(activity_summary.activity_type) = 'ride')
	) AS detail_count,
	COALESCE(athlete_hugel_count.count, 0) AS hugel_count
FROM
    athlete_load
INNER JOIN
    athletes ON athletes.id = athlete_load.athlete_id
LEFT JOIN
	athlete_hugel_count ON athlete_hugel_count.athlete_id = athletes.id
WHERE
		athlete_load.athlete_id = @athlete_id;

-- name: AthleteSyncedActivities :many
SELECT
	sqlc.embed(activity_summary),
	COUNT(*) OVER() AS total,
	activity_detail.id IS NOT NULL :: boolean AS detail_exists,
	activity_detail.updated_at AS detail_updated_at
FROM
	activity_summary
LEFT JOIN
	activity_detail ON
		activity_summary.id = activity_detail.id
WHERE
	activity_summary.athlete_id = @athlete_id AND
	LOWER(activity_summary.activity_type) = 'ride'
ORDER BY
    activity_summary.start_date DESC
LIMIT @_limit
OFFSET @_offset
;

-- name: UpsertAthleteLoad :one
INSERT INTO
	athlete_load(
		athlete_id,
		last_backload_activity_start,
	    last_load_attempt,
		last_load_incomplete,
		last_load_error,
		activites_loaded_last_attempt,
		earliest_activity,
	    earliest_activity_id,
		earliest_activity_done,
		next_load_not_before
	)
VALUES
	($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
ON CONFLICT
	(athlete_id)
DO UPDATE SET
	last_backload_activity_start = $2,
	last_load_attempt = $3,
	last_load_incomplete = $4,
	last_load_error = $5,
	activites_loaded_last_attempt = $6,
	earliest_activity = $7,
	earliest_activity_id = $8,
	earliest_activity_done = $9,
    next_load_not_before = $10
RETURNING *;
;

-- name: GetAthleteNeedsForwardLoad :many
SELECT
	sqlc.embed(athlete_forward_load), sqlc.embed(athlete_logins)
FROM
	athlete_forward_load
INNER JOIN
	athlete_logins
	ON
		-- Ignore non-authed athletes
		athlete_load.athlete_id = athlete_logins.athlete_id
WHERE
	Now() > athlete_load.next_load_not_before
ORDER BY
	-- Athletes with oldest load attempt first.
	-- Order is [false, true].
	last_load_complete, activity_time_after, last_touched
LIMIT 5;

-- name: GetAthleteNeedsLoad :many
SELECT
    sqlc.embed(athlete_load), sqlc.embed(athlete_logins)
FROM
	athlete_load
INNER JOIN
	athlete_logins
ON
    athlete_load.athlete_id = athlete_logins.athlete_id
WHERE
    athlete_load.next_load_not_before < Now()
ORDER BY
	-- Athletes with oldest load attempt first.
	-- Order is [false, true]. 
	not last_load_incomplete, earliest_activity_done, last_load_attempt
LIMIT 5;

-- name: GetAthleteLogin :one
SELECT * FROM athlete_logins WHERE athlete_id = @athlete_id;

-- name: DeleteAthleteLogin :exec
DELETE FROM athlete_logins WHERE athlete_id = @athlete_id;


-- name: GetAthleteFull :one
SELECT
	sqlc.embed(athletes),
	COALESCE(athlete_hugel_count.count, 0) AS hugel_count
FROM
	athletes
	LEFT JOIN
	athlete_hugel_count ON athlete_hugel_count.athlete_id = athletes.id
WHERE
	athletes.id = @athlete_id;

-- name: GetAthleteLoginFull :one
SELECT
    sqlc.embed(athlete_logins),
    sqlc.embed(athletes),
    COALESCE(athlete_hugel_count.count, 0) AS hugel_count
FROM
    athlete_logins
INNER JOIN
    athletes ON athlete_logins.athlete_id = athletes.id
LEFT JOIN
	athlete_hugel_count ON athlete_hugel_count.athlete_id = athletes.id
WHERE
	athlete_logins.athlete_id = @athlete_id;

-- name: GetAthlete :one
SELECT * FROM athletes WHERE id = @athlete_id;

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
		follow_count, friend_count, measurement_preference, ftp, weight, clubs,
		profile_pic_link, profile_pic_link_medium
)
VALUES
	(Now(), $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19)
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
		clubs = $17,
		profile_pic_link = $18,
		profile_pic_link_medium = $19
RETURNING *;


