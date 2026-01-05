-- name: CreateAuditLog :one
INSERT INTO audit_logs (
    user_id, action, table_name, record_id, old_value, new_value, ip_address, user_agent, created_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, NOW()
) RETURNING *;

-- name: ListAuditLogs :many
SELECT al.*, u.username
FROM audit_logs al
LEFT JOIN users u ON al.user_id = u.id
WHERE (al.user_id = $1 OR $1 IS NULL)
  AND (al.action = $2 OR $2 IS NULL)
  AND (al.table_name = $3 OR $3 IS NULL)
ORDER BY al.created_at DESC
LIMIT $4 OFFSET $5;
