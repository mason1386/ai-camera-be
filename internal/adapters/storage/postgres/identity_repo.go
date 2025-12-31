package postgres

import (
	"context"

	"app/internal/adapters/storage/postgres/generated"
	"app/internal/core/domain"
	"app/internal/core/ports"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
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
	params := generated.CreateIdentityParams{
		Code:               identity.Code,
		FullName:           identity.FullName,
		PhoneNumber:        pgtype.Text{String: identity.PhoneNumber, Valid: identity.PhoneNumber != ""},
		IdentityCardNumber: pgtype.Text{String: identity.IdentityCardNumber, Valid: identity.IdentityCardNumber != ""},
		FaceImageUrl:       identity.FaceImageURL,
		Type:               identity.Type,
		Status:             generated.NullIdentityStatus{IdentityStatus: generated.IdentityStatus(identity.Status), Valid: true},
		Note:               pgtype.Text{String: identity.Note, Valid: identity.Note != ""},
	}

	if identity.CreatedBy != nil {
		params.CreatedBy = pgtype.UUID{Bytes: *identity.CreatedBy, Valid: true}
	}

	row, err := r.db.Query.CreateIdentity(ctx, params)
	if err != nil {
		return nil, err
	}

	return r.toDomain(row), nil
}

func (r *IdentityRepository) GetIdentity(ctx context.Context, id uuid.UUID) (*domain.Identity, error) {
	row, err := r.db.Query.GetIdentity(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return nil, err
	}
	return r.toDomain(row), nil
}

func (r *IdentityRepository) GetIdentityByCode(ctx context.Context, code string) (*domain.Identity, error) {
	row, err := r.db.Query.GetIdentityByCode(ctx, code)
	if err != nil {
		return nil, err
	}
	return r.toDomain(row), nil
}

func (r *IdentityRepository) ListIdentities(ctx context.Context, limit, offset int32) ([]*domain.Identity, error) {
	rows, err := r.db.Query.ListIdentities(ctx, generated.ListIdentitiesParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	identities := make([]*domain.Identity, len(rows))
	for i, row := range rows {
		identities[i] = r.toDomain(row)
	}
	return identities, nil
}

func (r *IdentityRepository) CountIdentities(ctx context.Context) (int64, error) {
	return r.db.Query.CountIdentities(ctx)
}

func (r *IdentityRepository) UpdateIdentity(ctx context.Context, identity *domain.Identity) (*domain.Identity, error) {
	params := generated.UpdateIdentityParams{
		ID:           pgtype.UUID{Bytes: identity.ID, Valid: true},
		FullName:     identity.FullName,
		PhoneNumber:  pgtype.Text{String: identity.PhoneNumber, Valid: identity.PhoneNumber != ""},
		FaceImageUrl: identity.FaceImageURL,
		// Status is updated via separate method
	}

	row, err := r.db.Query.UpdateIdentity(ctx, params)
	if err != nil {
		return nil, err
	}
	return r.toDomain(row), nil
}

func (r *IdentityRepository) UpdateIdentityStatus(ctx context.Context, id uuid.UUID, status domain.IdentityStatus, approvedBy *uuid.UUID) (*domain.Identity, error) {
	params := generated.UpdateIdentityStatusParams{
		ID:     pgtype.UUID{Bytes: id, Valid: true},
		Status: generated.NullIdentityStatus{IdentityStatus: generated.IdentityStatus(status), Valid: true},
	}
	if approvedBy != nil {
		params.ApprovedBy = pgtype.UUID{Bytes: *approvedBy, Valid: true}
	}

	row, err := r.db.Query.UpdateIdentityStatus(ctx, params)
	if err != nil {
		return nil, err
	}
	return r.toDomain(row), nil
}

func (r *IdentityRepository) DeleteIdentity(ctx context.Context, id uuid.UUID) error {
	return r.db.Query.DeleteIdentity(ctx, pgtype.UUID{Bytes: id, Valid: true})
}

func (r *IdentityRepository) toDomain(row generated.Identity) *domain.Identity {
	// Handle Null Status safely
	status := domain.IdentityStatusPending
	if row.Status.Valid {
		status = domain.IdentityStatus(row.Status.IdentityStatus)
	}

	id := &domain.Identity{
		ID:                 row.ID.Bytes,
		Code:               row.Code,
		FullName:           row.FullName,
		PhoneNumber:        row.PhoneNumber.String,
		IdentityCardNumber: row.IdentityCardNumber.String,
		FaceImageURL:       row.FaceImageUrl,
		Type:               row.Type,
		Status:             status,
		Note:               row.Note.String,
		CreatedAt:          row.CreatedAt.Time,
		UpdatedAt:          row.UpdatedAt.Time,
	}

	if row.CreatedBy.Valid {
		uid := uuid.UUID(row.CreatedBy.Bytes)
		id.CreatedBy = &uid
	}
	if row.ApprovedBy.Valid {
		uid := uuid.UUID(row.ApprovedBy.Bytes)
		id.ApprovedBy = &uid
	}
	if row.UserAccountID.Valid {
		uid := uuid.UUID(row.UserAccountID.Bytes)
		id.UserAccountID = &uid
	}
	if row.DeletedAt.Valid {
		t := row.DeletedAt.Time
		id.DeletedAt = &t
	}

	return id
}
