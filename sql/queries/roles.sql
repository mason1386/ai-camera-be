-- name: CreateRole :one
INSERT INTO roles (
    name, description, permissions, is_system
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetRole :one
SELECT * FROM roles
WHERE id = $1 LIMIT 1;

-- name: ListRoles :many
SELECT * FROM roles
ORDER BY name ASC;

-- name: UpdateRole :one
UPDATE roles
SET name = $2, description = $3, permissions = $4, is_system = $5
WHERE id = $1
RETURNING *;

-- name: DeleteRole :exec
DELETE FROM roles
WHERE id = $1 AND is_system = FALSE;
