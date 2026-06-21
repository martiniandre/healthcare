package auth

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/healthcare/backend/internal/api/render"
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

func (handler *HTTPHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("POST /api/auth/login", handler.HandleLogin)
	mux.HandleFunc("POST /api/auth/logout", handler.HandleLogout)
	mux.HandleFunc("GET /api/auth/me", handler.HandleMe)
}

func (handler *HTTPHandler) HandleLogin(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	var loginRequestPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if decodeError := json.NewDecoder(httpRequest.Body).Decode(&loginRequestPayload); decodeError != nil {
		render.Error(httpResponseWriter, http.StatusBadRequest, "Payload de login inválido.")
		return
	}

	authenticatedUser, jsonWebToken, authError := handler.service.Login(httpRequest.Context(), loginRequestPayload.Email, loginRequestPayload.Password)
	if authError != nil {
		slog.Warn("Login failed", "email", loginRequestPayload.Email, "error", authError)
		render.Error(httpResponseWriter, http.StatusUnauthorized, "Credenciais inválidas.")
		return
	}

	crossSiteRequestForgeryToken := uuid.New().String()

	http.SetCookie(httpResponseWriter, &http.Cookie{
		Name:     "token",
		Value:    jsonWebToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   handler.secureCookies,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400,
	})

	http.SetCookie(httpResponseWriter, &http.Cookie{
		Name:     "csrf_token",
		Value:    crossSiteRequestForgeryToken,
		Path:     "/",
		HttpOnly: false,
		Secure:   handler.secureCookies,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   86400,
	})

	render.JSON(httpResponseWriter, http.StatusOK, map[string]interface{}{
		"userId": authenticatedUser.ID.String(),
		"role":   string(authenticatedUser.Role),
		"email":  authenticatedUser.Email,
	})
}

func (handler *HTTPHandler) HandleLogout(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
	cookie, cookieError := httpRequest.Cookie("token")
	if cookieError != nil {
		render.Error(httpResponseWriter, http.StatusUnauthorized, "Não autenticado.")
		return
	}

	_, jwtValidationErr := ValidateJWT(cookie.Value)
	if jwtValidationErr != nil {
		render.Error(httpResponseWriter, http.StatusUnauthorized, "Sessão expirada.")
		return
	}

	csrfHeader := httpRequest.Header.Get("X-CSRF-Token")
	csrfCookie, csrfCookieErr := httpRequest.Cookie("csrf_token")
	if csrfCookieErr != nil || csrfHeader == "" || csrfHeader != csrfCookie.Value {
		render.Error(httpResponseWriter, http.StatusForbidden, "Token CSRF inválido ou ausente.")
		return
	}

	http.SetCookie(httpResponseWriter, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   handler.secureCookies,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	http.SetCookie(httpResponseWriter, &http.Cookie{
		Name:     "csrf_token",
		Value:    "",
		Path:     "/",
		HttpOnly: false,
		Secure:   handler.secureCookies,
		SameSite: http.SameSiteLaxMode,
		MaxAge:   -1,
	})

	render.JSON(httpResponseWriter, http.StatusOK, map[string]string{"message": "Logged out successfully"})
}

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

	render.JSON(httpResponseWriter, http.StatusOK, map[string]interface{}{
		"userId":    user.ID.String(),
		"email":     user.Email,
		"fullName":  user.FullName,
		"role":      string(user.Role),
		"isActive":  user.IsActive,
		"createdAt": user.CreatedAt,
		"updatedAt": user.UpdatedAt,
	})
}

