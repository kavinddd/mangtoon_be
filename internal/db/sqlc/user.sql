-- name: CreateUser :one
INSERT INTO "user"(
	username, email, password
) VALUES (
	$1, $2, $3
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM "user"
WHERE id = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM "user"
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: ListActiveUsers :many
SELECT * FROM "user"
WHERE is_active = true
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateUserPassword :exec
UPDATE "user" SET password = $2 WHERE id = $1;

-- name: UpdateUserEmail :exec
UPDATE "user" SET email = $2 WHERE id = $1;

-- name: UpdateUserIsActive :exec
UPDATE "user" SET is_active = $2 WHERE id = $1;
