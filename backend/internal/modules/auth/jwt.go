package auth

import (
	"errors"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtSecret     []byte
	jwtInitOnce   sync.Once
	jwtInitErr    error
)

func InitJWT(secretKey string) error {
	jwtInitOnce.Do(func() {
		if secretKey == "" {
			jwtInitErr = errors.New("JWT_SECRET cannot be empty")
			return
		}
		jwtSecret = []byte(secretKey)
	})
	return jwtInitErr
}

func GenerateJWT(userID, role string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	return token.SignedString(jwtSecret)
}

func ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	token, validationError := jwt.Parse(tokenString, func(parsedToken *jwt.Token) (interface{}, error) {
		if _, isMethodHMAC := parsedToken.Method.(*jwt.SigningMethodHMAC); !isMethodHMAC {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})

	if validationError != nil || !token.Valid {
		return nil, validationError
	}

	claims, isClaimsValid := token.Claims.(jwt.MapClaims)
	if !isClaimsValid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
