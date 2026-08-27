package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"sync"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	jwtSecret       []byte
	jwtInitOnce     sync.Once
	jwtInitErr      error
	tokenBlacklist  sync.Map
	revocationStore RevocationStore
)

type RevocationStore interface {
	Revoke(ctx context.Context, tokenDigest string, expiresAt time.Time) error
	IsRevoked(ctx context.Context, tokenDigest string) (bool, error)
}

func SetRevocationStore(store RevocationStore) {
	revocationStore = store
}

func digestToken(tokenString string) string {
	tokenDigest := sha256.Sum256([]byte(tokenString))
	return hex.EncodeToString(tokenDigest[:])
}

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
	claims, claimsOk := token.Claims.(jwt.MapClaims)
	if !claimsOk {
		return
	}
	exp, expOk := claims["exp"].(float64)
	if !expOk {
		return
	}

	tokenDigest := digestToken(tokenString)
	tokenBlacklist.Store(tokenDigest, int64(exp))

	if revocationStore != nil {
		_ = revocationStore.Revoke(context.Background(), tokenDigest, time.Unix(int64(exp), 0))
	}
}

func isTokenRevoked(tokenString string) bool {
	tokenDigest := digestToken(tokenString)

	expValue, loaded := tokenBlacklist.Load(tokenDigest)
	if loaded {
		if exp, ok := expValue.(int64); ok {
			if time.Now().Unix() > exp {
				tokenBlacklist.Delete(tokenDigest)
			} else {
				return true
			}
		}
	}

	if revocationStore == nil {
		return false
	}

	revoked, storeError := revocationStore.IsRevoked(context.Background(), tokenDigest)
	if storeError != nil {
		return false
	}
	return revoked
}
