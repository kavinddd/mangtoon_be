-- name: ListRolesByUserId :many
SELECT r.name
FROM user_role ur
JOIN role r ON ur.role_id = r.id
WHERE ur.user_id = $1;