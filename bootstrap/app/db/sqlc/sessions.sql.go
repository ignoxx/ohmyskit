// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: sessions.sql

package sqlc

import (
	"context"
	"database/sql"
	"time"
)

const createSession = `-- name: CreateSession :one
INSERT INTO sessions (
    user_id, token, expires_at, created_at, updated_at
) VALUES (
    ?1, ?2, ?3, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP
)
RETURNING id, token, user_id, ip_address, user_agent, expires_at, created_at, updated_at, deleted_at
`

type CreateSessionParams struct {
	UserID    int64
	Token     interface{}
	ExpiresAt time.Time
}

func (q *Queries) CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error) {
	row := q.db.QueryRowContext(ctx, createSession, arg.UserID, arg.Token, arg.ExpiresAt)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.Token,
		&i.UserID,
		&i.IpAddress,
		&i.UserAgent,
		&i.ExpiresAt,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
	)
	return i, err
}

const deleteSessionByToken = `-- name: DeleteSessionByToken :exec
DELETE FROM sessions WHERE token = ?1
`

func (q *Queries) DeleteSessionByToken(ctx context.Context, token interface{}) error {
	_, err := q.db.ExecContext(ctx, deleteSessionByToken, token)
	return err
}

const findSessionByTokenAndExpiration = `-- name: FindSessionByTokenAndExpiration :one
SELECT sessions.id, sessions.token, sessions.user_id, sessions.ip_address, sessions.user_agent, sessions.expires_at, sessions.created_at, sessions.updated_at, sessions.deleted_at, users.id as user_id, users.email as user_email
FROM sessions
JOIN users ON sessions.user_id = users.id
WHERE token = ?1 AND expires_at > CURRENT_TIMESTAMP
`

type FindSessionByTokenAndExpirationRow struct {
	ID        int64
	Token     interface{}
	UserID    int64
	IpAddress sql.NullString
	UserAgent sql.NullString
	ExpiresAt time.Time
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt sql.NullTime
	UserID_2  int64
	UserEmail string
}

// and also return the associated user
func (q *Queries) FindSessionByTokenAndExpiration(ctx context.Context, token interface{}) (FindSessionByTokenAndExpirationRow, error) {
	row := q.db.QueryRowContext(ctx, findSessionByTokenAndExpiration, token)
	var i FindSessionByTokenAndExpirationRow
	err := row.Scan(
		&i.ID,
		&i.Token,
		&i.UserID,
		&i.IpAddress,
		&i.UserAgent,
		&i.ExpiresAt,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.DeletedAt,
		&i.UserID_2,
		&i.UserEmail,
	)
	return i, err
}