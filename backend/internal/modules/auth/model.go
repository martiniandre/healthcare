package auth

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleAdmin     Role = "ADMIN"
	RoleDoctor    Role = "DOCTOR"
	RoleNurse     Role = "NURSE"
	RoleReception Role = "RECEPTION"
	RolePatient   Role = "PATIENT"
)

func ParseRole(rawRole string) (Role, bool) {
	switch rawRole {
	case string(RoleAdmin), "RoleAdmin":
		return RoleAdmin, true
	case string(RoleDoctor), "RoleDoctor":
		return RoleDoctor, true
	case string(RoleNurse), "RoleNurse":
		return RoleNurse, true
	case string(RoleReception), "RoleReception":
		return RoleReception, true
	case string(RolePatient), "RolePatient":
		return RolePatient, true
	default:
		return "", false
	}
}

type User struct {
	ID           uuid.UUID `db:"id"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	FullName     string    `db:"full_name"`
	Role         Role      `db:"role"`
	IsActive     bool      `db:"is_active"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
