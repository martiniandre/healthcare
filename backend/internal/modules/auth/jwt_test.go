package auth

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type failingRevocationStore struct {
	revokeError    error
	isRevokedError error
}

func (store *failingRevocationStore) Revoke(ctx context.Context, tokenDigest string, expiresAt time.Time) error {
	return store.revokeError
}

func (store *failingRevocationStore) IsRevoked(ctx context.Context, tokenDigest string) (bool, error) {
	return false, store.isRevokedError
}

func initJWTForTest(testingInstance *testing.T) {
	testingInstance.Helper()
	require.NoError(testingInstance, InitJWT("jwt-package-test-secret-0123456789"))
}

func TestGenerateJWTIncludesIssuerJTIAndIssuedAt(testingInstance *testing.T) {
	initJWTForTest(testingInstance)

	tokenString, generateError := GenerateJWT("user-1", "admin", "admin@healthcare.com")
	require.NoError(testingInstance, generateError)

	claims, validateError := ValidateJWT(tokenString)
	require.NoError(testingInstance, validateError)
	assert.Equal(testingInstance, jwtIssuer, claims["iss"])
	assert.NotEmpty(testingInstance, claims["jti"])
	assert.NotZero(testingInstance, claims["iat"])
}

func TestValidateJWTRejectsNonHS256Algorithms(testingInstance *testing.T) {
	initJWTForTest(testingInstance)

	claims := jwt.MapClaims{"iss": jwtIssuer, "exp": time.Now().Add(time.Hour).Unix()}
	hs512Token, signError := jwt.NewWithClaims(jwt.SigningMethodHS512, claims).SignedString(jwtSecret)
	require.NoError(testingInstance, signError)

	_, validateError := ValidateJWT(hs512Token)
	assert.Error(testingInstance, validateError)
}

func TestValidateJWTRejectsWrongIssuer(testingInstance *testing.T) {
	initJWTForTest(testingInstance)

	claims := jwt.MapClaims{"iss": "someone-else", "exp": time.Now().Add(time.Hour).Unix()}
	foreignToken, signError := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecret)
	require.NoError(testingInstance, signError)

	_, validateError := ValidateJWT(foreignToken)
	assert.Error(testingInstance, validateError)
}

func TestValidateJWTFailsClosedWhenRevocationStoreErrors(testingInstance *testing.T) {
	initJWTForTest(testingInstance)
	previousStore := revocationStore
	defer SetRevocationStore(previousStore)

	tokenString, generateError := GenerateJWT("user-1", "admin", "admin@healthcare.com")
	require.NoError(testingInstance, generateError)

	SetRevocationStore(&failingRevocationStore{isRevokedError: errors.New("store unavailable")})

	_, validateError := ValidateJWT(tokenString)
	assert.Error(testingInstance, validateError)
}

func TestRevokeTokenStillBlocksTokenWhenStorePersistenceFails(testingInstance *testing.T) {
	initJWTForTest(testingInstance)
	previousStore := revocationStore
	defer SetRevocationStore(previousStore)
	SetRevocationStore(nil)

	tokenString, generateError := GenerateJWT("user-1", "admin", "admin@healthcare.com")
	require.NoError(testingInstance, generateError)

	SetRevocationStore(&failingRevocationStore{revokeError: errors.New("write failed")})
	RevokeToken(tokenString)
	SetRevocationStore(nil)

	_, validateError := ValidateJWT(tokenString)
	assert.Error(testingInstance, validateError)
}
