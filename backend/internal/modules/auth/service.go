package auth

import (
	"context"
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
)

type Service interface {
	Register(ctx context.Context, email, password, fullName, role string) (*User, error)
	Login(ctx context.Context, email, password string) (*User, string, error)
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
	existingUser, _ := authService.repo.GetUserByEmail(ctx, email)
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

	err = authService.repo.CreateUser(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

var dummyPasswordHash, _ = bcrypt.GenerateFromPassword([]byte("__dummy_constant_time__"), bcrypt.DefaultCost)

func (authService *service) Login(ctx context.Context, email, password string) (*User, string, error) {
	user, err := authService.repo.GetUserByEmail(ctx, email)
	if err != nil {
		bcrypt.CompareHashAndPassword(dummyPasswordHash, []byte(password))
		return nil, "", ErrUserNotFound
	}
	if !user.IsActive {
		bcrypt.CompareHashAndPassword(dummyPasswordHash, []byte(password))
		return nil, "", ErrUserInactive
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, "", ErrInvalidPassword
	}

	token, err := GenerateJWT(user.ID.String(), string(user.Role), user.Email)
	if err != nil {
		return nil, "", err
	}

	if authService.eventBus != nil {
		authService.eventBus.Publish(ctx, eventbus.Event{
			Name: "system.notification",
			Data: map[string]any{
				"title":       "Login Realizado",
				"body":        "Login realizado com sucesso como " + string(user.Role),
				"resource_id": user.ID.String(),
				"actor_id":    user.ID.String(),
			},
		})
	}

	return user, token, nil
}

func (authService *service) Me(ctx context.Context, userID string) (*User, error) {
	return authService.repo.GetUserByID(ctx, userID)
}
