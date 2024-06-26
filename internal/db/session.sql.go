// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: session.sql

package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createSession = `-- name: CreateSession :one
INSERT INTO session (id, user_id, user_agent, client_ip, refresh_token, is_blocked, expires_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, user_id, user_agent, client_ip, refresh_token, is_blocked, created_at, expires_at
`

type CreateSessionParams struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	UserAgent    string    `json:"user_agent"`
	ClientIp     string    `json:"client_ip"`
	RefreshToken string    `json:"refresh_token"`
	IsBlocked    bool      `json:"is_blocked"`
	ExpiresAt    time.Time `json:"expires_at"`
}

func (q *Queries) CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error) {
	row := q.db.QueryRowContext(ctx, createSession,
		arg.ID,
		arg.UserID,
		arg.UserAgent,
		arg.ClientIp,
		arg.RefreshToken,
		arg.IsBlocked,
		arg.ExpiresAt,
	)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.UserAgent,
		&i.ClientIp,
		&i.RefreshToken,
		&i.IsBlocked,
		&i.CreatedAt,
		&i.ExpiresAt,
	)
	return i, err
}

const getSessionById = `-- name: GetSessionById :one
SELECT id, user_id, user_agent, client_ip, refresh_token, is_blocked, created_at, expires_at FROM session
WHERE id = $1 LIMIT 1
`

func (q *Queries) GetSessionById(ctx context.Context, id uuid.UUID) (Session, error) {
	row := q.db.QueryRowContext(ctx, getSessionById, id)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.UserAgent,
		&i.ClientIp,
		&i.RefreshToken,
		&i.IsBlocked,
		&i.CreatedAt,
		&i.ExpiresAt,
	)
	return i, err
}

const getSessionByUserId = `-- name: GetSessionByUserId :one
SELECT id, user_id, user_agent, client_ip, refresh_token, is_blocked, created_at, expires_at FROM session
WHERE user_id = $1 LIMIT 1
`

func (q *Queries) GetSessionByUserId(ctx context.Context, userID uuid.UUID) (Session, error) {
	row := q.db.QueryRowContext(ctx, getSessionByUserId, userID)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.UserAgent,
		&i.ClientIp,
		&i.RefreshToken,
		&i.IsBlocked,
		&i.CreatedAt,
		&i.ExpiresAt,
	)
	return i, err
}
