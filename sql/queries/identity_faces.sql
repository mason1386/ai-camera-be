-- name: CreateIdentityFace :one
INSERT INTO identity_faces (
    identity_id, image_url, is_primary, quality_score, blur_score
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: ListIdentityFaces :many
SELECT * FROM identity_faces
WHERE identity_id = $1;

-- name: DeleteIdentityFace :exec
DELETE FROM identity_faces
WHERE id = $1;

-- name: SetPrimaryFace :exec
UPDATE identity_faces
SET is_primary = (id = $2)
WHERE identity_id = $1;
