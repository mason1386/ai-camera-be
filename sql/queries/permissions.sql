-- name: GrantCameraPermission :exec
INSERT INTO user_camera_permissions (user_id, camera_id, created_at)
VALUES ($1, $2, NOW())
ON CONFLICT DO NOTHING;

-- name: RevokeCameraPermission :exec
DELETE FROM user_camera_permissions
WHERE user_id = $1 AND camera_id = $2;

-- name: ListUserCameraPermissions :many
SELECT camera_id FROM user_camera_permissions
WHERE user_id = $1;

-- name: GrantZonePermission :exec
INSERT INTO user_zone_permissions (user_id, zone_id, created_at)
VALUES ($1, $2, NOW())
ON CONFLICT DO NOTHING;

-- name: RevokeZonePermission :exec
DELETE FROM user_zone_permissions
WHERE user_id = $1 AND zone_id = $2;

-- name: ListUserZonePermissions :many
SELECT zone_id FROM user_zone_permissions
WHERE user_id = $1;
