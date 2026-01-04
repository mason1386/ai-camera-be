-- name: GetAIConfigByCamera :one
SELECT * FROM ai_configs
WHERE camera_id = $1 LIMIT 1;

-- name: CreateOrUpdateAIConfig :one
INSERT INTO ai_configs (
    camera_id, ai_enabled, ai_types, roi_zones, active_hours, sensitivity, min_confidence, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, NOW()
)
ON CONFLICT (camera_id) DO UPDATE SET
    ai_enabled = EXCLUDED.ai_enabled,
    ai_types = EXCLUDED.ai_types,
    roi_zones = EXCLUDED.roi_zones,
    active_hours = EXCLUDED.active_hours,
    sensitivity = EXCLUDED.sensitivity,
    min_confidence = EXCLUDED.min_confidence,
    updated_at = NOW()
RETURNING *;

-- name: CreateAIEvent :one
INSERT INTO ai_events (
    camera_id, event_type, confidence, snapshot_url, metadata, status, created_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: ListAIEvents :many
SELECT e.*, c.name as camera_name
FROM ai_events e
JOIN cameras c ON e.camera_id = c.id
WHERE (e.camera_id = $1 OR $1 IS NULL)
  AND (e.event_type = $2 OR $2 IS NULL)
  AND (e.status = $3 OR $3 IS NULL)
  AND (e.created_at >= $4 OR $4 IS NULL)
  AND (e.created_at <= $5 OR $5 IS NULL)
ORDER BY e.created_at DESC
LIMIT $6 OFFSET $7;

-- name: UpdateAIEventStatus :one
UPDATE ai_events
SET status = $2, resolved_by = $3, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: GetDashboardStats :one
SELECT 
    COUNT(id) as total_cameras,
    COUNT(CASE WHEN status = 'online' THEN 1 END) as online_cameras,
    COUNT(CASE WHEN status = 'offline' THEN 1 END) as offline_cameras,
    COUNT(CASE WHEN status = 'maintenance' THEN 1 END) as maintenance_cameras
FROM cameras;

-- name: GetTodayEventsCount :one
SELECT COUNT(*) FROM ai_events
WHERE created_at >= CURRENT_DATE;
