package auth

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository interface {
	CreateUser(ctx context.Context, user *User) error
	GetUserByEmail(ctx context.Context, email string) (*User, error)
	GetUserByID(ctx context.Context, userID string) (*User, error)
	Revoke(ctx context.Context, tokenDigest string, expiresAt time.Time) error
	IsRevoked(ctx context.Context, tokenDigest string) (bool, error)
}

type repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &repository{db: db}
}

func (authRepository *repository) CreateUser(ctx context.Context, user *User) error {
	query := `INSERT INTO users (id, email, password_hash, full_name, role, is_active, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	_, err := authRepository.db.Exec(ctx, query, user.ID, user.Email, user.PasswordHash, user.FullName, user.Role, user.IsActive, user.CreatedAt, user.UpdatedAt)
	return err
}

func (authRepository *repository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	query := `SELECT id, email, password_hash, full_name, role, is_active, created_at, updated_at FROM users WHERE email = $1`

	user := &User{}
	err := authRepository.db.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FullName,
		&user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (authRepository *repository) GetUserByID(ctx context.Context, userID string) (*User, error) {
	query := `SELECT id, email, password_hash, full_name, role, is_active, created_at, updated_at FROM users WHERE id = $1`

	user := &User{}
	err := authRepository.db.QueryRow(ctx, query, userID).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.FullName,
		&user.Role, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (authRepository *repository) Revoke(ctx context.Context, tokenDigest string, expiresAt time.Time) error {
	query := `INSERT INTO revoked_tokens (token_digest, expires_at) VALUES ($1, $2)
		ON CONFLICT (token_digest) DO UPDATE SET expires_at = EXCLUDED.expires_at`
	_, err := authRepository.db.Exec(ctx, query, tokenDigest, expiresAt)
	return err
}

func (authRepository *repository) IsRevoked(ctx context.Context, tokenDigest string) (bool, error) {
	query := `SELECT EXISTS(SELECT 1 FROM revoked_tokens WHERE token_digest = $1 AND expires_at > NOW())`
	var revoked bool
	err := authRepository.db.QueryRow(ctx, query, tokenDigest).Scan(&revoked)
	return revoked, err
}

func (authRepository *repository) CreatePatientLink(ctx context.Context, userID uuid.UUID) error {
	query := `INSERT INTO patient_user_links (user_id, patient_fhir_id) VALUES ($1, $1::text)
		ON CONFLICT (user_id) DO NOTHING`
	_, err := authRepository.db.Exec(ctx, query, userID)
	return err
}
