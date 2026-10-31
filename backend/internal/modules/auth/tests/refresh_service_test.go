package tests

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/healthcare/backend/internal/modules/auth"
	"github.com/healthcare/backend/internal/modules/auth/mocks"
	"github.com/healthcare/backend/internal/shared/role"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newLoggedInSession(testingInstance *testing.T) (auth.Service, *mocks.MockRepository, *auth.User, *auth.SessionTokens) {
	testingInstance.Helper()
	mockRepository := mocks.NewMockRepository()
	authService := auth.NewService(mockRepository, nil)
	contextParam := context.Background()

	_, registerError := authService.Register(contextParam, "refresh@example.com", "securepass", "Refresh User", string(role.RoleNurse))
	require.NoError(testingInstance, registerError)

	user, sessionTokens, loginError := authService.Login(contextParam, "refresh@example.com", "securepass", auth.ClientInfo{UserAgent: "test-agent", IP: "10.0.0.1"})
	require.NoError(testingInstance, loginError)
	return authService, mockRepository, user, sessionTokens
}

func TestService_Login_StoresOnlyRefreshTokenHash(testingInstance *testing.T) {
	_, mockRepository, user, sessionTokens := newLoggedInSession(testingInstance)

	require.Len(testingInstance, mockRepository.RefreshTokens, 1)
	storedToken, exists := mockRepository.RefreshTokens[auth.HashRefreshToken(sessionTokens.RefreshToken)]
	require.True(testingInstance, exists)
	assert.NotEqual(testingInstance, sessionTokens.RefreshToken, storedToken.TokenHash)
	assert.Equal(testingInstance, user.ID, storedToken.UserID)
	assert.Equal(testingInstance, "test-agent", storedToken.UserAgent)
	assert.WithinDuration(testingInstance, time.Now().Add(12*time.Hour), storedToken.AbsoluteExpiresAt, time.Minute)
}

func TestService_Refresh_RotatesToken(testingInstance *testing.T) {
	authService, mockRepository, _, sessionTokens := newLoggedInSession(testingInstance)

	_, refreshedTokens, refreshError := authService.Refresh(context.Background(), sessionTokens.RefreshToken, auth.ClientInfo{})

	require.NoError(testingInstance, refreshError)
	assert.NotEqual(testingInstance, sessionTokens.RefreshToken, refreshedTokens.RefreshToken)
	assert.NotEmpty(testingInstance, refreshedTokens.AccessToken)
	assert.Len(testingInstance, mockRepository.RefreshTokens, 2)

	oldToken := mockRepository.RefreshTokens[auth.HashRefreshToken(sessionTokens.RefreshToken)]
	newToken := mockRepository.RefreshTokens[auth.HashRefreshToken(refreshedTokens.RefreshToken)]
	assert.NotNil(testingInstance, oldToken.ReplacedBy)
	assert.Equal(testingInstance, oldToken.FamilyID, newToken.FamilyID)
	assert.Equal(testingInstance, oldToken.AbsoluteExpiresAt, newToken.AbsoluteExpiresAt)
}

func TestService_Refresh_ReuseRevokesWholeFamily(testingInstance *testing.T) {
	authService, mockRepository, _, sessionTokens := newLoggedInSession(testingInstance)
	contextParam := context.Background()

	_, refreshedTokens, firstRefreshError := authService.Refresh(contextParam, sessionTokens.RefreshToken, auth.ClientInfo{})
	require.NoError(testingInstance, firstRefreshError)

	_, _, reuseError := authService.Refresh(contextParam, sessionTokens.RefreshToken, auth.ClientInfo{})
	assert.ErrorIs(testingInstance, reuseError, auth.ErrRefreshTokenReused)

	for _, storedToken := range mockRepository.RefreshTokens {
		assert.NotNil(testingInstance, storedToken.RevokedAt)
	}

	_, _, legitimateAfterReuseError := authService.Refresh(contextParam, refreshedTokens.RefreshToken, auth.ClientInfo{})
	assert.ErrorIs(testingInstance, legitimateAfterReuseError, auth.ErrRefreshTokenReused)
}

func TestService_Refresh_UnknownTokenIsInvalid(testingInstance *testing.T) {
	authService, _, _, _ := newLoggedInSession(testingInstance)

	_, _, refreshError := authService.Refresh(context.Background(), "not-a-real-token", auth.ClientInfo{})
	assert.ErrorIs(testingInstance, refreshError, auth.ErrRefreshTokenInvalid)

	_, _, emptyError := authService.Refresh(context.Background(), "", auth.ClientInfo{})
	assert.ErrorIs(testingInstance, emptyError, auth.ErrRefreshTokenInvalid)
}

