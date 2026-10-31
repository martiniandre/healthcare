package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/shared/eventbus"
	"github.com/healthcare/backend/internal/shared/role"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")
	ErrUserExists      = errors.New("user already exists")
	ErrUserInactive    = errors.New("user inactive")

	ErrRefreshTokenInvalid = errors.New("refresh token invalid")
	ErrRefreshTokenReused  = errors.New("refresh token reused")
	ErrRefreshTokenExpired = errors.New("refresh token expired")
)

const (
	absoluteSessionLifetime = 12 * time.Hour
	refreshTokenLifetime    = 12 * time.Hour
	refreshTokenByteLength  = 32
)

type Service interface {
	Register(ctx context.Context, email, password, fullName, role string) (*User, error)
	Login(ctx context.Context, email, password string, clientInfo ClientInfo) (*User, *SessionTokens, error)
	Refresh(ctx context.Context, presentedRefreshToken string, clientInfo ClientInfo) (*User, *SessionTokens, error)
	Logout(ctx context.Context, presentedRefreshToken string) error
	Me(ctx context.Context, userID string) (*User, error)
}

type service struct {
	repo     Repository
	eventBus eventbus.Bus
}

func NewService(repo Repository, eventBus eventbus.Bus) Service {
	return &service{repo: repo, eventBus: eventBus}
}

