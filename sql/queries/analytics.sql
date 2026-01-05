-- name: CreateRecognitionLog :one
INSERT INTO recognition_logs (
    event_id, camera_id, identity_id, match_score, face_image_url, recognized_at
) VALUES (
    $1, $2, $3, $4, $5, $6
) RETURNING *;

-- name: ListRecognitionLogs :many
SELECT rl.*, i.full_name as identity_name, c.name as camera_name
FROM recognition_logs rl
JOIN identities i ON rl.identity_id = i.id
JOIN cameras c ON rl.camera_id = c.id
WHERE (rl.identity_id = $1 OR $1 IS NULL)
  AND (rl.camera_id = $2 OR $2 IS NULL)
  AND (rl.recognized_at >= $3 OR $3 IS NULL)
  AND (rl.recognized_at <= $4 OR $4 IS NULL)
ORDER BY rl.recognized_at DESC
LIMIT $5 OFFSET $6;

-- name: ListAttendanceRecords :many
SELECT ar.*, i.full_name as identity_name
FROM attendance_records ar
JOIN identities i ON ar.identity_id = i.id
WHERE (ar.identity_id = $1 OR $1 IS NULL)
  AND (ar.date >= $2 OR $2 IS NULL)
  AND (ar.date <= $3 OR $3 IS NULL)
  AND (ar.status = $4 OR $4 IS NULL)
ORDER BY ar.date DESC, ar.check_in DESC
LIMIT $5 OFFSET $6;

-- name: GetAttendanceStats :many
SELECT status, COUNT(*) as count
FROM attendance_records
WHERE date = $1
GROUP BY status;
