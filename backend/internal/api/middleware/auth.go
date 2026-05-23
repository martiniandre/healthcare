package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/healthcare/backend/internal/modules/auth"
	"github.com/healthcare/backend/internal/shared/ctxkeys"
)

func ValidateHTTPAuth(httpResponseWriter http.ResponseWriter, httpRequest *http.Request, allowedRoles []auth.Role) (context.Context, bool) {
	cookie, cookieError := httpRequest.Cookie("token")
	if cookieError != nil {
		httpResponseWriter.Header().Set("Content-Type", "application/json")
		httpResponseWriter.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Não autenticado."})
		return nil, false
	}

	claims, jwtValidationErr := auth.ValidateJWT(cookie.Value)
	if jwtValidationErr != nil {
		httpResponseWriter.Header().Set("Content-Type", "application/json")
		httpResponseWriter.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Sessão expirada."})
		return nil, false
	}

	roleStr, roleClaimExists := claims["role"].(string)
	if !roleClaimExists {
		httpResponseWriter.Header().Set("Content-Type", "application/json")
		httpResponseWriter.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Função de usuário inválida."})
		return nil, false
	}

	callerRole := auth.Role(roleStr)
	roleAllowed := false
	for _, allowedRole := range allowedRoles {
		if callerRole == allowedRole {
			roleAllowed = true
			break
		}
	}

	if !roleAllowed {
		httpResponseWriter.Header().Set("Content-Type", "application/json")
		httpResponseWriter.WriteHeader(http.StatusForbidden)
		json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Acesso negado."})
		return nil, false
	}

	if httpRequest.Method == http.MethodPost || httpRequest.Method == http.MethodPut || httpRequest.Method == http.MethodDelete {
		csrfHeader := httpRequest.Header.Get("X-CSRF-Token")
		csrfCookie, csrfCookieErr := httpRequest.Cookie("csrf_token")
		if csrfCookieErr != nil || csrfHeader == "" || csrfHeader != csrfCookie.Value {
			httpResponseWriter.Header().Set("Content-Type", "application/json")
			httpResponseWriter.WriteHeader(http.StatusForbidden)
			json.NewEncoder(httpResponseWriter).Encode(map[string]string{"error": "Token CSRF inválido ou ausente."})
			return nil, false
		}
	}

	userIDStr, _ := claims["user_id"].(string)
	contextWithValues := context.WithValue(httpRequest.Context(), ctxkeys.UserIDKey, userIDStr)
	contextWithValues = context.WithValue(contextWithValues, ctxkeys.RoleKey, roleStr)
	return contextWithValues, true
}
