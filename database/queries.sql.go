// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2

package database

import (
	"context"
	"time"
)

const getAthletes = `-- name: GetAthletes :many
SELECT id, premium, username, firstname, lastname, sex, provider_id, created_at, updated_at, oauth_access_token, oauth_refresh_token, oauth_expiry, raw FROM athletes
`

func (q *sqlQuerier) GetAthletes(ctx context.Context) ([]Athlete, error) {
	rows, err := q.db.QueryContext(ctx, getAthletes)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Athlete
	for rows.Next() {
		var i Athlete
		if err := rows.Scan(
			&i.ID,
			&i.Premium,
			&i.Username,
			&i.Firstname,
			&i.Lastname,
			&i.Sex,
			&i.ProviderID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.OauthAccessToken,
			&i.OauthRefreshToken,
			&i.OauthExpiry,
			&i.Raw,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const upsertAthlete = `-- name: UpsertAthlete :one
INSERT INTO
    athletes(
		created_at, updated_at,
             id,
             premium, username, firstname, lastname, sex,
             provider_id, oauth_access_token, oauth_refresh_token, oauth_expiry,
             raw
	)
VALUES
    (Now(), Now(), $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
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
	raw = $11
RETURNING id, premium, username, firstname, lastname, sex, provider_id, created_at, updated_at, oauth_access_token, oauth_refresh_token, oauth_expiry, raw
`

type UpsertAthleteParams struct {
	ID                int64     `db:"id" json:"id"`
	Premium           bool      `db:"premium" json:"premium"`
	Username          string    `db:"username" json:"username"`
	Firstname         string    `db:"firstname" json:"firstname"`
	Lastname          string    `db:"lastname" json:"lastname"`
	Sex               string    `db:"sex" json:"sex"`
	ProviderID        string    `db:"provider_id" json:"provider_id"`
	OauthAccessToken  string    `db:"oauth_access_token" json:"oauth_access_token"`
	OauthRefreshToken string    `db:"oauth_refresh_token" json:"oauth_refresh_token"`
	OauthExpiry       time.Time `db:"oauth_expiry" json:"oauth_expiry"`
	Raw               string    `db:"raw" json:"raw"`
}

func (q *sqlQuerier) UpsertAthlete(ctx context.Context, arg UpsertAthleteParams) (Athlete, error) {
	row := q.db.QueryRowContext(ctx, upsertAthlete,
		arg.ID,
		arg.Premium,
		arg.Username,
		arg.Firstname,
		arg.Lastname,
		arg.Sex,
		arg.ProviderID,
		arg.OauthAccessToken,
		arg.OauthRefreshToken,
		arg.OauthExpiry,
		arg.Raw,
	)
	var i Athlete
	err := row.Scan(
		&i.ID,
		&i.Premium,
		&i.Username,
		&i.Firstname,
		&i.Lastname,
		&i.Sex,
		&i.ProviderID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.OauthAccessToken,
		&i.OauthRefreshToken,
		&i.OauthExpiry,
		&i.Raw,
	)
	return i, err
}

const insertWebhookDump = `-- name: InsertWebhookDump :one
INSERT INTO
	webhook_dump(
	id, recorded_at, raw
)
VALUES
	(gen_random_uuid(), Now(), $1)
RETURNING id, recorded_at, raw
`

func (q *sqlQuerier) InsertWebhookDump(ctx context.Context, rawJson string) (WebhookDump, error) {
	row := q.db.QueryRowContext(ctx, insertWebhookDump, rawJson)
	var i WebhookDump
	err := row.Scan(&i.ID, &i.RecordedAt, &i.Raw)
	return i, err
}
