-- name: CreateZone :one
INSERT INTO zones (
    name, description, created_at
) VALUES (
    $1, $2, NOW()
)
RETURNING id, created_at;

-- name: GetZone :one
SELECT * FROM zones
WHERE id = $1 LIMIT 1;

-- name: ListZones :many
SELECT * FROM zones
ORDER BY created_at DESC;

-- name: UpdateZone :one
UPDATE zones
SET name = $2, description = $3
WHERE id = $1
RETURNING *;

-- name: DeleteZone :exec
DELETE FROM zones
WHERE id = $1;
