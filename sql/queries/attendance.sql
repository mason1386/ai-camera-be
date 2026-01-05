-- name: CreateAttendanceRecord :one
INSERT INTO attendance_records (
    identity_id, date, check_in, status
) VALUES (
    $1, $2, $3, $4
) RETURNING *;

-- name: GetAttendanceRecord :one
SELECT * FROM attendance_records
WHERE identity_id = $1 AND date = $2
LIMIT 1;

-- name: UpdateAttendanceRecord :one
UPDATE attendance_records
SET check_out = $2, work_hours = $3, status = $4, updated_at = NOW()
WHERE id = $1
RETURNING *;
