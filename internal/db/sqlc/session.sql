-- name: GetSessionById :one
SELECT * FROM session
WHERE id = $1 LIMIT 1;

-- name: GetSessionByUserId :one
SELECT * FROM session
WHERE user_id = $1 LIMIT 1;

-- name: CreateSession :one
INSERT INTO session (id, user_id, user_agent, client_ip, refresh_token, is_blocked, expires_at)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

