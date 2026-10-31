package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/api/render"
)

const (
	accessTokenCookieName  = "token"
	csrfTokenCookieName    = "csrf_token"
	refreshTokenCookieName = "refresh_token"
	accessTokenCookiePath  = "/"
	csrfTokenCookiePath    = "/"
	refreshTokenCookiePath = "/api/v1/auth"
)

type HTTPHandler struct {
	service       Service
	secureCookies bool
}

func NewHTTPHandler(service Service, secureCookies bool) *HTTPHandler {
	return &HTTPHandler{
		service:       service,
		secureCookies: secureCookies,
	}
}

func hashEmailForLogging(email string) string {
	emailDigest := sha256.Sum256([]byte(strings.ToLower(strings.TrimSpace(email))))
	return hex.EncodeToString(emailDigest[:])[:12]
}

func (handler *HTTPHandler) sameSiteMode() http.SameSite {
	return http.SameSiteLaxMode
}

func (handler *HTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/v1/auth/login", handler.HandleLogin)
	mux.HandleFunc("POST /api/v1/auth/refresh", handler.HandleRefresh)
	mux.HandleFunc("POST /api/v1/auth/logout", handler.HandleLogout)
	mux.HandleFunc("GET /api/v1/auth/me", handler.HandleMe)
}

func isRefreshRejection(refreshError error) bool {
	return errors.Is(refreshError, ErrRefreshTokenInvalid) ||
		errors.Is(refreshError, ErrRefreshTokenReused) ||
		errors.Is(refreshError, ErrRefreshTokenExpired) ||
		errors.Is(refreshError, ErrUserInactive)
}

func extractClientInfo(httpRequest *http.Request) ClientInfo {
	remoteHost, _, splitError := net.SplitHostPort(httpRequest.RemoteAddr)
	if splitError != nil {
		remoteHost = httpRequest.RemoteAddr
	}
	return ClientInfo{UserAgent: httpRequest.UserAgent(), IP: remoteHost}
}

func hasValidCSRFToken(httpRequest *http.Request) bool {
	csrfHeader := httpRequest.Header.Get("X-CSRF-Token")
	csrfCookie, csrfCookieError := httpRequest.Cookie(csrfTokenCookieName)
	return csrfCookieError == nil && csrfHeader != "" && csrfHeader == csrfCookie.Value
}

func (handler *HTTPHandler) writeSessionCookies(httpResponseWriter http.ResponseWriter, sessionTokens *SessionTokens) string {
	remainingSessionSeconds := int(time.Until(sessionTokens.SessionExpiresAt).Seconds())
	if remainingSessionSeconds < 1 {
		remainingSessionSeconds = 1
	}
	accessTokenSeconds := int(accessTokenLifetime.Seconds())
	if remainingSessionSeconds < accessTokenSeconds {
		accessTokenSeconds = remainingSessionSeconds
	}

	crossSiteRequestForgeryToken := uuid.New().String()

	http.SetCookie(httpResponseWriter, &http.Cookie{
		Name:     accessTokenCookieName,
		Value:    sessionTokens.AccessToken,
		Path:     accessTokenCookiePath,
		HttpOnly: true,
		Secure:   handler.secureCookies,
		SameSite: handler.sameSiteMode(),
		MaxAge:   accessTokenSeconds,
	})

	http.SetCookie(httpResponseWriter, &http.Cookie{
		Name:     refreshTokenCookieName,
		Value:    sessionTokens.RefreshToken,
		Path:     refreshTokenCookiePath,
		HttpOnly: true,
		Secure:   handler.secureCookies,
		SameSite: handler.sameSiteMode(),
		MaxAge:   remainingSessionSeconds,
	})

	http.SetCookie(httpResponseWriter, &http.Cookie{
		Name:     csrfTokenCookieName,
		Value:    crossSiteRequestForgeryToken,
		Path:     csrfTokenCookiePath,
		HttpOnly: false,
		Secure:   handler.secureCookies,
		SameSite: handler.sameSiteMode(),
		MaxAge:   remainingSessionSeconds,
	})

	return crossSiteRequestForgeryToken
}