func (authService *service) Register(ctx context.Context, email, password, fullName, requestedRole string) (*User, error) {
	existingUser, lookupError := authService.repo.GetUserByEmail(ctx, email)
	if lookupError != nil && !errors.Is(lookupError, ErrUserNotFound) {
		return nil, lookupError
	}
	if existingUser != nil {
		return nil, ErrUserExists
	}

	parsedRole, roleIsValid := role.ParseRole(requestedRole)
	if !roleIsValid {
		return nil, role.ErrInvalidRole
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: string(hashedPassword),
		FullName:     fullName,
		Role:         parsedRole,
		IsActive:     true,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if parsedRole == role.RolePatient {
		err = authService.repo.CreateUserWithPatientLink(ctx, user)
	} else {
		err = authService.repo.CreateUser(ctx, user)
	}
	if err != nil {
		return nil, err
	}

	return user, nil
}

var dummyPasswordHash, _ = bcrypt.GenerateFromPassword([]byte("__dummy_constant_time__"), bcrypt.DefaultCost)

func generateOpaqueRefreshToken() (string, error) {
	randomBytes := make([]byte, refreshTokenByteLength)
	if _, readError := rand.Read(randomBytes); readError != nil {
		return "", readError
	}
	return base64.RawURLEncoding.EncodeToString(randomBytes), nil
}

func HashRefreshToken(refreshToken string) string {
	refreshTokenDigest := sha256.Sum256([]byte(refreshToken))
	return hex.EncodeToString(refreshTokenDigest[:])
}

func buildRefreshTokenRecord(user *User, familyID uuid.UUID, absoluteExpiresAt time.Time, clientInfo ClientInfo) (*RefreshToken, string, error) {
	opaqueToken, generateError := generateOpaqueRefreshToken()
	if generateError != nil {
		return nil, "", generateError
	}

	issuedAt := time.Now()
	expiresAt := issuedAt.Add(refreshTokenLifetime)
	if expiresAt.After(absoluteExpiresAt) {
		expiresAt = absoluteExpiresAt
	}

	return &RefreshToken{
		ID:                uuid.New(),
		FamilyID:          familyID,
		UserID:            user.ID,
		TokenHash:         HashRefreshToken(opaqueToken),
		IssuedAt:          issuedAt,
		ExpiresAt:         expiresAt,
		AbsoluteExpiresAt: absoluteExpiresAt,
		UserAgent:         clientInfo.UserAgent,
		IP:                clientInfo.IP,
	}, opaqueToken, nil
}

func (authService *service) Login(ctx context.Context, email, password string, clientInfo ClientInfo) (*User, *SessionTokens, error) {
	user, err := authService.repo.GetUserByEmail(ctx, email)
	if err != nil {
		bcrypt.CompareHashAndPassword(dummyPasswordHash, []byte(password))
		return nil, nil, ErrUserNotFound
	}
	if !user.IsActive {
		bcrypt.CompareHashAndPassword(dummyPasswordHash, []byte(password))
		return nil, nil, ErrUserInactive
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, nil, ErrInvalidPassword
	}

	absoluteExpiresAt := time.Now().Add(absoluteSessionLifetime)
	refreshTokenRecord, opaqueRefreshToken, buildError := buildRefreshTokenRecord(user, uuid.New(), absoluteExpiresAt, clientInfo)
	if buildError != nil {
		return nil, nil, buildError
	}

	accessToken, err := GenerateJWTBoundedBy(user.ID.String(), string(user.Role), user.Email, absoluteExpiresAt)
	if err != nil {
		return nil, nil, err
	}

	if persistError := authService.repo.CreateRefreshToken(ctx, refreshTokenRecord); persistError != nil {
		return nil, nil, persistError
	}

	if authService.eventBus != nil {
		authService.eventBus.Publish(ctx, eventbus.Event{
			Name: "system.notification",
			Data: map[string]any{
				"title_key":     "notifications:event.loginSucceeded.title",
				"body_key":      "notifications:event.loginSucceeded.body",
				"params":        map[string]any{},
				"resource_type": "user",
				"resource_id":   user.ID.String(),
				"actor_id":      user.ID.String(),
			},
		})
	}

	return user, &SessionTokens{
		AccessToken:      accessToken,
		RefreshToken:     opaqueRefreshToken,
		SessionExpiresAt: absoluteExpiresAt,
	}, nil
}

func (authService *service) Refresh(ctx context.Context, presentedRefreshToken string, clientInfo ClientInfo) (*User, *SessionTokens, error) {
	if presentedRefreshToken == "" {
		return nil, nil, ErrRefreshTokenInvalid
	}

	storedToken, lookupError := authService.repo.GetRefreshTokenByHash(ctx, HashRefreshToken(presentedRefreshToken))
	if lookupError != nil {
		return nil, nil, lookupError
	}

	if storedToken.RevokedAt != nil || storedToken.ReplacedBy != nil {
		if revokeError := authService.repo.RevokeRefreshTokenFamily(ctx, storedToken.FamilyID); revokeError != nil {
			return nil, nil, revokeError
		}
		return nil, nil, ErrRefreshTokenReused
	}

	now := time.Now()
	if !now.Before(storedToken.ExpiresAt) || !now.Before(storedToken.AbsoluteExpiresAt) {
		return nil, nil, ErrRefreshTokenExpired
	}

	user, userError := authService.repo.GetUserByID(ctx, storedToken.UserID.String())
	if userError != nil {
		if errors.Is(userError, ErrUserNotFound) {
			return nil, nil, ErrRefreshTokenInvalid
		}
		return nil, nil, userError
	}
	if !user.IsActive {
		if revokeError := authService.repo.RevokeRefreshTokenFamily(ctx, storedToken.FamilyID); revokeError != nil {
			return nil, nil, revokeError
		}
		return nil, nil, ErrUserInactive
	}

	nextTokenRecord, nextOpaqueToken, buildError := buildRefreshTokenRecord(user, storedToken.FamilyID, storedToken.AbsoluteExpiresAt, clientInfo)
	if buildError != nil {
		return nil, nil, buildError
	}

	accessToken, generateError := GenerateJWTBoundedBy(user.ID.String(), string(user.Role), user.Email, storedToken.AbsoluteExpiresAt)
	if generateError != nil {
		return nil, nil, generateError
	}

	rotated, rotateError := authService.repo.RotateRefreshToken(ctx, storedToken.ID, nextTokenRecord)
	if rotateError != nil {
		return nil, nil, rotateError
	}
	if !rotated {
		if revokeError := authService.repo.RevokeRefreshTokenFamily(ctx, storedToken.FamilyID); revokeError != nil {
			return nil, nil, revokeError
		}
		return nil, nil, ErrRefreshTokenReused
	}

	return user, &SessionTokens{
		AccessToken:      accessToken,
		RefreshToken:     nextOpaqueToken,
		SessionExpiresAt: storedToken.AbsoluteExpiresAt,
	}, nil
}

func (authService *service) Logout(ctx context.Context, presentedRefreshToken string) error {
	if presentedRefreshToken == "" {
		return nil
	}

	storedToken, lookupError := authService.repo.GetRefreshTokenByHash(ctx, HashRefreshToken(presentedRefreshToken))
	if lookupError != nil {
		if errors.Is(lookupError, ErrRefreshTokenInvalid) {
			return nil
		}
		return lookupError
	}

	return authService.repo.RevokeRefreshTokenFamily(ctx, storedToken.FamilyID)
}

func (authService *service) Me(ctx context.Context, userID string) (*User, error) {
	return authService.repo.GetUserByID(ctx, userID)
}
