package auth

import (
	"time"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/shared/role"
)

type User struct {
	ID           uuid.UUID `db:"id"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	FullName     string    `db:"full_name"`
	Role         role.Role `db:"role"`
	IsActive     bool      `db:"is_active"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type RefreshToken struct {
	ID                uuid.UUID  `db:"id"`
	FamilyID          uuid.UUID  `db:"family_id"`
	UserID            uuid.UUID  `db:"user_id"`
	TokenHash         string     `db:"token_hash"`
	IssuedAt          time.Time  `db:"issued_at"`
	ExpiresAt         time.Time  `db:"expires_at"`
	AbsoluteExpiresAt time.Time  `db:"absolute_expires_at"`
	ReplacedBy        *uuid.UUID `db:"replaced_by"`
	RevokedAt         *time.Time `db:"revoked_at"`
	UserAgent         string     `db:"user_agent"`
	IP                string     `db:"ip"`
}

type ClientInfo struct {
	UserAgent string
	IP        string
}

type SessionTokens struct {
	AccessToken      string
	RefreshToken     string
	SessionExpiresAt time.Time
}
