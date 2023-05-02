-- name: DeleteActivity :one
DELETE FROM
	activity_summary
WHERE
	id = $1
RETURNING *
;

-- name: GetActivityDetail :one
SELECT
	*
FROM
	activity_detail
WHERE
	id = $1;

-- name: GetActivitySummary :one
SELECT
	*
FROM
    activity_summary
WHERE
    id = $1;

-- name: UpdateActivityName :exec
UPDATE activity_summary
SET
    name = $2
WHERE
    id = $1;

-- name: UpsertActivitySummary :one
INSERT INTO
	activity_summary(
		updated_at, id, athlete_id, upload_id, external_id, name,
	    distance, moving_time, elapsed_time, total_elevation_gain,
	    activity_type, sport_type, workout_type, start_date,
	    start_date_local, timezone, utc_offset, achievement_count,
	    kudos_count, comment_count, athlete_count, photo_count, map_id,
	    trainer, commute, manual, private, flagged, gear_id, average_speed,
	    max_speed, device_watts, has_heartrate, pr_count, total_photo_count,
	    average_heartrate, max_heartrate
)
VALUES
	(Now(), $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20,
	 $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32, $33, $34, $35, $36)
ON CONFLICT
	(id)
	DO UPDATE SET
		updated_at = Now(),
		athlete_id = $2,
		upload_id = $3,
		external_id = $4,
		name = $5,
		distance = $6,
		moving_time = $7,
		elapsed_time = $8,
		total_elevation_gain = $9,
		activity_type = $10,
		sport_type = $11,
		workout_type = $12,
		start_date = $13,
		start_date_local = $14,
		timezone = $15,
		utc_offset = $16,
		achievement_count = $17,
		kudos_count = $18,
		comment_count = $19,
		athlete_count = $20,
		photo_count = $21,
		map_id = $22,
		trainer = $23,
		commute = $24,
		manual = $25,
		private = $26,
		flagged = $27,
		gear_id = $28,
		average_speed = $29,
		max_speed = $30,
		device_watts = $31,
		has_heartrate = $32,
		pr_count = $33,
		total_photo_count = $34,
		average_heartrate = $35,
		max_heartrate = $36
RETURNING *;

-- name: UpsertActivityDetail :one
INSERT INTO
	activity_detail(
		updated_at, id, athlete_id, start_latlng,
		end_latlng, from_accepted_tag, average_cadence, average_temp,
		average_watts, weighted_average_watts, kilojoules, max_watts,
	    elev_high, elev_low, suffer_score, embed_token,
	    segment_leaderboard_opt_out, leaderboard_opt_out, num_segment_efforts,
	    premium_fetch, map_id, calories, source
)
VALUES
	(Now(), $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17,
	 $18, $19, $20, $21, $22)
ON CONFLICT
	(id)
	DO UPDATE SET
	updated_at = Now(),
	athlete_id = $2,
	start_latlng = $3,
	end_latlng = $4,
	from_accepted_tag = $5,
	average_cadence = $6,
	average_temp = $7,
	average_watts = $8,
	weighted_average_watts = $9,
	kilojoules = $10,
	max_watts = $11,
	elev_high = $12,
	elev_low = $13,
	suffer_score = $14,
	embed_token = $15,
	segment_leaderboard_opt_out = $16,
	leaderboard_opt_out = $17,
	num_segment_efforts = $18,
	premium_fetch = $19,
	map_id = $20,
	calories = $21,
	source = $22
RETURNING *;