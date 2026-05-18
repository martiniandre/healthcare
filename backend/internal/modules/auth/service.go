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

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: string(hashedPassword),
		FullName:     fullName,
		Role:         Role(role),
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

func (authService *service) Login(ctx context.Context, email, password string) (*User, string, error) {
	user, err := authService.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, "", ErrUserNotFound
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
