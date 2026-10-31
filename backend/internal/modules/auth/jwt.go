package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"sync/atomic"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	jwtIssuer              = "healthcare-api"
	accessTokenLifetime    = 15 * time.Minute
	jwtKeyIDLength         = 8
	jwtKeyIDHeader         = "kid"
	revocationStoreTimeout = 3 * time.Second
)

type jwtSigningKey struct {
	keyID  string
	secret []byte
}

type jwtKeyring struct {
	currentKey jwtSigningKey
	validKeys  []jwtSigningKey
}

var (
	activeJWTKeyring atomic.Pointer[jwtKeyring]
	jwtInitOnce      sync.Once
	jwtInitErr       error
	tokenBlacklist   sync.Map
	revocationStore  RevocationStore
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

func deriveJWTKeyID(secretKey string) string {
	secretDigest := sha256.Sum256([]byte(secretKey))
	return hex.EncodeToString(secretDigest[:])[:jwtKeyIDLength]
}

func configureJWTKeys(currentSecret string, previousSecrets []string) error {
	if currentSecret == "" {
		return errors.New("JWT_SECRET cannot be empty")
	}

	currentKey := jwtSigningKey{keyID: deriveJWTKeyID(currentSecret), secret: []byte(currentSecret)}
	validKeys := []jwtSigningKey{currentKey}
	for _, previousSecret := range previousSecrets {
		if previousSecret == "" {
			return errors.New("JWT_SECRET_PREVIOUS cannot contain empty entries")
		}
		if previousSecret == currentSecret {
			continue
		}
		validKeys = append(validKeys, jwtSigningKey{keyID: deriveJWTKeyID(previousSecret), secret: []byte(previousSecret)})
	}

	activeJWTKeyring.Store(&jwtKeyring{currentKey: currentKey, validKeys: validKeys})
	return nil
}

func InitJWT(secretKey string, previousSecretKeys ...string) error {
	jwtInitOnce.Do(func() {
		jwtInitErr = configureJWTKeys(secretKey, previousSecretKeys)
	})
	return jwtInitErr
}

func candidateKeysForToken(keyring *jwtKeyring, tokenHeader map[string]interface{}) []jwtSigningKey {
	presentedKeyID, hasKeyID := tokenHeader[jwtKeyIDHeader].(string)
	if !hasKeyID || presentedKeyID == "" {
		return keyring.validKeys
	}

	var matchingKeys []jwtSigningKey
	for _, candidateKey := range keyring.validKeys {
		if candidateKey.keyID == presentedKeyID {
			matchingKeys = append(matchingKeys, candidateKey)
		}
	}
	return matchingKeys
}

func GenerateJWT(userID, role, email string) (string, error) {
	return GenerateJWTBoundedBy(userID, role, email, time.Time{})
}

func GenerateJWTBoundedBy(userID, role, email string, notAfter time.Time) (string, error) {
	keyring := activeJWTKeyring.Load()
	if keyring == nil {
		return "", errors.New("JWT keys not initialized")
	}

	issuedAt := time.Now()
	expiresAt := issuedAt.Add(accessTokenLifetime)
	if !notAfter.IsZero() && notAfter.Before(expiresAt) {
		expiresAt = notAfter
	}
	if !expiresAt.After(issuedAt) {
		return "", errors.New("session lifetime exhausted")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"role":    role,
		"email":   email,
		"iss":     jwtIssuer,
		"jti":     uuid.New().String(),
		"iat":     issuedAt.Unix(),
		"exp":     expiresAt.Unix(),
	})
	token.Header[jwtKeyIDHeader] = keyring.currentKey.keyID

	return token.SignedString(keyring.currentKey.secret)
}

func ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	if isTokenRevoked(tokenString) {
		return nil, errors.New("token revoked")
	}

	keyring := activeJWTKeyring.Load()
	if keyring == nil {
		return nil, errors.New("JWT keys not initialized")
	}

	unverifiedToken, _, headerParseError := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	if headerParseError != nil {
		return nil, headerParseError
	}

	candidateKeys := candidateKeysForToken(keyring, unverifiedToken.Header)
	if len(candidateKeys) == 0 {
		return nil, fmt.Errorf("unknown signing key")
	}

	var validationError error
	for _, candidateKey := range candidateKeys {
		claims, candidateError := parseJWTWithKey(tokenString, candidateKey.secret)
		if candidateError == nil {
			return claims, nil
		}
		validationError = candidateError
	}
	return nil, validationError
}

func parseJWTWithKey(tokenString string, secret []byte) (jwt.MapClaims, error) {
	token, validationError := jwt.Parse(tokenString, func(parsedToken *jwt.Token) (interface{}, error) {
		return secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}), jwt.WithIssuer(jwtIssuer), jwt.WithExpirationRequired())

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
		revokeContext, cancelRevoke := context.WithTimeout(context.Background(), revocationStoreTimeout)
		defer cancelRevoke()
		if revokeError := revocationStore.Revoke(revokeContext, tokenDigest, time.Unix(int64(exp), 0)); revokeError != nil {
			slog.Error("Failed to persist token revocation", "error", revokeError)
		}
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

	checkContext, cancelCheck := context.WithTimeout(context.Background(), revocationStoreTimeout)
	defer cancelCheck()
	revoked, storeError := revocationStore.IsRevoked(checkContext, tokenDigest)
	if storeError != nil {
		slog.Error("Failed to check token revocation, rejecting token", "error", storeError)
		return true
	}
	return revoked
}