func TestService_Refresh_RejectsInactiveUserAndRevokesFamily(testingInstance *testing.T) {
	authService, mockRepository, user, sessionTokens := newLoggedInSession(testingInstance)
	user.IsActive = false

	_, _, refreshError := authService.Refresh(context.Background(), sessionTokens.RefreshToken, auth.ClientInfo{})

	assert.ErrorIs(testingInstance, refreshError, auth.ErrUserInactive)
	storedToken := mockRepository.RefreshTokens[auth.HashRefreshToken(sessionTokens.RefreshToken)]
	assert.NotNil(testingInstance, storedToken.RevokedAt)
}

func TestService_Refresh_ReflectsCurrentRoleFromDatabase(testingInstance *testing.T) {
	authService, _, user, sessionTokens := newLoggedInSession(testingInstance)
	user.Role = role.RoleDoctor

	refreshedUser, refreshedTokens, refreshError := authService.Refresh(context.Background(), sessionTokens.RefreshToken, auth.ClientInfo{})

	require.NoError(testingInstance, refreshError)
	assert.Equal(testingInstance, role.RoleDoctor, refreshedUser.Role)
	claims, validateError := auth.ValidateJWT(refreshedTokens.AccessToken)
	require.NoError(testingInstance, validateError)
	assert.Equal(testingInstance, string(role.RoleDoctor), claims["role"])
}

func TestService_Refresh_RejectsAfterAbsoluteExpiry(testingInstance *testing.T) {
	authService, mockRepository, _, sessionTokens := newLoggedInSession(testingInstance)
	storedToken := mockRepository.RefreshTokens[auth.HashRefreshToken(sessionTokens.RefreshToken)]
	storedToken.AbsoluteExpiresAt = time.Now().Add(-time.Minute)

	_, _, refreshError := authService.Refresh(context.Background(), sessionTokens.RefreshToken, auth.ClientInfo{})

	assert.ErrorIs(testingInstance, refreshError, auth.ErrRefreshTokenExpired)
}

func TestService_Refresh_AccessTokenNeverOutlivesAbsoluteExpiry(testingInstance *testing.T) {
	authService, mockRepository, _, sessionTokens := newLoggedInSession(testingInstance)
	storedToken := mockRepository.RefreshTokens[auth.HashRefreshToken(sessionTokens.RefreshToken)]
	deadline := time.Now().Add(3 * time.Minute)
	storedToken.AbsoluteExpiresAt = deadline
	storedToken.ExpiresAt = deadline

	_, refreshedTokens, refreshError := authService.Refresh(context.Background(), sessionTokens.RefreshToken, auth.ClientInfo{})

	require.NoError(testingInstance, refreshError)
	claims, validateError := auth.ValidateJWT(refreshedTokens.AccessToken)
	require.NoError(testingInstance, validateError)
	assert.LessOrEqual(testingInstance, claims["exp"].(float64), float64(deadline.Unix()))
}

func TestService_Refresh_PropagatesRepositoryFailure(testingInstance *testing.T) {
	authService, mockRepository, _, sessionTokens := newLoggedInSession(testingInstance)
	mockRepository.Err = errors.New("database unavailable")

	_, _, refreshError := authService.Refresh(context.Background(), sessionTokens.RefreshToken, auth.ClientInfo{})

	assert.ErrorIs(testingInstance, refreshError, mockRepository.Err)
}

func TestService_Logout_RevokesOnlyCurrentFamily(testingInstance *testing.T) {
	authService, mockRepository, _, firstSession := newLoggedInSession(testingInstance)
	contextParam := context.Background()
	_, secondSession, secondLoginError := authService.Login(contextParam, "refresh@example.com", "securepass", auth.ClientInfo{})
	require.NoError(testingInstance, secondLoginError)

	require.NoError(testingInstance, authService.Logout(contextParam, firstSession.RefreshToken))

	assert.NotNil(testingInstance, mockRepository.RefreshTokens[auth.HashRefreshToken(firstSession.RefreshToken)].RevokedAt)
	assert.Nil(testingInstance, mockRepository.RefreshTokens[auth.HashRefreshToken(secondSession.RefreshToken)].RevokedAt)

	_, _, refreshAfterLogoutError := authService.Refresh(contextParam, firstSession.RefreshToken, auth.ClientInfo{})
	assert.Error(testingInstance, refreshAfterLogoutError)
	_, _, otherSessionRefreshError := authService.Refresh(contextParam, secondSession.RefreshToken, auth.ClientInfo{})
	assert.NoError(testingInstance, otherSessionRefreshError)
}

func TestService_Logout_UnknownTokenIsNoop(testingInstance *testing.T) {
	authService, _, _, _ := newLoggedInSession(testingInstance)

	assert.NoError(testingInstance, authService.Logout(context.Background(), "unknown"))
	assert.NoError(testingInstance, authService.Logout(context.Background(), ""))
}
