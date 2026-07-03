package role

import "errors"

type Role string

const (
	RoleAdmin     Role = "ADMIN"
	RoleDoctor    Role = "DOCTOR"
	RoleNurse     Role = "NURSE"
	RoleReception Role = "RECEPTION"
	RolePatient   Role = "PATIENT"
)

var ErrInvalidRole = errors.New("invalid role")

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

func (r Role) String() string {
	return string(r)
}
