package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	rotationOldSecret = "rotation-old-secret-0123456789abcdef"
	rotationNewSecret = "rotation-new-secret-0123456789abcdef"
)

func restoreKeyringAfterTest(testingInstance *testing.T) {
	testingInstance.Helper()
	previousKeyring := activeJWTKeyring.Load()
	testingInstance.Cleanup(func() { activeJWTKeyring.Store(previousKeyring) })
}

func TestGenerateJWTCarriesKeyIDDerivedFromSecret(testingInstance *testing.T) {
	restoreKeyringAfterTest(testingInstance)
	require.NoError(testingInstance, configureJWTKeys(rotationNewSecret, nil))

	tokenString, generateError := GenerateJWT("user-1", "admin", "admin@healthcare.com")
	require.NoError(testingInstance, generateError)

	parsedToken, _, parseError := new(jwt.Parser).ParseUnverified(tokenString, jwt.MapClaims{})
	require.NoError(testingInstance, parseError)
	assert.Equal(testingInstance, deriveJWTKeyID(rotationNewSecret), parsedToken.Header["kid"])
	assert.Len(testingInstance, deriveJWTKeyID(rotationNewSecret), 8)
}

func TestValidateJWTAcceptsTokenSignedWithPreviousKeyAfterRotation(testingInstance *testing.T) {
	restoreKeyringAfterTest(testingInstance)
	require.NoError(testingInstance, configureJWTKeys(rotationOldSecret, nil))
	tokenSignedWithOldKey, generateError := GenerateJWT("user-1", "admin", "admin@healthcare.com")
	require.NoError(testingInstance, generateError)

	require.NoError(testingInstance, configureJWTKeys(rotationNewSecret, []string{rotationOldSecret}))

	claims, validateError := ValidateJWT(tokenSignedWithOldKey)
	require.NoError(testingInstance, validateError)
	assert.Equal(testingInstance, "user-1", claims["user_id"])

	tokenSignedWithNewKey, newGenerateError := GenerateJWT("user-2", "admin", "admin@healthcare.com")
	require.NoError(testingInstance, newGenerateError)
	parsedToken, _, _ := new(jwt.Parser).ParseUnverified(tokenSignedWithNewKey, jwt.MapClaims{})
	assert.Equal(testingInstance, deriveJWTKeyID(rotationNewSecret), parsedToken.Header["kid"])
	_, newValidateError := ValidateJWT(tokenSignedWithNewKey)
	assert.NoError(testingInstance, newValidateError)
}

func TestValidateJWTRejectsTokenSignedWithRetiredKey(testingInstance *testing.T) {
	restoreKeyringAfterTest(testingInstance)
	require.NoError(testingInstance, configureJWTKeys(rotationOldSecret, nil))
	tokenSignedWithOldKey, generateError := GenerateJWT("user-1", "admin", "admin@healthcare.com")
	require.NoError(testingInstance, generateError)

	require.NoError(testingInstance, configureJWTKeys(rotationNewSecret, nil))

	_, validateError := ValidateJWT(tokenSignedWithOldKey)
	assert.Error(testingInstance, validateError)
}

func TestValidateJWTWithoutKeyIDFallsBackToTryingAllKeys(testingInstance *testing.T) {
	restoreKeyringAfterTest(testingInstance)
	require.NoError(testingInstance, configureJWTKeys(rotationNewSecret, []string{rotationOldSecret}))

	legacyClaims := jwt.MapClaims{"iss": jwtIssuer, "exp": time.Now().Add(time.Hour).Unix(), "user_id": "legacy"}
	legacyToken, signError := jwt.NewWithClaims(jwt.SigningMethodHS256, legacyClaims).SignedString([]byte(rotationOldSecret))
	require.NoError(testingInstance, signError)

	claims, validateError := ValidateJWT(legacyToken)
	require.NoError(testingInstance, validateError)
	assert.Equal(testingInstance, "legacy", claims["user_id"])
}

func TestValidateJWTRejectsUnknownKeyID(testingInstance *testing.T) {
	restoreKeyringAfterTest(testingInstance)
	require.NoError(testingInstance, configureJWTKeys(rotationNewSecret, nil))

	foreignToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"iss": jwtIssuer, "exp": time.Now().Add(time.Hour).Unix()})
	foreignToken.Header["kid"] = "deadbeef"
	foreignTokenString, signError := foreignToken.SignedString([]byte(rotationNewSecret))
	require.NoError(testingInstance, signError)

	_, validateError := ValidateJWT(foreignTokenString)
	assert.Error(testingInstance, validateError)
}

func TestConfigureJWTKeysRejectsEmptySecrets(testingInstance *testing.T) {
	restoreKeyringAfterTest(testingInstance)
	assert.Error(testingInstance, configureJWTKeys("", nil))
	assert.Error(testingInstance, configureJWTKeys(rotationNewSecret, []string{""}))
}

func TestGenerateJWTBoundedByNeverExceedsSessionDeadline(testingInstance *testing.T) {
	restoreKeyringAfterTest(testingInstance)
	require.NoError(testingInstance, configureJWTKeys(rotationNewSecret, nil))

	deadline := time.Now().Add(2 * time.Minute)
	tokenString, generateError := GenerateJWTBoundedBy("user-1", "admin", "admin@healthcare.com", deadline)
	require.NoError(testingInstance, generateError)
	claims, validateError := ValidateJWT(tokenString)
	require.NoError(testingInstance, validateError)
	assert.Equal(testingInstance, float64(deadline.Unix()), claims["exp"])

	unboundedToken, _ := GenerateJWT("user-1", "admin", "admin@healthcare.com")
	unboundedClaims, _ := ValidateJWT(unboundedToken)
	assert.InDelta(testingInstance, float64(time.Now().Add(accessTokenLifetime).Unix()), unboundedClaims["exp"], 5)

	_, exhaustedError := GenerateJWTBoundedBy("user-1", "admin", "admin@healthcare.com", time.Now().Add(-time.Second))
	assert.Error(testingInstance, exhaustedError)
}
