-- name: CreateIdentity :one
INSERT INTO identities (
    code, full_name, group_name, phone_number, identity_card_number, department, metadata, status, note, created_by
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
) RETURNING *;

-- name: GetIdentity :one
SELECT * FROM identities
WHERE id = $1 AND deleted_at IS NULL LIMIT 1;

-- name: GetIdentityByCode :one
SELECT * FROM identities
WHERE code = $1 AND deleted_at IS NULL LIMIT 1;

-- name: ListIdentities :many
SELECT * FROM identities
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountIdentities :one
SELECT COUNT(*) FROM identities
WHERE deleted_at IS NULL;

-- name: UpdateIdentity :one
UPDATE identities
SET full_name = $2, group_name = $3, phone_number = $4, identity_card_number = $5, department = $6, metadata = $7, note = $8, updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: UpdateIdentityStatus :one
UPDATE identities
SET status = $2, updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: DeleteIdentity :exec
UPDATE identities
SET deleted_at = NOW(), updated_at = NOW()
WHERE id = $1;
