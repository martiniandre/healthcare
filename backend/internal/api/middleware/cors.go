package middleware

import (
	"net/http"
	"os"
)

func CORS(secureCookies bool) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			allowedOrigin := os.Getenv("FRONTEND_URL")
			if allowedOrigin == "" {
				allowedOrigin = "http://localhost:5173" // fallback
			}
			writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			writer.Header().Set("Access-Control-Allow-Credentials", "true")
			writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-CSRF-Token, Authorization")
			writer.Header().Set("Vary", "Origin")
			writer.Header().Set("X-Content-Type-Options", "nosniff")
			writer.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
			writer.Header().Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'")

			if request.Method == http.MethodOptions {
				writer.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(writer, request)
		})
	}
}
