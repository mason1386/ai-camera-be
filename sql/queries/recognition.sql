-- name: CreateRecognitionLog :one
INSERT INTO recognition_logs (
    camera_id, identity_id, snapshot_url, face_crop_url, confidence, label, occurred_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7
) RETURNING *;

-- name: ListRecognitionLogs :many
SELECT * FROM recognition_logs
ORDER BY occurred_at DESC
LIMIT $1 OFFSET $2;