func (handler *HTTPHandler) clearSessionCookies(httpResponseWriter http.ResponseWriter) {
	cookiesToClear := []struct {
		name string
		path string
	}{
		{accessTokenCookieName, accessTokenCookiePath},
		{refreshTokenCookieName, refreshTokenCookiePath},
		{csrfTokenCookieName, csrfTokenCookiePath},
	}
	for _, cookieToClear := range cookiesToClear {
		http.SetCookie(httpResponseWriter, &http.Cookie{
			Name:     cookieToClear.name,
			Value:    "",
			Path:     cookieToClear.path,
			HttpOnly: true,
			Secure:   handler.secureCookies,
			SameSite: handler.sameSiteMode(),
			MaxAge:   -1,
		})
	}
}

// HandleLogin godoc
//
//	@Summary		Login
//	@Description	Authenticates a user with email and password, returns JWT token via cookie
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		LoginRequest	true	"Login credentials"
//	@Success		200		{object}	LoginResponse
//	@Failure		400		{object}	map[string]string
//	@Failure		401		{object}	map[string]string
//	@Router			/auth/login [post]
func (handler *HTTPHandler) HandleLogin(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	var loginRequestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if decodeError := json.NewDecoder(httpRequest.Body).Decode(&loginRequestPayload); decodeError != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, "Payload de login inválido.")
		return
	}

	authenticatedUser, sessionTokens, authError := handler.service.Login(httpRequest.Context(), loginRequestPayload.Email, loginRequestPayload.Password, extractClientInfo(httpRequest))
	if authError != nil {
		slog.Warn("Login failed", "email_hash", hashEmailForLogging(loginRequestPayload.Email), "error", authError)
		render.Error(httpResponseWriter, http.StatusUnauthorized, "Credenciais inválidas.")
		return
	}

	crossSiteRequestForgeryToken := handler.writeSessionCookies(httpResponseWriter, sessionTokens)

	render.JSON(httpResponseWriter, http.StatusOK, map[string]interface{}{
		"userId":    authenticatedUser.ID.String(),
		"role":      string(authenticatedUser.Role),
		"email":     authenticatedUser.Email,
		"csrfToken": crossSiteRequestForgeryToken,
	})
}

// HandleRefresh godoc
//
//	@Summary		Refresh session
//	@Description	Rotates the refresh token cookie and issues a new short-lived access token cookie. Requires the X-CSRF-Token header matching the csrf_token cookie. Reuse of an already rotated refresh token revokes the whole session family.
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	LoginResponse
//	@Failure		401	{object}	map[string]string
//	@Failure		403	{object}	map[string]string
//	@Router			/auth/refresh [post]
func (handler *HTTPHandler) HandleRefresh(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	refreshCookie, refreshCookieError := httpRequest.Cookie(refreshTokenCookieName)
	if refreshCookieError != nil || refreshCookie.Value == "" {
		render.Error(httpResponseWriter, http.StatusUnauthorized, "Não autenticado.")
		return
	}

	if !hasValidCSRFToken(httpRequest) {
		render.Error(httpResponseWriter, http.StatusForbidden, "Token CSRF inválido ou ausente.")
		return
	}

	refreshedUser, sessionTokens, refreshError := handler.service.Refresh(httpRequest.Context(), refreshCookie.Value, extractClientInfo(httpRequest))
	if refreshError != nil {
		if isRefreshRejection(refreshError) {
			slog.Warn("Session refresh rejected", "error", refreshError)
			handler.clearSessionCookies(httpResponseWriter)
			render.Error(httpResponseWriter, http.StatusUnauthorized, "Sessão expirada.")
			return
		}
		slog.Error("Session refresh failed", "error", refreshError)
		render.Error(httpResponseWriter, http.StatusInternalServerError, "Erro interno do servidor.")
		return
	}

	crossSiteRequestForgeryToken := handler.writeSessionCookies(httpResponseWriter, sessionTokens)

	render.JSON(httpResponseWriter, http.StatusOK, map[string]interface{}{
		"userId":    refreshedUser.ID.String(),
		"role":      string(refreshedUser.Role),
		"email":     refreshedUser.Email,
		"csrfToken": crossSiteRequestForgeryToken,
	})
}

