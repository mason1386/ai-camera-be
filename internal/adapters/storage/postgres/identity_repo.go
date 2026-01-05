package postgres

import (
	"context"
	"fmt"

	"app/internal/core/domain"
	"app/internal/core/ports"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type IdentityRepository struct {
	db *PostgresDB
}

func NewIdentityRepository(db *PostgresDB) ports.IdentityRepository {
	return &IdentityRepository{
		db: db,
	}
}

func (r *IdentityRepository) CreateIdentity(ctx context.Context, identity *domain.Identity) (*domain.Identity, error) {
	query := `
		INSERT INTO identities (
			code, full_name, type, phone_number, identity_card_number, face_image_url, department, metadata, status, note, created_by, approved_by, user_account_id
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
		) RETURNING id, created_at, updated_at`

	err := r.db.Pool.QueryRow(ctx, query,
		identity.Code,
		identity.FullName,
		identity.Type,
		identity.PhoneNumber,
		identity.IdentityCardNumber,
		identity.FaceImageURL,
		identity.Department,
		identity.Metadata,
		identity.Status,
		identity.Note,
		identity.CreatedBy,
		identity.ApprovedBy,
		identity.UserAccountID,
	).Scan(&identity.ID, &identity.CreatedAt, &identity.UpdatedAt)

	if err != nil {
		return nil, fmt.Errorf("failed to create identity: %w", err)
	}

	return identity, nil
}

func (r *IdentityRepository) GetIdentity(ctx context.Context, id uuid.UUID) (*domain.Identity, error) {
	query := `SELECT id, COALESCE(code, ''), COALESCE(full_name, ''), COALESCE(type, ''), phone_number, identity_card_number, face_image_url, COALESCE(department, ''), metadata, COALESCE(status, 'active'), note, created_by, approved_by, user_account_id, created_at, updated_at, deleted_at 
	          FROM identities WHERE id = $1 AND deleted_at IS NULL`

	identity := &domain.Identity{}
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&identity.ID, &identity.Code, &identity.FullName, &identity.Type,
		&identity.PhoneNumber, &identity.IdentityCardNumber, &identity.FaceImageURL, &identity.Department,
		&identity.Metadata, &identity.Status, &identity.Note, &identity.CreatedBy,
		&identity.ApprovedBy, &identity.UserAccountID,
		&identity.CreatedAt, &identity.UpdatedAt, &identity.DeletedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return identity, nil
}

func (r *IdentityRepository) GetIdentityByCode(ctx context.Context, code string) (*domain.Identity, error) {
	query := `SELECT id, COALESCE(code, ''), COALESCE(full_name, ''), COALESCE(type, ''), phone_number, identity_card_number, face_image_url, COALESCE(department, ''), metadata, COALESCE(status, 'active'), note, created_by, approved_by, user_account_id, created_at, updated_at, deleted_at 
	          FROM identities WHERE code = $1 AND deleted_at IS NULL`

	identity := &domain.Identity{}
	err := r.db.Pool.QueryRow(ctx, query, code).Scan(
		&identity.ID, &identity.Code, &identity.FullName, &identity.Type,
		&identity.PhoneNumber, &identity.IdentityCardNumber, &identity.FaceImageURL, &identity.Department,
		&identity.Metadata, &identity.Status, &identity.Note, &identity.CreatedBy,
		&identity.ApprovedBy, &identity.UserAccountID,
		&identity.CreatedAt, &identity.UpdatedAt, &identity.DeletedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return identity, nil
}

func (r *IdentityRepository) ListIdentities(ctx context.Context, page, limit int, search string) ([]*domain.Identity, int64, error) {
	offset := (page - 1) * limit

	whereClause := "WHERE deleted_at IS NULL"
	var args []interface{}
	argIdx := 1

	if search != "" {
		whereClause += fmt.Sprintf(" AND full_name ILIKE $%d", argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}

	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM identities %s", whereClause)
	var total int64
	err := r.db.Pool.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	query := fmt.Sprintf(`
		SELECT id, COALESCE(code, ''), COALESCE(full_name, ''), COALESCE(type, ''), COALESCE(department, ''), COALESCE(status, 'active'), created_at, updated_at
		FROM identities
		%s
		ORDER BY created_at DESC
		LIMIT $%d OFFSET $%d
	`, whereClause, argIdx, argIdx+1)

	args = append(args, limit, offset)

	rows, err := r.db.Pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	identities := []*domain.Identity{}
	for rows.Next() {
		identity := &domain.Identity{}
		err := rows.Scan(
			&identity.ID, &identity.Code, &identity.FullName, &identity.Type,
			&identity.Department, &identity.Status,
			&identity.CreatedAt, &identity.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		identities = append(identities, identity)
	}

	return identities, total, nil
}

func (r *IdentityRepository) CountIdentities(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.Pool.QueryRow(ctx, "SELECT COUNT(*) FROM identities WHERE deleted_at IS NULL").Scan(&count)
	return count, err
}

func (r *IdentityRepository) UpdateIdentity(ctx context.Context, identity *domain.Identity) (*domain.Identity, error) {
	query := `
		UPDATE identities
		SET full_name = $2, type = $3, phone_number = $4, identity_card_number = $5, face_image_url = $6, department = $7, metadata = $8, note = $9, updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
		RETURNING updated_at`

	err := r.db.Pool.QueryRow(ctx, query,
		identity.ID, identity.FullName, identity.Type, identity.PhoneNumber,
		identity.IdentityCardNumber, identity.FaceImageURL, identity.Department, identity.Metadata, identity.Note,
	).Scan(&identity.UpdatedAt)

	if err != nil {
		return nil, err
	}
	return identity, nil
}

func (r *IdentityRepository) UpdateIdentityStatus(ctx context.Context, id uuid.UUID, status domain.IdentityStatus) (*domain.Identity, error) {
	query := `
		UPDATE identities
		SET status = $2, updated_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL`

	_, err := r.db.Pool.Exec(ctx, query, id, status)
	if err != nil {
		return nil, err
	}
	return r.GetIdentity(ctx, id)
}

func (r *IdentityRepository) DeleteIdentity(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE identities SET deleted_at = NOW(), updated_at = NOW() WHERE id = $1`
	_, err := r.db.Pool.Exec(ctx, query, id)
	return err
}

type IdentityFaceRepository struct {
	db *PostgresDB
}

func NewIdentityFaceRepository(db *PostgresDB) ports.IdentityFaceRepository {
	return &IdentityFaceRepository{db: db}
}

func (r *IdentityFaceRepository) CreateFace(ctx context.Context, face *domain.IdentityFace) (*domain.IdentityFace, error) {
	query := `INSERT INTO identity_faces (identity_id, image_url, is_primary, quality_score, blur_score) 
	          VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at`
	err := r.db.Pool.QueryRow(ctx, query, face.IdentityID, face.ImageURL, face.IsPrimary, face.QualityScore, face.BlurScore).
		Scan(&face.ID, &face.CreatedAt)
	if err != nil {
		return nil, err
	}
	return face, nil
}

func (r *IdentityFaceRepository) ListFaces(ctx context.Context, identityID uuid.UUID) ([]*domain.IdentityFace, error) {
	query := `SELECT id, identity_id, image_url, is_primary, quality_score, blur_score, created_at 
	          FROM identity_faces WHERE identity_id = $1`
	rows, err := r.db.Pool.Query(ctx, query, identityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	faces := []*domain.IdentityFace{}
	for rows.Next() {
		face := &domain.IdentityFace{}
		err := rows.Scan(&face.ID, &face.IdentityID, &face.ImageURL, &face.IsPrimary, &face.QualityScore, &face.BlurScore, &face.CreatedAt)
		if err != nil {
			return nil, err
		}
		faces = append(faces, face)
	}
	return faces, nil
}

func (r *IdentityFaceRepository) DeleteFace(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.Pool.Exec(ctx, "DELETE FROM identity_faces WHERE id = $1", id)
	return err
}

func (r *IdentityFaceRepository) SetPrimary(ctx context.Context, identityID, faceID uuid.UUID) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, "UPDATE identity_faces SET is_primary = false WHERE identity_id = $1", identityID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "UPDATE identity_faces SET is_primary = true WHERE id = $1", faceID)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}
