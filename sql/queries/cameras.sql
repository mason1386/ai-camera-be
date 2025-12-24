-- name: CreateCamera :one
INSERT INTO cameras (
    zone_id, name, ip_address, rtsp_url, status, ai_enabled, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, NOW(), NOW()
)
RETURNING id, created_at, updated_at;

-- name: GetCamera :one
SELECT * FROM cameras
WHERE id = $1 LIMIT 1;

-- name: ListCameras :many
SELECT * FROM cameras
ORDER BY created_at DESC;

-- name: ListCamerasByZone :many
SELECT * FROM cameras
WHERE zone_id = $1
ORDER BY created_at DESC;

-- name: UpdateCamera :one
UPDATE cameras
SET zone_id = $2, name = $3, ip_address = $4, rtsp_url = $5, status = $6, ai_enabled = $7, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: DeleteCamera :exec
DELETE FROM cameras
WHERE id = $1;
