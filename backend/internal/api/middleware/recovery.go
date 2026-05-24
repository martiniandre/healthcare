package middleware

import (
	"log/slog"
	"net/http"
	"runtime/debug"
)

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(httpResponseWriter http.ResponseWriter, httpRequest *http.Request) {
		defer func() {
			recoveredValue := recover()
			if recoveredValue != nil {
				slog.Error("panic recovered in HTTP handler",
					"panic", recoveredValue,
					"path", httpRequest.URL.Path,
					"method", httpRequest.Method,
					"stack", string(debug.Stack()),
				)
				http.Error(httpResponseWriter, "Internal Server Error", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(httpResponseWriter, httpRequest)
	})
}
