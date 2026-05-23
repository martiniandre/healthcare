package auth

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidPassword = errors.New("invalid password")
	ErrUserExists      = errors.New("user already exists")
	ErrInvalidRole     = errors.New("invalid role")
	ErrUserInactive    = errors.New("user inactive")
)

type Service interface {
	Register(ctx context.Context, email, password, fullName, role string) (*User, error)
	Login(ctx context.Context, email, password string) (*User, string, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (authService *service) Register(ctx context.Context, email, password, fullName, role string) (*User, error) {
	existingUser, _ := authService.repo.GetUserByEmail(ctx, email)
	if existingUser != nil {
		return nil, ErrUserExists
	}

	parsedRole, roleIsValid := ParseRole(role)
	if !roleIsValid {
		return nil, ErrInvalidRole
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

	token, err := GenerateJWT(user.ID.String(), string(user.Role))
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}