// HandleLogout godoc
//
//	@Summary		Logout
//	@Description	Revokes the current session (access token and refresh token family) and clears authentication, refresh and CSRF cookies
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]string
//	@Failure		401	{object}	map[string]string
//	@Failure		403	{object}	map[string]string
//	@Router			/auth/logout [post]
func (handler *HTTPHandler) HandleLogout(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	accessCookie, accessCookieError := httpRequest.Cookie(accessTokenCookieName)
	refreshCookie, refreshCookieError := httpRequest.Cookie(refreshTokenCookieName)

	accessTokenIsValid := false
	if accessCookieError == nil {
		_, jwtValidationErr := ValidateJWT(accessCookie.Value)
		accessTokenIsValid = jwtValidationErr == nil
	}
	hasRefreshCookie := refreshCookieError == nil && refreshCookie.Value != ""

	if !accessTokenIsValid && !hasRefreshCookie {
		if accessCookieError != nil {
			render.Error(httpResponseWriter, http.StatusUnauthorized, "Não autenticado.")
			return
		}
		render.Error(httpResponseWriter, http.StatusUnauthorized, "Sessão expirada.")
		return
	}

	if !hasValidCSRFToken(httpRequest) {
		render.Error(httpResponseWriter, http.StatusForbidden, "Token CSRF inválido ou ausente.")
		return
	}

	if accessTokenIsValid {
		RevokeToken(accessCookie.Value)
	}

	if hasRefreshCookie {
		if revokeFamilyError := handler.service.Logout(httpRequest.Context(), refreshCookie.Value); revokeFamilyError != nil {
			slog.Error("Failed to revoke refresh token family on logout", "error", revokeFamilyError)
		}
	}

	handler.clearSessionCookies(httpResponseWriter)

	render.JSON(httpResponseWriter, http.StatusOK, map[string]string{"message": "Logged out successfully"})
}

// HandleMe godoc
//
//	@Summary		Get current user
//	@Description	Returns the authenticated user's profile information
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}
//	@Failure		401	{object}	map[string]string
//	@Router			/auth/me [get]
func (handler *HTTPHandler) HandleMe(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	cookie, cookieError := httpRequest.Cookie("token")
	if cookieError != nil {
		render.Error(httpResponseWriter, http.StatusUnauthorized, "Não autenticado.")
		return
	}

	claims, jwtValidationError := ValidateJWT(cookie.Value)
	if jwtValidationError != nil {
		render.Error(httpResponseWriter, http.StatusUnauthorized, "Sessão expirada.")
		return
	}

	userID, userIDOk := claims["user_id"].(string)
	if !userIDOk || userID == "" {
		render.Error(httpResponseWriter, http.StatusUnauthorized, "Token inválido.")
		return
	}

	user, serviceError := handler.service.Me(httpRequest.Context(), userID)
	if serviceError != nil {
		render.Error(httpResponseWriter, http.StatusUnauthorized, "Usuário não encontrado.")
		return
	}

	csrfCookie, _ := httpRequest.Cookie("csrf_token")
	responsePayload := map[string]interface{}{
		"userId":    user.ID.String(),
		"email":     user.Email,
		"fullName":  user.FullName,
		"role":      string(user.Role),
		"isActive":  user.IsActive,
		"createdAt": user.CreatedAt,
		"updatedAt": user.UpdatedAt,
	}
	if csrfCookie != nil {
		responsePayload["csrfToken"] = csrfCookie.Value
	}

	render.JSON(httpResponseWriter, http.StatusOK, responsePayload)
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	UserID    string `json:"userId"`
	Role      string `json:"role"`
	Email     string `json:"email"`
	CSRFToken string `json:"csrfToken"`
}
