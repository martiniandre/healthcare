package auth

import (
	"errors"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtSecret      []byte
	jwtInitOnce    sync.Once
	jwtInitErr     error
	tokenBlacklist sync.Map
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

func GenerateJWT(userID, role, email string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"email":   email,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	})

	return token.SignedString(jwtSecret)
}

func ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	if isTokenRevoked(tokenString) {
		return nil, errors.New("token revoked")
	}

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

	expirationClaim, hasExp := claims["exp"].(float64)
	if !hasExp {
		return nil, errors.New("token missing expiration")
	}
	if time.Now().Unix() > int64(expirationClaim) {
		return nil, errors.New("token expired")
	}

	return claims, nil
}

func RevokeToken(tokenString string) {
	token, _, parseError := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if parseError != nil {
		return
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		if exp, ok := claims["exp"].(float64); ok {
			tokenBlacklist.Store(tokenString, int64(exp))
		}
	}
}

func isTokenRevoked(tokenString string) bool {
	expValue, loaded := tokenBlacklist.Load(tokenString)
	if !loaded {
		return false
	}
	exp, ok := expValue.(int64)
	if !ok {
		return false
	}
	if time.Now().Unix() > exp {
		tokenBlacklist.Delete(tokenString)
		return false
	}
	return true
}
